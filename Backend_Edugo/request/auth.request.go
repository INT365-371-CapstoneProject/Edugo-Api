package request

type LoginRequest struct {
	Email    string `json:"email" validate:"required_without=Username,omitempty,email" message:"Email is required when username is not provided"`
	Username string `json:"username" validate:"required_without=Email,omitempty" message:"Username is required when email is not provided"`
	Password string `json:"password" validate:"required" message:"Password is required"`
	RememberMe bool `json:"remember_me"`
}