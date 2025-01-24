package entity

import (
	"time"
)

type Comment struct {
	Comments_ID     uint      `json:"id" gorm:"primaryKey"`
	Comments_Text  string    `json:"comments_text"`
	Comments_Image []byte    `gorm:"type:longblob" json:"comments_image"`
	Publish_Date   time.Time `json:"publish_date" gorm:"autoCreateTime"`
	Posts_ID       uint      `json:"post_id"`    // Changed from Post_ID to Posts_ID
	Account_ID     uint      `json:"account_id"`

	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
	Post    Post    `gorm:"foreignKey:Posts_ID;references:Posts_ID"`
}