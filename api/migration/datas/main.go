package users

import (
	"github.com/litondev/gin-react-crud/api/db"
	model "github.com/litondev/gin-react-crud/api/models"
)

func Migrate() {
	database := db.Database()
	database.AutoMigrate(&model.Data{})
}
