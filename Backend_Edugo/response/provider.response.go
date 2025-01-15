package response

import "time"

type ProviderResponse struct {
	Provider_ID  uint       `json:"id"`
	Company_Name string     `json:"company_name"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	URL          string     `json:"url"`
	Address      string     `json:"address"`
	Phone        string     `json:"phone"`
	Status       string     `json:"status"`
	Verify       string     `json:"verify"`
	Create_On    time.Time  `json:"create_on"`
	Last_Login   *time.Time `json:"last_login"`
	Update_On    time.Time  `json:"update_on"`
	Role         string     `json:"role"`
}