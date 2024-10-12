package request

type UserCreateRequest struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	AccountAccountID int  `json:"account_account_id"`
}