package entity

import "time"

type Announce_Post struct {
	Announce_ID uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Url         *string    `json:"url"`
	Attach_Name *string    `json:"attach_name"`
	Attach_File []byte     `gorm:"type:longblob" json:"attach_file"`
	Close_Date *time.Time `json:"close_date"`
	Posts_ID     uint       `json:"post_id"`
	Category_ID uint       `json:"category_id"`
	Country_ID  uint       `json:"country_id"`

	// ตัวแปร Post ใช้เพื่อเก็บข้อมูลจากตาราง Post
	Post 	  Post       `gorm:"foreignKey:Posts_ID;references:Posts_ID "`
	// ตัวแปร Category ใช้เพื่อเก็บข้อมูลจากตาราง Category
	Category  Category   `gorm:"foreignKey:Category_ID;references:Category_ID"`
	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country   Country    `gorm:"foreignKey:Country_ID;references:Country_ID"`
}