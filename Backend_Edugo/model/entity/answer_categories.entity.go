package entity

type AnswerCategories struct {
    Answer_ID   uint     `json:"id" gorm:"primaryKey"`
    User_ID     uint     `json:"user_id"`
    Category_ID uint     `json:"category_id"`

    // Relations
    User     Users     `gorm:"foreignKey:User_ID;references:User_ID"`
    Category Category  `gorm:"foreignKey:Category_ID;references:Category_ID"`
}
