package request

import (
	"mime/multipart"
	"time"
)

type CreateCommentRequest struct {
	Comments_Text  string          `form:"comments_text" validate:"required"`
	Comments_Image *multipart.File `form:"comments_image"` // เปลี่ยนเป็น *multipart.File
	Posts_ID       uint            `form:"posts_id" validate:"required"`
	// Account_ID    uint           `form:"account_id" validate:"required"`
}

type UpdateCommentRequest struct {
	Comments_Text string     `form:"comments_text" validate:"required"`
	Publish_Date  *time.Time `json:"publish_date"`
}
