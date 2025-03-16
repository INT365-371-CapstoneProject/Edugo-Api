package entity

type AnswerCountries struct {
	Answer_ID uint `json:"id" gorm:"primaryKey"`
	User_ID uint `json:"user_id"`
	Country_ID uint `json:"country_id"`

	// ตัวแปร User ใช้เพื่อเก็บข้อมูลจากตาราง User
	User Users `gorm:"foreignKey:User_ID;references:User_ID"`

	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country Country `gorm:"foreignKey:Country_ID;references:Country_ID"`
}