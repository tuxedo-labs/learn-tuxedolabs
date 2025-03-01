package config

import (
	"fmt"
	entity "learn-tuxedolabs/internal/entity/model"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	dsn := os.Getenv("MYSQL")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database!")
		return err
	}

	DB = db

  // auto migrate
	DB.AutoMigrate(entity.Users{}, entity.Contacts{})

	return nil
}
