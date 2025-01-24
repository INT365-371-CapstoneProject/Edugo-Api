package request


type ProviderCreateRequest struct {
	Company_name string `json:"company_name" validate:"required,min=5,max=100"`
	URL string `json:"url" validate:"required,min=5,max=100"`
	Address string `json:"address" validate:"required,min=5,max=100"`
	Phone string `json:"phone" validate:"required,min=5,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}