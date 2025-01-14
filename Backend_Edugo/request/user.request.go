package request

type UserRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Username        string `json:"username" validate:"required"`
    Password        string `json:"password" validate:"required"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}