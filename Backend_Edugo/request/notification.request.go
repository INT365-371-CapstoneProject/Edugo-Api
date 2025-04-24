package request

type CreateNotificationRequest struct {
	Title       string `gorm:"size:100;not null" json:"title"`
	Message     string `gorm:"size:500;not null" json:"message"`
	Is_Read     uint   `gorm:"type:TINYINT(1);not null" json:"is_read"`
	Announce_ID uint   `gorm:"not null" json:"announce_id"`
}
