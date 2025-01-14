package request

import "time"

type AnnouncePostCreateRequest struct {
	Title        string     `json:"title" validate:"required,min=5,max=100"`
	Description  string     `json:"description" validate:"required,min=10,max=3000"`
	URL          *string    `json:"url" validate:"omitempty,min=10,max=255,url"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Posts_Type   string     `json:"posts_type" validate:"required,oneof=Announce"`
	Publish_Date *time.Time `json:"publish_date"`
	Close_Date   *time.Time `json:"close_date" validate:"required"`
	Category_ID  uint       `json:"category_id" validate:"required"`
	Country_ID   uint       `json:"country_id" validate:"required"`
}

type PostCreateRequest struct {
	Title        string     `json:"title" validate:"required,min=5,max=100"`
	Description  string     `json:"description" validate:"required,min=10,max=3000"`
	Image        *string    `json:"image"`
	Posts_Type   string     `json:"posts_type" validate:"required,oneof=Subject"`
	Publish_Date *time.Time `json:"publish_date"`
	Country_ID   uint       `json:"country_id" validate:"required"`
}

type PostUpdateRequest struct {
	Title        string     `json:"title" validate:"omitempty,min=5,max=100"`
	Description  string     `json:"description" validate:"omitempty,min=10,max=3000"`
	Image        *string    `json:"image"`
	Publish_Date *time.Time `json:"publish_date"`
	Country_ID   uint       `json:"country_id"`
}

type AnnouncePostUpdateRequest struct {
	Title        string     `json:"title" validate:"omitempty,min=5,max=100"`
	Description  string     `json:"description" validate:"omitempty,min=10,max=3000"`
	URL          *string    `json:"url" validate:"omitempty,min=10,max=255,url"`
	Attach_File  *string    `json:"attach_file"`
	Image        *string    `json:"image"`
	Publish_Date *time.Time `json:"publish_date"`
	Close_Date   *time.Time `json:"close_date"`
	Category_ID  uint       `json:"category_id"`
	Country_ID   uint       `json:"country_id"`
}
