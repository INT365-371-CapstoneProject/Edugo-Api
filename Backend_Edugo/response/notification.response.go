package response

import (
	"time"
)

type NotificationResponse struct {
	Notification_ID uint      `json:"id"`
	Title           string    `json:"title"`
	Message         string    `json:"message"`
	IsRead          uint      `json:"is_read"`
	CreatedAt       time.Time `json:"created_at"`
	Announce_ID     uint      `json:"announce_id"`
	Account_ID      uint      `json:"account_id"`
}
