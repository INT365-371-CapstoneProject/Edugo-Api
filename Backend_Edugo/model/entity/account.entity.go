package entity

import (
	"time"
)

type Account struct {
	Account_ID uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Create_On  time.Time `json:"create_on" gorm:"autoCreateTime"`
	Last_Login *time.Time `json:"last_login"`
	Update_On  time.Time `json:"update_on" gorm:"autoUpdateTime"`
	Role       string    `json:"role"`
}
