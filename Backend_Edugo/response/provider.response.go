package response

import "time"

type ProviderResponse struct {
	Provider_ID  uint       `json:"id"`
	Company_Name string     `json:"company_name"`
	FirstName   *string     `json:"first_name"`
	LastName    *string     `json:"last_name"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	URL          string     `json:"url"`
	Address      string     `json:"address"`
	City         string     `json:"city"`
	Country      string     `json:"country"`
	Postal_Code  string     `json:"postal_code"`
	Phone        string     `json:"phone"`
	Status       string     `json:"status"`
	Verify       string     `json:"verify"`
	Create_On    time.Time  `json:"create_on"`
	Last_Login   *time.Time `json:"last_login"`
	Update_On    time.Time  `json:"update_on"`
	Role         string     `json:"role"`
}

type ProviderDetails struct {
	Provider_ID  uint   `json:"id"`
	Company_Name string `json:"company_name"`
	URL          string `json:"url"`
	Address      string `json:"address"`
	City         string `json:"city"`
	Country      string `json:"country"`
	Postal_Code  string `json:"postal_code"`
	Phone        string `json:"phone"`
	Phone_Person *string `json:"phone_person"`
	Verify	   string `json:"verify"`
}