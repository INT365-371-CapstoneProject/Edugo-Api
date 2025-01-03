package entity

import "time"

type Post struct {
	Posts_ID        uint   `json:"id" gorm:"primaryKey"`
	Description    string `json:"description"`
	Image 		*string `json:"image"`
	Posts_Type      string `json:"post_type"`
	Publish_Date *time.Time `json:"published_date"`
}