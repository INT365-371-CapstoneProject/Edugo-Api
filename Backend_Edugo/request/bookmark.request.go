package request

import "mime/multipart"

type CreateBookmarkRequest struct {
	Comments_Text  string          `form:"comments_text" validate:"required"`
	Comments_Image *multipart.File `form:"comments_image"` // เปลี่ยนเป็น *multipart.File
	Posts_ID       uint            `form:"posts_id" validate:"required"`
	// Account_ID    uint           `form:"account_id" validate:"required"`
}
