package response

import (
	"time"
)

type FCMTokenResponse struct {
	Token_ID   uint      `json:"id"`
	Account_ID uint      `json:"account_id"`
	FCM_Token  string    `json:"fcm_token"`
	Created_At time.Time `json:"created_at"`
}
