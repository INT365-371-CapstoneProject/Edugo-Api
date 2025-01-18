package entity


type Comment struct {
	Comment_ID uint `json:"id" gorm:"primaryKey"`
	Comments_Text string `json:"comments_text"`
	Comments_Image *string `json:"comments_image"`
	Comments_Type string `json:"comments_type"`
	Publish_Date string `json:"publish_date" gorm:"autoCreateTime"`
	Post_ID uint `json:"post_id"`
	Account_ID uint `json:"account_id"`

	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
	Post Post `gorm:"foreignKey:Post_ID;references:Post_ID"`
}