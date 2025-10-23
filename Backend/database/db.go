package database

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	if err := godotenv.Load(); err != nil {
		log.Println("tidak menemukan file .env")
	}
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Data Source Na,=me (dSn)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s/%s)?charset=utf8mb4&parseTime=True&loc=Local",
		dbUser,
		dbPass,
		dbHost,
		dbPort,
		dbName)

	var err error

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(
			"gagal connect ke DB: %v", err,
		)
	}
	fmt.Println("berhasil connect ke Db")

	err = DB.AutoMigrate(
		&models.User{},
		&models.ScheduledJob{},
		&models.Log{},
		&models.Monitoring{},
	)
	if err != nil {
		log.Fatalf("gagal migrasi tabel: %v", err)
	}
	fmt.Println("migrasi tabel selesai")

	return DB
}
