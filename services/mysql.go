package services

import (
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"gorm.io/driver/mysql"
)

var DB *gorm.DB

func InitMySQL() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dsn := dbUser+ ":" +dbPass+ "@tcp(127.0.0.1:" +dbPort+ ")/" +dbName+ "?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db
}