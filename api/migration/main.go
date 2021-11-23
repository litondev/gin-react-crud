package migration

import (
	"github.com/litondev/gin-react-crud/api/db"
)

func main() {
	database := db.Database()
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.Data{})
}
