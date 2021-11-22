import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
)

func MigrateUser() {
	dsn := "host=" + os.Getenv("DB_HOST") + 
			" user=" + os.Getenv("DB_USER") + 
			" password=" + os.Getenv("DB_PASSWORD") + 
			" dbname=" + os.Getenv("DB_NAME") + 
			" port=" + os.Getenv("DB_PORT") + 
			" sslmode=disable TimeZone=Asia/Jakarta";

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
}