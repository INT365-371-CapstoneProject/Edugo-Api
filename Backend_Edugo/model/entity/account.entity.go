package entity

type Account struct {
	Account_ID uint   `json:"id" gorm:"primaryKey"`
	Phone_Number string `json:"phone_number"`
	Create_On string `json:"create_on"`
	Update_On string `json:"update_on"`
	Last_Login string `json:"last_login"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email string `json:"email"`
	Role string `json:"role"`
}