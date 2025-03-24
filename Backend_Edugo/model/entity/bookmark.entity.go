package entity

import (
	"time"
)

type Bookmark struct {
	Bookmark_ID uint      `gorm:"primaryKey;autoIncrement" json:"bookmark_id"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	Account_ID  uint      `gorm:"not null" json:"account_id"`
	Announce_ID uint      `gorm:"not null" json:"announce_id"`

	Account  Account       `gorm:"foreignKey:Account_ID;references:Account_ID"`
	Announce Announce_Post `gorm:"foreignKey:Announce_ID;references:Announce_ID"`
}
