package main

import (
	"github.com/litondev/gin-react-crud/api/config"
	"github.com/litondev/gin-react-crud/api/models"
)

func main() {
	database, _ := config.Database()
	database.AutoMigrate(&models.User{})
	database.AutoMigrate(&models.Product{})
	database.AutoMigrate(&models.Data{})
}
