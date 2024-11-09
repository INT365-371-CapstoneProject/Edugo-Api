package entity

import "time"

type Post struct {
	Posts_ID        uint   `json:"id" gorm:"primaryKey"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	URL            *string `json:"url"`
	Attach_File    *string `json:"attach_file"`
	Image 		*string `json:"image"`
	Posts_Type      string `json:"post_type"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date    *time.Time `json:"closed_date"`
	Provider_ID	*uint   `json:"provider_id"`
	User_ID       *uint   `json:"user_id"`
}