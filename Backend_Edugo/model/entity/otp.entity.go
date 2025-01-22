package entity

import "time"

type OTP struct {
	ID           uint      `gorm:"primaryKey;column:otp_id" json:"id"`
	Code         string    `gorm:"size:6;not null" json:"code"`
	IsUsed       bool      `gorm:"default:false" json:"is_used"`
	AttemptCount int       `gorm:"default:0" json:"attempt_count"`
	ExpiredAt    time.Time `gorm:"not null" json:"expired_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	AccountID    uint      `gorm:"column:account_id;not null" json:"account_id"`
	Account      Account   `gorm:"foreignKey:AccountID;references:Account_ID"`
}
