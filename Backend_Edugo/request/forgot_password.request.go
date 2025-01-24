package request

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type VerifyOTPRequest struct {
	Email       string `json:"email" validate:"required,email"`
	OTPCode     string `json:"otp_code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
