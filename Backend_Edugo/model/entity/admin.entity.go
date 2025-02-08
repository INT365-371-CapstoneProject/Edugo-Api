package entity

type Admin struct {
	Admin_ID uint `json:"id" gorm:"primaryKey"`
	Phone   string `json:"phone_number"`
	Status  string `json:"status"`
	Account_ID uint `json:"account_id"`
	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}