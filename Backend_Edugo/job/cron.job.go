package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/tk-neng/demo-go-fiber/model/entity"
	"gorm.io/gorm"
)

// InitCronJob ตั้งค่าและเริ่ม Cron Job
func InitCronJob(db *gorm.DB) {
	c := cron.New()

	// ตั้งให้รันทุกๆ 10 วินาที
	c.AddFunc("@every 30m", func() {
		// fmt.Println("Running scheduled check for announcements...")

		// คำนวณเวลาปัจจุบันและอีก 1 ชั่วโมงข้างหน้า
		now := time.Now()
		oneHourLater := now.Add(1 * time.Hour)

		// ค้นหาประกาศที่จะหมดอายุในอีก 1 ชั่วโมง
		var closingAnnouncements []entity.Announce_Post
		closingResult := db.Where("close_date BETWEEN ? AND ?", now, oneHourLater).Find(&closingAnnouncements)

		if closingResult.Error != nil {
			// fmt.Println("Error fetching announcements closing soon:", closingResult.Error)
			return
		}

		if len(closingAnnouncements) == 0 {
			// fmt.Println("No announcements closing in 1 hour")
		} else {
			// fmt.Printf("Found %d announcements closing in 1 hour:\n", len(closingAnnouncements))
			announceIDs := make([]int, 0) // เก็บ ID ของประกาศที่กำลังจะปิด

			for _, announcement := range closingAnnouncements {
				// fmt.Printf("ID: %d, Title: %s, Close Date: %s\n",
				// 	announcement.Announce_ID, announcement.Title, announcement.Close_Date)
				announceIDs = append(announceIDs, int(announcement.Announce_ID))
			}

			// ค้นหา Bookmark ที่เกี่ยวข้องกับประกาศที่กำลังจะปิด
			var bookmarks []entity.Bookmark
			bookmarkResult := db.
				Where("announce_id IN ?", announceIDs).
				Select("MIN(bookmark_id) as bookmark_id, account_id, announce_id").
				Group("account_id, announce_id").
				Find(&bookmarks)

			if bookmarkResult.Error != nil {
				// fmt.Println("Error fetching bookmarks related to closing announcements:", bookmarkResult.Error)
				return
			}

			if len(bookmarks) == 0 {
				// fmt.Println("No bookmarks found for closing announcements")
			} else {
				// fmt.Printf("Found %d bookmarks for closing announcements:\n", len(bookmarks))
				for _, bookmark := range bookmarks {
					// fmt.Printf("Bookmark ID: %d, Account ID: %d, Announce ID: %d\n",
					// 	bookmark.Bookmark_ID, bookmark.Account_ID, bookmark.Announce_ID)

					// เช็คว่าใน notification มีที่ตรงกับ bookmark หรือไม่
					var existingNotification entity.Notification
					notificationResult := db.Where("account_id = ? AND announce_id = ?", bookmark.Account_ID, bookmark.Announce_ID).First(&existingNotification)

					// ถ้าไม่มี notification ที่ตรงกันให้สร้างใหม่
					if notificationResult.Error != nil {
						// ค้นหาข้อมูลของประกาศ
						var announcement entity.Announce_Post
						announceResult := db.Where("announce_id = ?", bookmark.Announce_ID).First(&announcement)
						if announceResult.Error != nil {
							// fmt.Println("Error fetching announcement:", announceResult.Error)
							continue
						}

						// สร้าง Notification ใหม่
						newNotification := entity.Notification{
							Title:       announcement.Title,
							Message:     "Add New Noti",
							IsRead:      0,
							CreatedAt:   time.Now(),
							Account_ID:  bookmark.Account_ID,
							Announce_ID: bookmark.Announce_ID,
						}

						// บันทึก Notification ใหม่
						createResult := db.Create(&newNotification)
						if createResult.Error != nil {
							fmt.Println("Error creating new notification:", createResult.Error)
						} else {
							// fmt.Printf("Created new notification for Account ID: %d, Announce ID: %d\n", bookmark.Account_ID, bookmark.Announce_ID)
						}
					} else {
						// fmt.Printf("Notification already exists for Account ID: %d, Announce ID: %d\n", bookmark.Account_ID, bookmark.Announce_ID)
					}
				}
			}
		}

		// ค้นหาประกาศที่จะหมดอายุในอีก 1 ชั่วโมง
		var notificationList []entity.Notification
		notificationResult := db.Find(&notificationList)

		if notificationResult.Error != nil {
			// fmt.Println("Error fetching notifications:", notificationResult.Error)
			return
		}

		if len(notificationList) == 0 {
			// fmt.Println("No notifications found")
		} else {
			// fmt.Printf("Found %d notifications\n", len(notificationList))
			// for _, notification := range notificationList {
			// 	fmt.Printf("Notification ID: %d, Account ID: %d, Announce ID: %d\n",
			// 		notification.NotificationID, notification.Account_ID, notification.Announce_ID)
			// }
		}
	})

	// เริ่ม Cron Job
	c.Start()
}
