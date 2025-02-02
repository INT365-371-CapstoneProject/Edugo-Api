package response

import "time"

type AnnouncePostResponse struct {
	Announce_ID     uint       `json:"id"`
	Title           string     `json:"title"`
	Attach_Name     *string    `json:"attach_name"`
	Description     string     `json:"description"`
	URL             *string    `json:"url"`
	Posts_Type      string     `json:"posts_type"`
	Publish_Date    *time.Time `json:"publish_date"`
	Close_Date      *time.Time `json:"close_date"`
	Education_Level string     `json:"education_level"`

	// ตัวแปร Category ใช้เพื่อเก็บข้อมูลจากตาราง Category
	Category string `json:"category"`
	// ตัวแปร Country ใช้เพื่อเก็บข้อมูลจากตาราง Country
	Country string `json:"country"`
	Post_ID uint   `json:"post_id"`
}

type AnnouncePostResponseAdd struct {
	Announce_ID     uint       `json:"id"`
	Title           string     `json:"title"`
	Description     string     `json:"description"`
	URL             *string    `json:"url"`
	Attach_Name     *string    `json:"attach_name"`
	Posts_Type      string     `json:"posts_type"`
	Publish_Date    *time.Time `json:"publish_date"`
	Close_Date      *time.Time `json:"close_date"`
	Category_ID     uint       `json:"category_id"`
	Country_ID      uint       `json:"country_id"`
	Account_ID      uint       `json:"account_id"`
	Education_Level string     `json:"education_level"`
}

type PostResponse struct {
	Post_ID      uint       `json:"id"`
	Description  string     `json:"description"`
	Publish_Date *time.Time `json:"publish_date"`
	Posts_Type   string     `json:"posts_type"`
	Account_ID   uint       `json:"account_id"`
}

type PaginatedPostResponse struct {
	Data     []PostResponse `json:"data"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	LastPage int            `json:"last_page"`
	PerPage  int            `json:"per_page"`
}

type PaginatedAnnouncePostResponse struct {
	Data     []AnnouncePostResponse `json:"data"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	LastPage int                    `json:"last_page"`
	PerPage  int                    `json:"per_page"`
}

type PostResponseAdd struct {
	Post_ID      uint       `json:"id"`
	Description  string     `json:"description"`
	Publish_Date *time.Time `json:"publish_date"`
	Posts_Type   string     `json:"posts_type"`
	Account_ID   uint       `json:"account_id"`
}
