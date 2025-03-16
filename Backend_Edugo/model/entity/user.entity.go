package entity

type Users struct {
	User_ID         uint    `json:"id" gorm:"primaryKey"`
	Education_level *string `json:"education_level"` // เปลี่ยนเป็น pointer เพื่อให้รองรับค่า nil
	Account_ID      uint    `json:"account_id"`

	// ตัวแปร Account ใช้เพื่อเก็บข้อมูลจากตาราง Account
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}
