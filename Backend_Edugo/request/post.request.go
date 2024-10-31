package request

import "time"

type PostCreateRequest struct {
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	URL          *string     `json:"url"`
	Attach_File  *[]byte     `json:"attach_file"`
	Posts_Type   string     `json:"post_type"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date   *time.Time `json:"closed_date"`
	Provider_ID  *uint        `json:"provider_id"`
	User_ID      *uint        `json:"user_id"`
}