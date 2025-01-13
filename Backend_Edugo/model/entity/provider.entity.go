package entity

type Provider struct {
	Provider_ID uint   `json:"id" gorm:"primaryKey"`
	Provider_Name string `json:"provider_name"`
	URL string `json:"url"`
	Address string `json:"address"`
	Status string `json:"status"`
	Verify string `json:"verify"`
	Account_ID uint `json:"account_id"`

	// ตัวแปร Account ใช้เพื่อเก็บข้อมูลจากตาราง Account
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}