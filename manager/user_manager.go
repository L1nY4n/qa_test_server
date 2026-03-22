package manager

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrDBUnavailable      = errors.New("database is unavailable")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserDisabled       = errors.New("user is disabled")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUsernameTaken      = errors.New("username already exists")
	ErrInvalidRole        = errors.New("invalid role")
	ErrWeakPassword       = errors.New("password must be at least 6 characters")
	ErrInvalidUsername    = errors.New("username must be 3-32 chars, letters/digits/._-")
)

var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_.-]{3,32}$`)

type TokenClaims struct {
	UserID   uint   `json:"uid"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Exp      int64  `json:"exp"`
	Iat      int64  `json:"iat"`
}

type UserUpdate struct {
	DisplayName *string
	Role        *string
	Enabled     *bool
}

type UserManager struct {
	mu       sync.RWMutex
	secret   []byte
	tokenTTL time.Duration
}

var UserManagerGlobal = &UserManager{}

func (m *UserManager) Init(secret string, tokenTTL time.Duration, defaultAdminUser, defaultAdminPassword string) error {
	if db.DB == nil {
		return ErrDBUnavailable
	}
	secret = strings.TrimSpace(secret)
	if secret == "" {
		return errors.New("auth secret is required")
	}
	if tokenTTL <= 0 {
		tokenTTL = 24 * time.Hour
	}

	m.mu.Lock()
	m.secret = []byte(secret)
	m.tokenTTL = tokenTTL
	m.mu.Unlock()

	if err := db.DB.AutoMigrate(&model.User{}); err != nil {
		return fmt.Errorf("migrate users failed: %w", err)
	}

	var total int64
	if err := db.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		return fmt.Errorf("count users failed: %w", err)
	}
	if total > 0 {
		return nil
	}

	defaultAdminUser = strings.TrimSpace(defaultAdminUser)
	if defaultAdminUser == "" {
		defaultAdminUser = "admin"
	}
	defaultAdminPassword = strings.TrimSpace(defaultAdminPassword)
	if defaultAdminPassword == "" {
		defaultAdminPassword = "Admin@123456"
	}

	_, err := m.Create(defaultAdminUser, "系统管理员", defaultAdminPassword, model.RoleAdmin)
	if err != nil {
		return fmt.Errorf("create default admin failed: %w", err)
	}
	return nil
}

func (m *UserManager) Register(username, displayName, password string) (model.PublicUser, error) {
	return m.Create(username, displayName, password, model.RoleViewer)
}

func (m *UserManager) Create(username, displayName, password, role string) (model.PublicUser, error) {
	if db.DB == nil {
		return model.PublicUser{}, ErrDBUnavailable
	}

	username = strings.TrimSpace(username)
	if !usernameRegexp.MatchString(username) {
		return model.PublicUser{}, ErrInvalidUsername
	}
	if len(password) < 6 {
		return model.PublicUser{}, ErrWeakPassword
	}

	role, err := normalizeRole(role)
	if err != nil {
		return model.PublicUser{}, err
	}
	if strings.TrimSpace(displayName) == "" {
		displayName = username
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return model.PublicUser{}, fmt.Errorf("hash password failed: %w", err)
	}

	user := model.User{
		Username:     username,
		DisplayName:  strings.TrimSpace(displayName),
		PasswordHash: string(hash),
		Role:         role,
		Enabled:      true,
	}
	if err := db.DB.Create(&user).Error; err != nil {
		if isUniqueViolation(err) {
			return model.PublicUser{}, ErrUsernameTaken
		}
		return model.PublicUser{}, err
	}
	return user.Public(), nil
}

func (m *UserManager) Login(username, password string) (string, time.Time, model.PublicUser, error) {
	if db.DB == nil {
		return "", time.Time{}, model.PublicUser{}, ErrDBUnavailable
	}
	username = strings.TrimSpace(username)
	if username == "" || strings.TrimSpace(password) == "" {
		return "", time.Time{}, model.PublicUser{}, ErrInvalidCredentials
	}

	var user model.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", time.Time{}, model.PublicUser{}, ErrInvalidCredentials
		}
		return "", time.Time{}, model.PublicUser{}, err
	}
	if !user.Enabled {
		return "", time.Time{}, model.PublicUser{}, ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", time.Time{}, model.PublicUser{}, ErrInvalidCredentials
	}

	token, expiresAt, err := m.issueToken(user)
	if err != nil {
		return "", time.Time{}, model.PublicUser{}, err
	}
	return token, expiresAt, user.Public(), nil
}

func (m *UserManager) ParseToken(token string) (TokenClaims, error) {
	secret, _, err := m.authConfig()
	if err != nil {
		return TokenClaims{}, err
	}

	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return TokenClaims{}, errors.New("invalid token")
	}

	payloadEncoded := parts[0]
	providedSigEncoded := parts[1]
	providedSig, err := base64.RawURLEncoding.DecodeString(providedSigEncoded)
	if err != nil {
		return TokenClaims{}, errors.New("invalid token signature")
	}
	expectedSig := signPayload(secret, payloadEncoded)
	if !hmac.Equal(providedSig, expectedSig) {
		return TokenClaims{}, errors.New("invalid token signature")
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEncoded)
	if err != nil {
		return TokenClaims{}, errors.New("invalid token payload")
	}

	claims := TokenClaims{}
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return TokenClaims{}, errors.New("invalid token claims")
	}
	if claims.Exp <= time.Now().Unix() {
		return TokenClaims{}, errors.New("token expired")
	}
	return claims, nil
}

func (m *UserManager) GetByID(id uint) (model.PublicUser, error) {
	if db.DB == nil {
		return model.PublicUser{}, ErrDBUnavailable
	}
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PublicUser{}, ErrUserNotFound
		}
		return model.PublicUser{}, err
	}
	return user.Public(), nil
}

func (m *UserManager) List(keyword string, offset, limit int) ([]model.PublicUser, int64, error) {
	if db.DB == nil {
		return nil, 0, ErrDBUnavailable
	}
	offset, limit = normalizePageRange(offset, limit)
	keyword = strings.TrimSpace(keyword)

	query := db.DB.Model(&model.User{})
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("username LIKE ? OR display_name LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]model.User, 0, limit)
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}

	result := make([]model.PublicUser, 0, len(items))
	for _, item := range items {
		result = append(result, item.Public())
	}
	return result, total, nil
}

func (m *UserManager) Update(id uint, actorID uint, patch UserUpdate) (model.PublicUser, error) {
	if db.DB == nil {
		return model.PublicUser{}, ErrDBUnavailable
	}
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.PublicUser{}, ErrUserNotFound
		}
		return model.PublicUser{}, err
	}

	updates := map[string]interface{}{}
	if patch.DisplayName != nil {
		name := strings.TrimSpace(*patch.DisplayName)
		if name == "" {
			name = user.Username
		}
		updates["display_name"] = name
	}

	if patch.Role != nil {
		newRole, err := normalizeRole(*patch.Role)
		if err != nil {
			return model.PublicUser{}, err
		}
		if actorID == user.ID && user.Role == model.RoleAdmin && newRole != model.RoleAdmin {
			return model.PublicUser{}, errors.New("cannot downgrade current admin account")
		}
		if user.Role == model.RoleAdmin && newRole != model.RoleAdmin {
			ok, err := m.hasOtherEnabledAdmin(user.ID)
			if err != nil {
				return model.PublicUser{}, err
			}
			if !ok {
				return model.PublicUser{}, errors.New("at least one enabled admin is required")
			}
		}
		updates["role"] = newRole
	}

	if patch.Enabled != nil {
		enabled := *patch.Enabled
		if actorID == user.ID && !enabled {
			return model.PublicUser{}, errors.New("cannot disable current account")
		}
		if user.Role == model.RoleAdmin && !enabled {
			ok, err := m.hasOtherEnabledAdmin(user.ID)
			if err != nil {
				return model.PublicUser{}, err
			}
			if !ok {
				return model.PublicUser{}, errors.New("at least one enabled admin is required")
			}
		}
		updates["enabled"] = enabled
	}

	if len(updates) == 0 {
		return user.Public(), nil
	}

	if err := db.DB.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return model.PublicUser{}, err
	}

	if err := db.DB.First(&user, id).Error; err != nil {
		return model.PublicUser{}, err
	}
	return user.Public(), nil
}

func (m *UserManager) Delete(id uint, actorID uint) error {
	if db.DB == nil {
		return ErrDBUnavailable
	}
	if id == actorID {
		return errors.New("cannot delete current account")
	}
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}
	if user.Role == model.RoleAdmin {
		ok, err := m.hasOtherEnabledAdmin(user.ID)
		if err != nil {
			return err
		}
		if !ok {
			return errors.New("at least one enabled admin is required")
		}
	}

	return db.DB.Delete(&model.User{}, id).Error
}

func (m *UserManager) ResetPassword(id uint, newPassword string) error {
	if db.DB == nil {
		return ErrDBUnavailable
	}
	if len(newPassword) < 6 {
		return ErrWeakPassword
	}
	var user model.User
	if err := db.DB.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password failed: %w", err)
	}
	return db.DB.Model(&model.User{}).Where("id = ?", id).Update("password_hash", string(hash)).Error
}

func (m *UserManager) issueToken(user model.User) (string, time.Time, error) {
	secret, ttl, err := m.authConfig()
	if err != nil {
		return "", time.Time{}, err
	}
	now := time.Now()
	expiresAt := now.Add(ttl)
	claims := TokenClaims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     user.Role,
		Iat:      now.Unix(),
		Exp:      expiresAt.Unix(),
	}

	payloadBytes, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}
	payload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	sig := base64.RawURLEncoding.EncodeToString(signPayload(secret, payload))
	return payload + "." + sig, expiresAt, nil
}

func (m *UserManager) authConfig() ([]byte, time.Duration, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if len(m.secret) == 0 {
		return nil, 0, errors.New("user manager is not initialized")
	}
	if m.tokenTTL <= 0 {
		return nil, 0, errors.New("invalid token ttl")
	}
	return m.secret, m.tokenTTL, nil
}

func (m *UserManager) hasOtherEnabledAdmin(excludeID uint) (bool, error) {
	var count int64
	err := db.DB.Model(&model.User{}).
		Where("role = ? AND enabled = ? AND id <> ?", model.RoleAdmin, true, excludeID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func normalizeRole(role string) (string, error) {
	role = strings.ToLower(strings.TrimSpace(role))
	if role == "" {
		role = model.RoleViewer
	}
	if !model.IsValidRole(role) {
		return "", ErrInvalidRole
	}
	return role, nil
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	message := strings.ToLower(err.Error())
	return strings.Contains(message, "duplicate") || strings.Contains(message, "unique")
}

func normalizePageRange(offset, limit int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 200 {
		limit = 200
	}
	return offset, limit
}

func signPayload(secret []byte, payload string) []byte {
	h := hmac.New(sha256.New, secret)
	_, _ = h.Write([]byte(payload))
	return h.Sum(nil)
}
