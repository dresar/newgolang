package models

import (
	"time"
)

type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Username     string    `gorm:"uniqueIndex;size:64" json:"username"`
	PasswordHash string    `json:"-"`
	APIToken     string    `gorm:"size:128" json:"api_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Device struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:128" json:"name"`
	Status       string    `gorm:"size:64" json:"status"`
	LastQRCode   string    `gorm:"type:text" json:"last_qr_code"`
	Connected    bool      `json:"connected"`
	PhoneNumber  string    `gorm:"size:32" json:"phone_number"`
	LastSeenAt   time.Time `json:"last_seen_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AutoReply struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Keyword   string    `gorm:"size:128;index" json:"keyword"`
	ReplyText string    `gorm:"type:text" json:"reply_text"`
	MatchType string    `gorm:"size:16" json:"match_type"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type MessageLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	To        string    `gorm:"size:64;index" json:"to"`
	Message   string    `gorm:"type:text" json:"message"`
	Status    string    `gorm:"size:32" json:"status"`
	Error     string    `gorm:"type:text" json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

