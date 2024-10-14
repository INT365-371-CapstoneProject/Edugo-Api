package entity

type User struct {
	User_ID       uint   `json:"id" gorm:"primaryKey"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	AccountID int  `json:"account_account_id"`
}