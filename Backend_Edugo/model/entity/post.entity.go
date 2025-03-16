package entity

import "time"

type Post struct {
	Posts_ID     uint       `json:"id" gorm:"primaryKey"`
	Description  string     `json:"description"`
	Image        []byte     `gorm:"type:longblob" json:"image"`
	Publish_Date *time.Time `json:"published_date"`
	Account_ID   uint       `json:"account_id"`

	Account Account `gorm:"foreignKey:Account_ID;references:Account_ID"`
}
