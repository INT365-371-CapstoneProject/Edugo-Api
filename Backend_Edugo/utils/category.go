package utils

import (
	"github.com/tk-neng/demo-go-fiber/database"
	"github.com/tk-neng/demo-go-fiber/model/entity"
)

func GetCategoryName(posts []entity.Announce_Post) error {
	for i, announcePost := range posts {
		var category entity.Category
		database.DB.First(&category, "category_id = ?", announcePost.Category_ID)
		posts[i].Category = category
	}
	return nil
}