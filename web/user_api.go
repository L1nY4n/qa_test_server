package web

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"qa_test_server/manager"

	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type registerRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
}

type createUserRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Password    string `json:"password"`
	Role        string `json:"role"`
}

type updateUserRequest struct {
	DisplayName *string `json:"displayName"`
	Role        *string `json:"role"`
	Enabled     *bool   `json:"enabled"`
}

type resetPasswordRequest struct {
	Password string `json:"password"`
}

func register(c *gin.Context) {
	req := registerRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid register payload")
		return
	}

	if _, err := manager.UserManagerGlobal.Register(req.Username, req.DisplayName, req.Password); err != nil {
		handleUserErr(c, err)
		return
	}

	token, expiresAt, user, err := manager.UserManagerGlobal.Login(req.Username, req.Password)
	if err != nil {
		handleUserErr(c, err)
		return
	}

	ok(c, gin.H{
		"token":     token,
		"expiresAt": expiresAt,
		"user":      user,
	})
}

func login(c *gin.Context) {
	req := loginRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid login payload")
		return
	}

	token, expiresAt, user, err := manager.UserManagerGlobal.Login(req.Username, req.Password)
	if err != nil {
		handleUserErr(c, err)
		return
	}

	ok(c, gin.H{
		"token":     token,
		"expiresAt": expiresAt,
		"user":      user,
	})
}

func profile(c *gin.Context) {
	claims, hasClaims := currentClaims(c)
	if !hasClaims {
		fail(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := manager.UserManagerGlobal.GetByID(claims.UserID)
	if err != nil {
		handleUserErr(c, err)
		return
	}

	ok(c, gin.H{
		"user": user,
	})
}

func userList(c *gin.Context) {
	keyword := c.Query("keyword")
	offset := queryInt(c, "offset", 0)
	limit := queryInt(c, "limit", 20)

	items, total, err := manager.UserManagerGlobal.List(keyword, offset, limit)
	if err != nil {
		handleUserErr(c, err)
		return
	}

	ok(c, gin.H{
		"items":  items,
		"total":  total,
		"offset": offset,
		"limit":  limit,
	})
}

func createUser(c *gin.Context) {
	req := createUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid create user payload")
		return
	}

	user, err := manager.UserManagerGlobal.Create(req.Username, req.DisplayName, req.Password, req.Role)
	if err != nil {
		handleUserErr(c, err)
		return
	}
	ok(c, user)
}

func updateUser(c *gin.Context) {
	id, validID := parseUserID(c.Param("id"))
	if !validID {
		fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	req := updateUserRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid update user payload")
		return
	}

	claims, _ := currentClaims(c)
	user, err := manager.UserManagerGlobal.Update(id, claims.UserID, manager.UserUpdate{
		DisplayName: req.DisplayName,
		Role:        req.Role,
		Enabled:     req.Enabled,
	})
	if err != nil {
		handleUserErr(c, err)
		return
	}
	ok(c, user)
}

func resetUserPassword(c *gin.Context) {
	id, validID := parseUserID(c.Param("id"))
	if !validID {
		fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	req := resetPasswordRequest{}
	if err := c.ShouldBindJSON(&req); err != nil {
		fail(c, http.StatusBadRequest, "invalid reset password payload")
		return
	}

	if err := manager.UserManagerGlobal.ResetPassword(id, req.Password); err != nil {
		handleUserErr(c, err)
		return
	}
	ok(c, gin.H{"id": id, "message": "password updated"})
}

func deleteUser(c *gin.Context) {
	id, validID := parseUserID(c.Param("id"))
	if !validID {
		fail(c, http.StatusBadRequest, "invalid user id")
		return
	}

	claims, _ := currentClaims(c)
	if err := manager.UserManagerGlobal.Delete(id, claims.UserID); err != nil {
		handleUserErr(c, err)
		return
	}
	ok(c, gin.H{"id": id, "message": "deleted"})
}

func parseUserID(raw string) (uint, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, false
	}
	id64, err := strconv.ParseUint(raw, 10, 64)
	if err != nil || id64 == 0 {
		return 0, false
	}
	return uint(id64), true
}

func handleUserErr(c *gin.Context, err error) {
	switch {
	case err == nil:
		return
	case errors.Is(err, manager.ErrDBUnavailable):
		fail(c, http.StatusServiceUnavailable, err.Error())
	case errors.Is(err, manager.ErrInvalidCredentials), errors.Is(err, manager.ErrUserDisabled):
		fail(c, http.StatusUnauthorized, err.Error())
	case errors.Is(err, manager.ErrInvalidRole),
		errors.Is(err, manager.ErrWeakPassword),
		errors.Is(err, manager.ErrInvalidUsername):
		fail(c, http.StatusBadRequest, err.Error())
	case errors.Is(err, manager.ErrUsernameTaken):
		fail(c, http.StatusConflict, err.Error())
	case errors.Is(err, manager.ErrUserNotFound):
		fail(c, http.StatusNotFound, err.Error())
	default:
		fail(c, http.StatusInternalServerError, err.Error())
	}
}
