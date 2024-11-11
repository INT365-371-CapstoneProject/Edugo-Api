package request

import "time"

type AnnouncePostCreateRequest struct {
	Title        string     `json:"title" validate:"required"`
	Description  string     `json:"description" validate:"required"`
	URL          *string     `json:"url"`
	Attach_File  *string     `json:"attach_file"`
	Image        *string     `json:"image"`
	Posts_Type   string     `json:"posts_type" validate:"required" enum:"Announce,Subject"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date   *time.Time `json:"closed_date"`
	Category_ID  uint       `json:"category_id" validate:"required"`
	Country_ID   uint       `json:"country_id" validate:"required"`
}


type AnnouncePostUpdateRequest struct {
	Title        string     `json:"title" validate:"required"`
	Description  string     `json:"description" validate:"required"`
	URL          *string     `json:"url"`
	Attach_File  *string    `json:"attach_file"`
	Image		*string     `json:"image"`
	Close_Date   *time.Time `json:"closed_date"`
}