package response

import "time"

type UserResponse struct {
	Account_ID uint   `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Create_On  time.Time `json:"create_on"`
	Last_Login *time.Time `json:"last_login"`
	Update_On  time.Time `json:"update_on"`
	Role       string `json:"role"`
}