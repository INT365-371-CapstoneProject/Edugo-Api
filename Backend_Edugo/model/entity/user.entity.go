package entity

type User struct {
	User_ID       uint   `json:"id" gorm:"primaryKey"`
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	AccountAccountID int  `json:"account_account_id"`
}