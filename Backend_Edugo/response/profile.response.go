package response

// ProfileResponse struct
type ProfileResponse struct {
	ID        uint    `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      string  `json:"role"`
}

// AdminProfileResponse struct
type AdminProfileResponse struct {
	ID        uint    `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      string  `json:"role"`
	Phone     *string `json:"phone_number"`
}

// ProviderProfileResponse struct
type ProviderProfileResponse struct {
	ID           uint   `json:"id"`
	Email        string `json:"email"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Role         string `json:"role"`
	Company_Name string `json:"company_name"`
	Phone        string `json:"phone"`
	Phone_Person string `json:"phone_person"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Postal_Code  string `json:"postal_code"`
}

// UserProfileResponse struct
type UserProfileResponse struct {
	ID        uint    `json:"id"`
	Email     string  `json:"email"`
	Username  string  `json:"username"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Role      string  `json:"role"`
}
