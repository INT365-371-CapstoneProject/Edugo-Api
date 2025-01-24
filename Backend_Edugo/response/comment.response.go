package response

import (
	"time"
)

type CommentResponse struct {
	Comments_ID     uint      `json:"id"`
	Comments_Text  string    `json:"comments_text"`
	Publish_Date   time.Time `json:"publish_date"`
	Posts_ID        uint      `json:"post_id"`
	Account_ID     uint      `json:"account_id"`
}