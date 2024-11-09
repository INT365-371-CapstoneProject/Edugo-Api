package request

import "time"

type PostCreateRequest struct {
	Title        string     `json:"title" validate:"required"`
	Description  string     `json:"description" validate:"required"`
	URL          *string     `json:"url"`
	Attach_File  *string     `json:"attach_file"`
	Image        *string     `json:"image"`
	Posts_Type   string     `json:"post_type" validate:"required" enum:"Annouce,Subject"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date   *time.Time `json:"closed_date"`
	Provider_ID  *uint        `json:"provider_id"`
	User_ID      *uint        `json:"user_id"`
}


type PostUpdateRequest struct {
	Title        string     `json:"title" validate:"required"`
	Description  string     `json:"description" validate:"required"`
	URL          *string     `json:"url"`
	Attach_File  *string    `json:"attach_file"`
	Image		*string     `json:"image"`
	Close_Date   *time.Time `json:"closed_date"`
}