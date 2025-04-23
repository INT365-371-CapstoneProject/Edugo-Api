package response

import (
	"time"
)

type NotificationResponse struct {
	Notification_ID uint      `json:"id"`
	Title           string    `json:"title"`
	Message         string    `json:"message"`
	Is_Read         uint      `json:"is_read"`
	Created_At      time.Time `json:"created_at"`
	Announce_ID     uint      `json:"announce_id"`
	Account_ID      uint      `json:"account_id"`
}

type NotificationUpdateResponse struct {
	Notification_ID uint `json:"id"`
	Is_Read         uint `json:"is_read"`
	Announce_ID     uint `json:"announce_id"`
	Account_ID      uint `json:"account_id"`
}
