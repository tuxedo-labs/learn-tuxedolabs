package service

import (
	"encoding/json"
	"fmt"
	"learn-tuxedolabs/internal/middleware"
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/internal/repository"
	"learn-tuxedolabs/pkg/config"
	"learn-tuxedolabs/pkg/utils"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-playground/validator/v10"
	"github.com/markbates/goth"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateLogin(loginRequest *request.LoginRequest) error {
	return validate.Struct(loginRequest)
}

func ValidateRegister(registerRequest *request.RegisterRequest) error {
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

	return utils.GenerateToken(&claims)
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

func SendForgetPasswordEmail(email string) {
  // sementara stop
  return nil
}

func FetchGitHubEmail(accessToken string) (string, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var emails []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return "", fmt.Errorf("failed to decode user emails: %w", err)
	}

	for _, email := range emails {
		if email["primary"].(bool) && email["verified"].(bool) {
			return email["email"].(string), nil
		}
	}

	return "", fmt.Errorf("no primary verified email found")
}
