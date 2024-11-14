package entity

type Country struct {
	Country_ID uint   `json:"id" gorm:"primaryKey"`
	Name		string `json:"country_name"`
}