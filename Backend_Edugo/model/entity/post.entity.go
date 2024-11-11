package entity

import "time"

type Post struct {
	Posts_ID        uint   `json:"id" gorm:"primaryKey"`
	Title          string `json:"title"`
	Description    string `json:"description"`
	Image 		*string `json:"image"`
	Posts_Type      string `json:"post_type"`
	Publish_Date *time.Time `json:"published_date"`
	Country_ID     uint   `json:"country_id"`

	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country Country `gorm:"foreignKey:Country_ID"`
}