package response

import (
	"time"
)

type BookmarkResponse struct {
	Bookmark_ID uint      `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	Announce_ID uint      `json:"annouce_id"`
	Account_ID  uint      `json:"account_id"`
}
