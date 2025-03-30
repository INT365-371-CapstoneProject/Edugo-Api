package entity

import (
	"time"
)

type Notification struct {
	Notification_ID uint      `gorm:"primaryKey;autoIncrement" json:"notification_id"`
	Title           string    `gorm:"size:100;not null" json:"title"`
	Message         string    `gorm:"size:500;not null" json:"message"`
	Is_Read         uint      `gorm:"type:TINYINT(1);not null" json:"is_read"`
	Created_At      time.Time `gorm:"not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	Account_ID      uint      `gorm:"not null" json:"account_id"`
	Announce_ID     uint      `gorm:"not null" json:"announce_id"`

	Account       Account       `gorm:"foreignKey:Account_ID;references:Account_ID"`
	Announce_Post Announce_Post `gorm:"foreignKey:Announce_ID;references:Announce_ID"`
}
