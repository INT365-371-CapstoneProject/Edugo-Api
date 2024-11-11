package entity

import "time"

type Announce_Post struct {
	Announce_ID uint       `json:"id" gorm:"primaryKey"`
	Url         *string    `json:"url"`
	Attach_File *string    `json:"attach_file"`
	Close_Date *time.Time `json:"close_date"`
	Posts_ID     uint       `json:"post_id"`
	Category_ID uint       `json:"category_id"`

	// ตัวแปร Post ใช้เพื่อเก็บข้อมูลจากตาราง Post
	Post 	  Post       `gorm:"foreignKey:Posts_ID"`
	// ตัวแปร Category ใช้เพื่อเก็บข้อมูลจากตาราง Category
	Category  Category   `gorm:"foreignKey:Category_ID"`
}