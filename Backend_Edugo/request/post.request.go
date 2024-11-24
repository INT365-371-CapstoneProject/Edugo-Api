package request

import "time"

type AnnouncePostCreateRequest struct {
	Title        string     `json:"title" validate:"required,min=5,max=100" message:"Title is required and must be between 5-100 characters"`
	Description  string     `json:"description" validate:"required,min=10,max=500" message:"Description is required and must be between 10-500 characters"`
	URL          *string    `json:"url" validate:"omitempty,min=10,max=255,url" message:"URL must be between 10-255 characters"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Posts_Type   string     `json:"posts_type" validate:"required,oneof=Announce Subject" message:"Posts_Type is required and must be either 'Announce' or 'Subject'"`
	Publish_Date *time.Time `json:"publish_date"` 
	Close_Date   *time.Time `json:"close_date" validate:"required" message:"Close_Date is required"`
	Category_ID  uint       `json:"category_id" validate:"required" message:"Category_ID is required"`
	Country_ID   uint       `json:"country_id" validate:"required" message:"Country_ID is required"`
}


type AnnouncePostUpdateRequest struct {
	Title        string     `json:"title" validate:"omitempty,min=5,max=100" message:"Title is must be between 5-100 characters"`
	Description  string     `json:"description" validate:"omitempty,min=10,max=500" message:"Description is must be between 10-500 characters"`
	URL          *string    `json:"url" validate:"omitempty,min=10,max=255,url" message:"URL must be between 10-255 characters"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Publish_Date *time.Time `json:"publish_date"`
	Close_Date   *time.Time `json:"close_date"`
	Category_ID  uint       `json:"category_id"`
	Country_ID   uint       `json:"country_id"`
}
