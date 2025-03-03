package service

import (
	"fmt"
	"learn-tuxedolabs/internal/middleware"
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/internal/repository"
	"learn-tuxedolabs/pkg/utils"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth"
)

func ValidateLogin(loginRequest *request.LoginRequest) error {
	validate := validator.New()
	return validate.Struct(loginRequest)
}

func ValidateRegister(registerRequest *request.RegisterRequest) error {
	validate := validator.New()
	return validate.Struct(registerRequest)
}

func GetUserByEmail(email string) (*entity.Users, error) {
	return repository.GetUserByEmail(email)
}

func AuthenticateUser(email, password string) (*entity.Users, error) {
	user, err := GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if !middleware.CheckPassword(user.Password, password) {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
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

func RegisterUser(registerRequest *request.RegisterRequest) (*entity.Users, error) {
	_, err := repository.HashAndStoreUser(registerRequest)
	if err != nil {
		return nil, err
	}

	user, err := GetUserByEmail(registerRequest.Email)
	if err != nil {
		return nil, err
	}

	// Send verification email (implementation not shown)
	// sendVerificationEmail(user)

	return user, nil
}

func SaveOAuthUser(oauthUser goth.User) error {
	firstName := oauthUser.FirstName
	if firstName == "" {
		firstName = strings.Split(oauthUser.Email, "@")[0]
	}

	user := entity.Users{
		Email:     oauthUser.Email,
		Name:      firstName,
		FirstName: firstName,
		LastName:  &oauthUser.LastName,
		Password:  "",
    Avatar:    oauthUser.AvatarURL,
		Role:      "member",
		Verify:    true,
	}

	existingUser, err := GetUserByEmail(oauthUser.Email)
	if err == nil && existingUser != nil {
		user.ID = existingUser.ID
	}

	return repository.SaveUser(&user)
}
