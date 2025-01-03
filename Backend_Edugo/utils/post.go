package utils

import (

	// "github.com/tk-neng/demo-go-fiber/database"
	// "github.com/tk-neng/demo-go-fiber/model/entity"
)

// func GetPostByAnnounceID(posts []entity.Announce_Post) error {
// 	// คิวรี่ข้อมูลใน Post โดยใช้เงื่อนไข posts_id ในตาราง Announce_Post
// 	for i, announcePost := range posts {
// 		var post entity.Post
// 		database.DB.First(&post, "posts_id = ?", announcePost.Posts_ID)
// 		//เพิ่มข้อมูลลงในตัวแปร posts
// 		posts[i].Post = post
// 		// ใช้ฟังก์ชัน GetCountryName จากไฟล์ utils/post.go
// 		if err := GetCountryName(&post); err != nil {
// 			return err
// 		}
// 		// กำหนดค่าให้กับตัวแปร Country ใน posts
// 		posts[i].Post.Country = post.Country
// 	}
// 	return nil
// }

// func GetCountryName(post *entity.Post) error {
// 	var country entity.Country
// 	database.DB.First(&country, "country_id = ?", post.Country_ID)
// 	post.Country = country
// 	return nil
// }

// func GetCountryNamePost(post *entity.Post) error {
// 	var country entity.Country
// 	database.DB.First(&country, "country_id = ?", post.Country_ID)
// 	post.Country = country
// 	return nil
// }