package handler

import (
	"encoding/json"
	"learn-tuxedolabs/internal/model/request"
	"learn-tuxedolabs/internal/service"
	"learn-tuxedolabs/pkg/utils"
	"net/http"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
)

var validate *validator.Validate

func init() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	goth.UseProviders(
		google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), os.Getenv("GOOGLE_CALLBACK_URL"), "openid", "profile", "email"),
		github.New(os.Getenv("GITHUB_KEY"), os.Getenv("GITHUB_SECRET"), os.Getenv("GITHUB_CALLBACK_URL"), "user:email"),
	)

	gothic.Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET")))

	validate = validator.New()
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginRequest request.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := validate.Struct(&loginRequest); errValidate != nil {
		validationErrors := utils.ParseValidationErrors(errValidate)
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Validation failed", "errors": validationErrors})
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

	if errValidate := validate.Struct(&registerRequest); errValidate != nil {
		validationErrors := utils.ParseValidationErrors(errValidate)
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Validation failed", "errors": validationErrors})
		return
	}

	user, err := service.RegisterUser(&registerRequest)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.RespondJSON(w, http.StatusConflict, map[string]string{"message": "Email already in use"})
			return
		}
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to register user"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "success",
		"user_id": user.ID,
		"message": "Registration successful! Please check your email for the verification code",
	})
}

func OAuthLogin(w http.ResponseWriter, r *http.Request) {
	fromURL := r.URL.Query().Get("from")
	if fromURL != "" {
		session, _ := gothic.Store.Get(r, "redirect")
		session.Values["from"] = fromURL
		session.Save(r, w)
	} else {
		session, _ := gothic.Store.Get(r, "redirect")
		delete(session.Values, "from")
		session.Save(r, w)
	}

	gothic.BeginAuthHandler(w, r)
}

func OAuthCallback(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		utils.RespondJSON(w, http.StatusTemporaryRedirect, map[string]string{"message": "Error completing user authentication"})
		return
	}

	if user.Provider == "github" && user.Email == "" {
		email, err := service.FetchGitHubEmail(user.AccessToken)
		if err != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error fetching user email"})
			return
		}
		user.Email = email
	}

	existingUser, err := service.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		accessToken, err := service.GenerateAccessToken(*existingUser)
		if err != nil {
			utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
			return
		}
		redirectOrRespond(w, r, accessToken, "User already registered")
		return
	}

	err = service.SaveOAuthUser(user)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to save user data"})
		return
	}

	newUser, err := service.GetUserByEmail(user.Email)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error retrieving user data"})
		return
	}

	accessToken, err := service.GenerateAccessToken(*newUser)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Error generating access token"})
		return
	}

	redirectOrRespond(w, r, accessToken, "Successfully registered")
}

func ForgetPassword(w http.ResponseWriter, r *http.Request) {
	var forgetPasswordRequest request.ForgetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&forgetPasswordRequest); err != nil {
		utils.RespondJSON(w, http.StatusBadRequest, map[string]string{"message": "Invalid request payload"})
		return
	}

	if errValidate := validate.Struct(&forgetPasswordRequest); errValidate != nil {
		validationErrors := utils.ParseValidationErrors(errValidate)
		utils.RespondJSON(w, http.StatusBadRequest, map[string]interface{}{"message": "Validation failed", "errors": validationErrors})
		return
	}

	existingUser, err := service.GetUserByEmail(forgetPasswordRequest.Email)
	if err != nil || existingUser == nil {
		utils.RespondJSON(w, http.StatusNotFound, map[string]string{"message": "User not found"})
		return
	}

	if err := service.SendForgetPasswordEmail(existingUser); err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to send password reset email"})
		return
	}

	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "Password reset instructions sent to your email."})
}

func Logout(w http.ResponseWriter, r *http.Request) {
	err := gothic.Logout(w, r)
	if err != nil {
		utils.RespondJSON(w, http.StatusInternalServerError, map[string]string{"message": "Failed to logout"})
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "You have been logged out."})
}

func redirectOrRespond(w http.ResponseWriter, r *http.Request, token string, message string) {
	session, _ := gothic.Store.Get(r, "redirect")
	fromURL, ok := session.Values["from"].(string)
	if !ok || fromURL == "" {
		utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
			"access_token": token,
			"message":      message,
		})
		return
	}

	http.Redirect(w, r, fromURL+"?accessToken="+token+"&message="+message, http.StatusSeeOther)
}
