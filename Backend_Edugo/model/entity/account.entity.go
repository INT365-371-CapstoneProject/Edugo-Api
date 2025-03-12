package entity

import (
	"time"
)

type Account struct {
	Account_ID uint       `json:"id" gorm:"primaryKey"`
	Username   string     `json:"username"`
	Password   string     `json:"password"`
	Email      string     `json:"email"`
	FirstName  *string    `json:"first_name" gorm:"column:first_name"` // เปลี่ยนเป็น first_name
	LastName   *string    `json:"last_name" gorm:"column:last_name"`   // เปลี่ยนเป็น last_name
	Avatar     []byte     `gorm:"type:longblob" json:"avatar"`
	Status     string     `json:"status"`
	Create_On  time.Time  `json:"create_on" gorm:"autoCreateTime"`
	Last_Login *time.Time `json:"last_login"`
	Update_On  time.Time  `json:"update_on" gorm:"autoUpdateTime"`
	Role       string     `json:"role"`
}
