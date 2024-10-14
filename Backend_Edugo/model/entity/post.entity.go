package entity

type Post struct {
	Post_ID       uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
	Description  string `json:"description"`
	URL  string `json:"url"`
	ATTACH_FILE  []byte `json:"attach_file"`
	CATEGORY_ID  uint `json:"category_id"`
	ProviderId  uint `json:"provider_id"`
}