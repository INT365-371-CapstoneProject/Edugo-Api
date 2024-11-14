package response

import "time"

type AnnouncePostResponse struct {
	Announce_ID    uint    `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	URL            *string `json:"url"`
	Attach_File    *string `json:"attach_file"`
	Image          *string `json:"image"`
	Post_Type      string  `json:"post_type"`
	Published_Date *time.Time `json:"published_date"`
	Close_Date    *time.Time  `json:"close_date"`

	// ตัวแปร Category ใช้เพื่อเก็บข้อมูลจากตาราง Category
	Category    string `json:"category"`
	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country     string `json:"country"`
}

type AnnouncePostResponseAdd struct {
	Announce_ID    uint    `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	URL            *string `json:"url"`
	Attach_File    *string `json:"attach_file"`
	Image          *string `json:"image"`
	Post_Type      string  `json:"post_type"`
	Published_Date *time.Time `json:"published_date"`
	Close_Date    *time.Time  `json:"close_date"`
	Category_ID    uint    `json:"category_id"`
	Country_ID     uint    `json:"country_id"`
}