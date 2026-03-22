package model

import "time"

// DecryptLog records every time-key decrypt attempt for auditing.
type DecryptLog struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	OperatorID    uint      `gorm:"index" json:"operatorId"`
	OperatorName  string    `gorm:"size:64;index" json:"operatorName"`
	OperatorRole  string    `gorm:"size:32;index" json:"operatorRole"`
	Operation     string    `gorm:"size:16;index" json:"operation"`
	DeviceSn      string    `gorm:"size:128;index" json:"deviceSn"`
	InputMode     string    `gorm:"size:16" json:"inputMode"`
	KeyPreview    string    `gorm:"size:64" json:"keyPreview"`
	KeyHash       string    `gorm:"size:64;index" json:"keyHash"`
	Success       bool      `gorm:"index" json:"success"`
	ErrorMessage  string    `gorm:"size:255" json:"errorMessage"`
	DecodedYear   int       `json:"decodedYear"`
	DecodedMonth  int       `json:"decodedMonth"`
	DecodedDay    int       `json:"decodedDay"`
	DecodedHour   int       `json:"decodedHour"`
	DecodedMinute int       `json:"decodedMinute"`
	DecodedSecond int       `json:"decodedSecond"`
	SourceIP      string    `gorm:"size:64" json:"sourceIp"`
	UserAgent     string    `gorm:"size:255" json:"userAgent"`
	CreatedAt     time.Time `json:"createdAt"`
}

func (DecryptLog) TableName() string {
	return "qa_decrypt_logs"
}
