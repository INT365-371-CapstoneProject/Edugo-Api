package request

type AdminCreateRequest struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=5,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
	FirstName       string `json:"first_name" validate:"required,min=5,max=100"`
	LastName        string `json:"last_name" validate:"required,min=5,max=100"`
	Phone           string `json:"phone" validate:"required,min=5,max=100"`
}
