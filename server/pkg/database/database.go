package database

import (
	"fmt"
	"learn-tuxedolabs/internal/model/entity"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DBConnect() error {
	dsn := os.Getenv("DATABASE")

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database!")
		return err
	}

	DB = db

  // auto migrate
	DB.AutoMigrate(entity.Users{}, entity.Contacts{}, entity.ResetPasswordToken{})
	return nil
}
