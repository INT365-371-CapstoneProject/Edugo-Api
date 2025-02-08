package request


type ProviderCreateRequest struct {
	Company_name string `json:"company_name" validate:"required,min=5,max=100"`
	FirstName string `json:"first_name" validate:"required,min=5,max=100"`
	LastName string `json:"last_name" validate:"required,min=5,max=100"`
	URL string `json:"url" validate:"required,min=5,max=100"`
	Address string `json:"address" validate:"required,min=5,max=100"`
	City string `json:"city" validate:"required,min=5,max=100"`
	Country string `json:"country" validate:"required,min=5,max=100"`
	Postal_code string `json:"postal_code" validate:"required,min=5,max=100"`
	Phone string `json:"phone" validate:"required,min=5,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Username string `json:"username" validate:"required,min=5,max=100"`
	Password string `json:"password" validate:"required,min=5,max=100"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}