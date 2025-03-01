package handlers

import (
	"encoding/json"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/internal/repository"
	"learn-tuxedolabs/internal/service"
	"learn-tuxedolabs/pkg/utils"
	"net/http"
	"strings"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := service.ValidateLogin(&loginRequest); errValidate != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Validation failed", "error": errValidate.Error()})
		return
	}

	user, err := service.AuthenticateUser(loginRequest.Email, loginRequest.Password)
	if err != nil {
		utils.RespondJSON(w, http.StatusUnauthorized, map[string]string{"message": "Invalid email or password"})
		return
	}

	if !user.Verify {
		utils.RespondJSON(w, http.StatusForbidden, map[string]string{"message": "Account not verified. Please check your email for verification instructions."})
		return
	}

	accessToken, err := service.GenerateAccessToken(*user)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"access_token": accessToken,
	})
}

func Register(w http.ResponseWriter, r *http.Request) {
	var registerRequest request.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := service.ValidateRegister(&registerRequest); errValidate != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Validation failed", "error": errValidate.Error()})
		return
	}

	result, err := repository.HashAndStoreUser(&registerRequest)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Email already in use"})
			return
		}
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to register user"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  result,
		"message": "Registration successful! Please check your email for the verification code",
	})
}
