package entity

import "time"

type FCMToken struct {
	Token_ID   uint      `gorm:"primaryKey;autoIncrement" json:"token_id"`
	Account_ID uint      `gorm:"not null;index" json:"account_id"`
	FCM_Token  string    `gorm:"size:255;not null" json:"fcm_token"` // แก้พิมพ์ผิดจาก varcher เป็น size
	Created_At time.Time `gorm:"autoCreateTime" json:"created_at"`

	// เชื่อมกับตาราง Account และให้ลบอัตโนมัติเมื่อ Account ถูกลบ
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID;constraint:OnDelete:CASCADE;"`
}
