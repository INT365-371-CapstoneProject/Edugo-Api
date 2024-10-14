package request

type UserCreateRequest struct {
	Firstname string `json:"first_name"`
	Lastname  string `json:"last_name"`
	AccountID int  `json:"account_account_id"`
}