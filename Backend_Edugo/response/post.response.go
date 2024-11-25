package response

import "time"

type AnnouncePostResponse struct {
	Announce_ID    uint    `json:"id"`
	Title          string  `json:"title"`
	Description    string  `json:"description"`
	URL            *string `json:"url"`
	Attach_File    *string `json:"attach_file"`
	Image          *string `json:"image"`
	Posts_Type      string  `json:"posts_type"`
	Publish_Date *time.Time `json:"publish_date"`
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
	Posts_Type      string  `json:"posts_type"`
	Publish_Date *time.Time `json:"publish_date"`
	Close_Date    *time.Time  `json:"close_date"`
	Category_ID    uint    `json:"category_id"`
	Country_ID     uint    `json:"country_id"`
}

type PostResponse struct {
	Post_ID    uint    `json:"id"`
	Title      string  `json:"title"`
	Description string  `json:"description"`
	Image	  *string `json:"image"`
	Publish_Date *time.Time `json:"publish_date"`
	Posts_Type  string  `json:"posts_type"`
	Country     string `json:"country"`
}

type PostResponseAdd struct {
	Post_ID    uint    `json:"id"`
	Title      string  `json:"title"`
	Description string  `json:"description"`
	Image	  *string `json:"image"`
	Publish_Date *time.Time `json:"publish_date"`
	Posts_Type  string  `json:"posts_type"`
	Country_ID  uint    `json:"country_id"`
}