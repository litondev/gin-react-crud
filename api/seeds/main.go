package seeds

import "github.com/litondev/gin-react-crud/api/db"

func main() {
	database := db.Database()
	database.AutoMigrate(&model.User{})
	database.AutoMigrate(&model.Product{})
	database.AutoMigrate(&model.Datas{})
}
