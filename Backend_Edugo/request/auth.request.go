package request

type LoginRequest struct {
	Email      string `json:"email" validate:"required_without=Username,omitempty,email" message:"Email is required when username is not provided"`
	Username   string `json:"username" validate:"required_without=Email,omitempty" message:"Username is required when email is not provided"`
	Password   string `json:"password" validate:"required" message:"Password is required"`
	RememberMe bool   `json:"remember_me"`
}

type EditProfileRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
}

type AdminUpdateRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2"`
	Phone     *string `json:"phone" validate:"omitempty"`
}

type ProviderUpdateRequest struct {
	Username    *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email       *string `json:"email" validate:"omitempty,email"`
	FirstName   *string `json:"first_name" validate:"omitempty,min=2"`
	LastName    *string `json:"last_name" validate:"omitempty,min=2"`
	CompanyName *string `json:"company_name" validate:"omitempty"`
	Phone       *string `json:"phone" validate:"omitempty"`
	Address     *string `json:"address" validate:"omitempty"`
	City        *string `json:"city" validate:"omitempty"`
	Country     *string `json:"country" validate:"omitempty"`
	PostalCode  *string `json:"postal_code" validate:"omitempty"`
}

type UserUpdateRequest struct {
	Username  *string `json:"username" validate:"omitempty,min=3,max=50"`
	Email     *string `json:"email" validate:"omitempty,email"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2"`
}
