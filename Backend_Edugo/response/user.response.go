package response

type UserResponse struct {
	User_ID uint   `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Phone_Number string `json:"phone_number"`
	Update_On string `json:"update_on"`
	Last_Login string `json:"last_login"`
	Username string `json:"username"`
	Email string `json:"email"`
	Role string `json:"role"`
}