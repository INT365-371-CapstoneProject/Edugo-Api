package entity

type Category struct {
	Category_ID uint   `json:"id" gorm:"primaryKey"`
	Name		string `json:"category_name"`
}