package main

import (
	"github.com/litondev/gin-react-crud/api/db"
	"github.com/litondev/gin-react-crud/api/model"
)

func main() {
	database := db.Database()
	database.AutoMigrate(&model.User{})
}
