package request

type CreateFCMTokenRequest struct {
	FCM_Token string `form:"fcm_token" validate:"required"`
}
