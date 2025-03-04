package repository

import (
	"fmt"
	"learn-tuxedolabs/internal/middleware"
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/pkg/database"
)

func UpdateUser(user *entity.Users) error {
	return database.DB.Save(user).Error
}

func HashAndStoreUser(registerRequest *request.RegisterRequest) (string, error) {
	var existingUser entity.Users
	if err := database.DB.First(&existingUser, "email = ?", registerRequest.Email).Error; err == nil {
		return "", fmt.Errorf("user with email %s already exists", registerRequest.Email)
	}

	hashedPassword, err := middleware.HashPassword(registerRequest.Password)
	if err != nil {
		return "", err
	}

	newUser := entity.Users{
		Name:      fmt.Sprintf("%s %s", registerRequest.FirstName, registerRequest.LastName),
		FirstName: registerRequest.FirstName,
		LastName:  &registerRequest.LastName,
		Email:     registerRequest.Email,
		Avatar:    "",
		Password:  hashedPassword,
		Role:      "member",
		Verify:    true,
	}

	if err := database.DB.Create(&newUser).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("User %s registered successfully", newUser.Email), nil
}

func GetUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	err := database.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func SaveUser(user *entity.Users) error {
	return database.DB.Save(user).Error
}
