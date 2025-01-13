package entity

type User struct {
	User_ID uint   `json:"id" gorm:"primaryKey"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Account_ID uint `json:"account_id"`

	// ตัวแปร Account ใช้เพื่อเก็บข้อมูลจากตาราง Account
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}