package main

import (
	datas "github.com/litondev/gin-react-crud/api/migration/datas"
	products "github.com/litondev/gin-react-crud/api/migration/products"
	users "github.com/litondev/gin-react-crud/api/migration/users"
)

func main() {
	users.Migrate()
	datas.Migrate()
	products.Migrate()
}
