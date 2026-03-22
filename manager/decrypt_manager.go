package manager

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"qa_test_server/db"
	"qa_test_server/model"
)

var (
	ErrDecryptDBUnavailable = errors.New("database is unavailable")
	ErrDecryptSNRequired    = errors.New("device sn is required")
	ErrDecryptKeyRequired   = errors.New("decrypt key is required")
	ErrDecryptKeyLength     = errors.New("decrypt key length must be 28")
	ErrDecryptSNLength      = errors.New("device sn length must be 1-20")
	ErrDecryptKeyInvalid    = errors.New("decrypt key verification failed")
	ErrDecryptTimeRequired  = errors.New("time is required")
)

var timeKeySecureMod = []byte("www.ultronphotonics.com")

type TimeKeyDecodeResult struct {
	DeviceSn      string    `json:"deviceSn"`
	Valid         bool      `json:"valid"`
	DecodedYear   int       `json:"decodedYear"`
	DecodedMonth  int       `json:"decodedMonth"`
	DecodedDay    int       `json:"decodedDay"`
	DecodedHour   int       `json:"decodedHour"`
	DecodedMinute int       `json:"decodedMinute"`
	DecodedSecond int       `json:"decodedSecond"`
	FullYear      int       `json:"fullYear"`
	DecodedAt     time.Time `json:"decodedAt"`
	DecodedAtText string    `json:"decodedAtText"`
	RawHead       string    `json:"rawHead"`
	RawTail       string    `json:"rawTail"`
}

type TimeKeyGenerateResult struct {
	DeviceSn      string    `json:"deviceSn"`
	InputTime     time.Time `json:"inputTime"`
	InputTimeText string    `json:"inputTimeText"`
	DecodedYear   int       `json:"decodedYear"`
	DecodedMonth  int       `json:"decodedMonth"`
	DecodedDay    int       `json:"decodedDay"`
	DecodedHour   int       `json:"decodedHour"`
	DecodedMinute int       `json:"decodedMinute"`
	DecodedSecond int       `json:"decodedSecond"`
	FullYear      int       `json:"fullYear"`
	Key           string    `json:"key"`
	KeyHex        string    `json:"keyHex"`
	RawHead       string    `json:"rawHead"`
	RawTail       string    `json:"rawTail"`
}

type TimeKeyDecodeLogInput struct {
	OperatorID    uint
	OperatorName  string
	OperatorRole  string
	Operation     string
	DeviceSn      string
	InputMode     string
	KeyRaw        string
	SourceIP      string
	UserAgent     string
	Result        *TimeKeyDecodeResult
	DecodedYear   int
	DecodedMonth  int
	DecodedDay    int
	DecodedHour   int
	DecodedMinute int
	DecodedSecond int
	Err           error
}

type DecryptLogQuery struct {
	DeviceSn  string
	Keyword   string
	Operation string
	Offset    int
	Limit     int
}

type DecryptManager struct{}

var DecryptManagerGlobal = &DecryptManager{}

func (m *DecryptManager) Init() error {
	if db.DB == nil {
		return ErrDecryptDBUnavailable
	}
	return db.DB.AutoMigrate(&model.DecryptLog{})
}

func (m *DecryptManager) DecodeTimeKey(sn, keyRaw string) (TimeKeyDecodeResult, error) {
	deviceSN := strings.TrimSpace(sn)
	if deviceSN == "" {
		return TimeKeyDecodeResult{}, ErrDecryptSNRequired
	}
	if len(deviceSN) < 1 || len(deviceSN) > 20 {
		return TimeKeyDecodeResult{}, ErrDecryptSNLength
	}

	keyBytes := []byte(strings.TrimSpace(keyRaw))
	if len(keyBytes) == 0 {
		return TimeKeyDecodeResult{}, ErrDecryptKeyRequired
	}
	if len(keyBytes) != 28 {
		return TimeKeyDecodeResult{}, ErrDecryptKeyLength
	}

	head := make([]byte, 14)
	tail := make([]byte, 14)
	copy(head, keyBytes[:14])
	copy(tail, keyBytes[14:28])

	xorPayload := make([]byte, 14)
	for i := 0; i < 14; i++ {
		xorPayload[i] = head[i] ^ deviceSN[i%len(deviceSN)]
	}
	expectedTail, err := calcUnlockKey(xorPayload, 0, 14)
	if err != nil {
		return TimeKeyDecodeResult{}, err
	}
	if string(expectedTail) != string(tail) {
		return TimeKeyDecodeResult{}, ErrDecryptKeyInvalid
	}

	year, ok := decodePair(head[13], head[12])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid year", ErrDecryptKeyInvalid)
	}
	month, ok := decodePair(head[11], head[10])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid month", ErrDecryptKeyInvalid)
	}
	day, ok := decodePair(head[9], head[8])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid day", ErrDecryptKeyInvalid)
	}
	hour, ok := decodePair(head[7], head[6])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid hour", ErrDecryptKeyInvalid)
	}
	minute, ok := decodePair(head[5], head[4])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid minute", ErrDecryptKeyInvalid)
	}
	second, ok := decodePair(head[3], head[2])
	if !ok {
		return TimeKeyDecodeResult{}, fmt.Errorf("%w: invalid second", ErrDecryptKeyInvalid)
	}

	fullYear := 2000 + year
	decodedAt, validDatetime := buildDecodedTime(fullYear, month, day, hour, minute, second)
	decodedText := ""
	if validDatetime {
		decodedText = decodedAt.Format("2006-01-02 15:04:05")
	}

	return TimeKeyDecodeResult{
		DeviceSn:      deviceSN,
		Valid:         true,
		DecodedYear:   year,
		DecodedMonth:  month,
		DecodedDay:    day,
		DecodedHour:   hour,
		DecodedMinute: minute,
		DecodedSecond: second,
		FullYear:      fullYear,
		DecodedAt:     decodedAt,
		DecodedAtText: decodedText,
		RawHead:       string(head),
		RawTail:       string(tail),
	}, nil
}

func (m *DecryptManager) GenerateTimeKey(sn string, t time.Time) (TimeKeyGenerateResult, error) {
	deviceSN := strings.TrimSpace(sn)
	if deviceSN == "" {
		return TimeKeyGenerateResult{}, ErrDecryptSNRequired
	}
	if len(deviceSN) < 1 || len(deviceSN) > 20 {
		return TimeKeyGenerateResult{}, ErrDecryptSNLength
	}
	if t.IsZero() {
		return TimeKeyGenerateResult{}, ErrDecryptTimeRequired
	}

	local := t.In(time.Local)
	year := local.Year() % 100
	month := int(local.Month())
	day := local.Day()
	hour := local.Hour()
	minute := local.Minute()
	second := local.Second()

	head := make([]byte, 14)
	for i := 0; i < 14; i++ {
		head[i] = 'A'
	}

	if !setTimePair(head, 2, second) ||
		!setTimePair(head, 4, minute) ||
		!setTimePair(head, 6, hour) ||
		!setTimePair(head, 8, day) ||
		!setTimePair(head, 10, month) ||
		!setTimePair(head, 12, year) {
		return TimeKeyGenerateResult{}, ErrDecryptKeyInvalid
	}

	xorPayload := make([]byte, 14)
	for i := 0; i < 14; i++ {
		xorPayload[i] = head[i] ^ deviceSN[i%len(deviceSN)]
	}
	tail, err := calcUnlockKey(xorPayload, 0, 14)
	if err != nil {
		return TimeKeyGenerateResult{}, err
	}

	keyBytes := append(append([]byte{}, head...), tail...)
	return TimeKeyGenerateResult{
		DeviceSn:      deviceSN,
		InputTime:     local,
		InputTimeText: local.Format("2006-01-02 15:04:05"),
		DecodedYear:   year,
		DecodedMonth:  month,
		DecodedDay:    day,
		DecodedHour:   hour,
		DecodedMinute: minute,
		DecodedSecond: second,
		FullYear:      local.Year(),
		Key:           string(keyBytes),
		KeyHex:        strings.ToUpper(hex.EncodeToString(keyBytes)),
		RawHead:       string(head),
		RawTail:       string(tail),
	}, nil
}

func (m *DecryptManager) WriteDecodeLog(ctx context.Context, input TimeKeyDecodeLogInput) error {
	if db.DB == nil {
		return ErrDecryptDBUnavailable
	}

	sn := strings.TrimSpace(input.DeviceSn)
	if len(sn) > 128 {
		sn = sn[:128]
	}
	mode := strings.TrimSpace(strings.ToLower(input.InputMode))
	if mode == "" {
		mode = "plain"
	}
	operation := strings.TrimSpace(strings.ToLower(input.Operation))
	if operation == "" {
		operation = "decode"
	}

	keyRaw := input.KeyRaw
	keyHash := sha256.Sum256([]byte(keyRaw))
	keyPreview := buildKeyPreview(keyRaw)

	entry := model.DecryptLog{
		OperatorID:   input.OperatorID,
		OperatorName: truncate(input.OperatorName, 64),
		OperatorRole: truncate(input.OperatorRole, 32),
		Operation:    truncate(operation, 16),
		DeviceSn:     sn,
		InputMode:    truncate(mode, 16),
		KeyPreview:   truncate(keyPreview, 64),
		KeyHash:      hex.EncodeToString(keyHash[:]),
		Success:      input.Err == nil,
		SourceIP:     truncate(strings.TrimSpace(input.SourceIP), 64),
		UserAgent:    truncate(strings.TrimSpace(input.UserAgent), 255),
	}

	if input.Err != nil {
		entry.ErrorMessage = truncate(input.Err.Error(), 255)
	} else {
		entry.DecodedYear = input.DecodedYear
		entry.DecodedMonth = input.DecodedMonth
		entry.DecodedDay = input.DecodedDay
		entry.DecodedHour = input.DecodedHour
		entry.DecodedMinute = input.DecodedMinute
		entry.DecodedSecond = input.DecodedSecond

		if input.Result != nil {
			entry.DecodedYear = input.Result.DecodedYear
			entry.DecodedMonth = input.Result.DecodedMonth
			entry.DecodedDay = input.Result.DecodedDay
			entry.DecodedHour = input.Result.DecodedHour
			entry.DecodedMinute = input.Result.DecodedMinute
			entry.DecodedSecond = input.Result.DecodedSecond
		}
	}

	return db.DB.WithContext(ctx).Create(&entry).Error
}

func (m *DecryptManager) QueryLogs(q DecryptLogQuery) ([]model.DecryptLog, int64, error) {
	if db.DB == nil {
		return nil, 0, ErrDecryptDBUnavailable
	}
	offset, limit := normalizeDecodeLogRange(q.Offset, q.Limit)
	sn := strings.TrimSpace(q.DeviceSn)
	keyword := strings.TrimSpace(q.Keyword)
	operation := strings.TrimSpace(strings.ToLower(q.Operation))

	query := db.DB.Model(&model.DecryptLog{})
	if sn != "" {
		query = query.Where("device_sn = ?", sn)
	}
	if operation != "" {
		query = query.Where("operation = ?", operation)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		query = query.Where("device_sn LIKE ? OR operator_name LIKE ?", like, like)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	items := make([]model.DecryptLog, 0, limit)
	if err := query.Order("id DESC").Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func calcUnlockKey(payload []byte, unlockTime uint16, length int) ([]byte, error) {
	if length <= 0 || length > 19 || len(payload) < length {
		return nil, errors.New("invalid payload length")
	}
	out := make([]byte, length)

	sum := byte(0)
	for i := 0; i < length; i++ {
		sum += payload[i]
	}

	for i := 0; i < length; i++ {
		mask := byte((int(timeKeySecureMod[i]) + int(unlockTime) + int(sum)) & 0xff)
		v := payload[i] ^ mask
		out[i] = 'A' + (v % 26)
	}
	return out, nil
}

func setTimePair(head []byte, lowIndex, value int) bool {
	if lowIndex < 0 || lowIndex+1 >= len(head) {
		return false
	}
	if value < 0 || value > 99 {
		return false
	}
	ones := value % 10
	tens := value / 10
	head[lowIndex] = byte(int('A') + ones)
	head[lowIndex+1] = byte(int('A') + tens)
	return true
}

func decodePair(highRaw, lowRaw byte) (int, bool) {
	high := int(highRaw) - int('A')
	low := int(lowRaw) - int('A')
	if high < 0 || high > 9 || low < 0 || low > 9 {
		return 0, false
	}
	return high*10 + low, true
}

func buildDecodedTime(year, month, day, hour, minute, second int) (time.Time, bool) {
	if month < 1 || month > 12 || day < 1 || day > 31 || hour < 0 || hour > 23 || minute < 0 || minute > 59 || second < 0 || second > 59 {
		return time.Time{}, false
	}
	ts := time.Date(year, time.Month(month), day, hour, minute, second, 0, time.Local)
	if ts.Year() != year || int(ts.Month()) != month || ts.Day() != day || ts.Hour() != hour || ts.Minute() != minute || ts.Second() != second {
		return time.Time{}, false
	}
	return ts, true
}

func normalizeDecodeLogRange(offset, limit int) (int, int) {
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

func buildKeyPreview(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return value
	}
	return value[:4] + "****" + value[len(value)-4:]
}

func truncate(value string, max int) string {
	if max <= 0 {
		return ""
	}
	if len(value) <= max {
		return value
	}
	return value[:max]
}
