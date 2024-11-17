package request

import "time"

type AnnouncePostCreateRequest struct {
	Title        string     `json:"title" validate:"required" message:"Title is required"`
	Description  string     `json:"description" validate:"required" message:"Description is required"`
	URL          *string    `json:"url"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Posts_Type   string     `json:"posts_type" validate:"required,oneof=Announce Subject" message:"Posts_Type is required and must be either 'Announce' or 'Subject'"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date   *time.Time `json:"closed_date" validate:"required" message:"Close_Date is required"`
	Category_ID  uint       `json:"category_id" validate:"required" message:"Category_ID is required"`
	Country_ID   uint       `json:"country_id" validate:"required" message:"Country_ID is required"`
}

type AnnouncePostUpdateRequest struct {
	Title        string     `json:"title" validate:"required" message:"Title is required"`
	Description  string     `json:"description" validate:"required" message:"Description is required"`
	URL          *string    `json:"url"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Publish_Date *time.Time `json:"published_date"`
	Close_Date   *time.Time `json:"closed_date" validate:"required" message:"Close_Date is required"`
	Category_ID  uint       `json:"category_id" validate:"required" message:"Category_ID is required"`
	Country_ID   uint       `json:"country_id" validate:"required" message:"Country_ID is required"`
}
