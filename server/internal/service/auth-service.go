package service

import (
	"fmt"
	"learn-tuxedolabs/internal/middleware"
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/pkg/database"
	"learn-tuxedolabs/pkg/utils"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func ValidateLogin(loginRequest *request.LoginRequest) error {
	validate := validator.New()
	return validate.Struct(loginRequest)
}

func GetUserByEmail(email string) (*entity.Users, error) {
	var user entity.Users
	err := database.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func GenerateJWTToken(user *entity.Users) (string, error) {
	claims := jwt.MapClaims{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"exp":   time.Now().Add(5 * time.Minute).Unix(),
		"role":  "member",
	}

	if user.Role == "admin" {
		claims["role"] = "admin"
	}

	return utils.GenerateToken(&claims)
}

func ValidateRegister(registerRequest *request.RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(registerRequest)
}

func AuthenticateUser(email, password string) (*entity.Users, error) {
	var user entity.Users
	err := database.DB.First(&user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	if !middleware.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

func GenerateAccessToken(user entity.Users) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}

	token, err := utils.GenerateToken(&claims)
	if err != nil {
		return "", err
	}

	return token, nil
}

