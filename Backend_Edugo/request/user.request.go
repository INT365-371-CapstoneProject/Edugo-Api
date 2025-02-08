package request

type UserRequest struct {
    Email           string `json:"email" validate:"required,email"`
    Username        string `json:"username" validate:"required"`
    Password        string `json:"password" validate:"required"`
    FirstName       string `json:"first_name" validate:"required"`
    LastName        string `json:"last_name" validate:"required"`
    ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}