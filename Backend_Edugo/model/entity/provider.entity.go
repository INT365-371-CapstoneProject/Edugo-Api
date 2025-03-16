package entity

type Provider struct {
	Provider_ID  uint    `json:"id" gorm:"primaryKey"`
	Company_Name string  `json:"company_name"`
	URL          string  `json:"url"`
	Address      string  `json:"address"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
	Postal_Code  string  `json:"postal_code"`
	Phone        string  `json:"phone"`
	Phone_Person *string `json:"phone_person"`
	Verify       string  `json:"verify"`
	Account_ID   uint    `json:"account_id"`

	// ตัวแปร Account ใช้เพื่อเก็บข้อมูลจากตาราง Account
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}
