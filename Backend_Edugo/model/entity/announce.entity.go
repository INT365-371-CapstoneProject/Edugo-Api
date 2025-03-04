package entity

import "time"

type Announce_Post struct {
	Announce_ID     uint       `json:"id" gorm:"primaryKey"`
	Title           string     `json:"title"`
	Url             *string    `json:"url"`
	Description     string     `json:"description"`
	Attach_Name     *string    `json:"attach_name"`
	Attach_File     []byte     `gorm:"type:longblob" json:"attach_file"`
	Image           []byte     `gorm:"type:longblob" json:"image"`
	Publish_Date    *time.Time `json:"published_date"`
	Close_Date      *time.Time `json:"close_date"`
	Category_ID     uint       `json:"category_id"`
	Country_ID      uint       `json:"country_id"`
	Education_Level string     `json:"education_level"`
	Provider_ID     uint       `json:"provider_id"`

	// ตัวแปร Provider ใช้เพื่อเก็บข้อมูลจากตาราง Provider
	Provider Provider `gorm:"foreignKey:Provider_ID;references:Provider_ID"`
	// ตัวแปร Category ใช้เพื่อเก็บข้อมูลจากตาราง Category
	Category Category `gorm:"foreignKey:Category_ID;references:Category_ID"`
	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country Country `gorm:"foreignKey:Country_ID;references:Country_ID"`
}
