package request

type UserRequest struct {
	Email           string `json:"email" validate:"required,email"`
	Username        string `json:"username" validate:"required"`
	Password        string `json:"password" validate:"required"`
	FirstName       string `json:"first_name" validate:"required"`
	LastName        string `json:"last_name" validate:"required"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type QuestionRequest struct {
	Education_Level string `json:"education_level" validate:"required,oneof=Undergraduate Master Doctorate"`
	Countries       []uint `json:"countries" validate:"required,min=1,max=3,dive,number"`
	Categories      []uint `json:"categories" validate:"required,min=1,max=3,dive,number"`
}
