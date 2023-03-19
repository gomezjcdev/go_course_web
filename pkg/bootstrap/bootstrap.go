package bootstrap

import (
	"fmt"
	"github.com/gomezjcdev/go_course_web/internal/user"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func InitLogger() *log.Logger {
	return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
}

func DBConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password='%s' dbname=%s port=%s sslmode=disable TimeZone=America/Bogota",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_USER"),
		os.Getenv("DATABASE_PASSWORD"),
		os.Getenv("DATABASE_NAME"),
		os.Getenv("DATABASE_PORT"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if os.Getenv("DATABASE_DEBUG") == "true" {
		db.Debug()
	}

	if os.Getenv("DATABASE_MIGRATE") == "true" {
		err := db.AutoMigrate(&user.User{})
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
