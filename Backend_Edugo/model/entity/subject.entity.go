package entity

type Subject struct {
	Subject_ID       uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Description  string `json:"description"`
	AttachFile []byte `json:"attach_file"`
	TagId  uint `json:"tag_id"`
	UserID  uint `json:"user_id"`
}