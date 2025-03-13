package response

type AdminResponse struct {
	Admin_ID       uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone"`
	Status    string `json:"status"`
}

type AdminDetails struct {
	Admin_ID  uint   `json:"id"`
	Phone	 string `json:"phone"`
}
