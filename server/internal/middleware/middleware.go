package middleware

import (
	"learn-tuxedolabs/internal/model/entity"
	"learn-tuxedolabs/pkg/database"
	"learn-tuxedolabs/pkg/utils"
	"net/http"

	"context"

	"golang.org/x/crypto/bcrypt"
)

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("x-token")
		if token == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.DecodeToken(token)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := uint(claims["user_id"].(float64))
		var user entity.Users
		if err := database.DB.First(&user, userID).Error; err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "usersInfo", claims)
		ctx = context.WithValue(ctx, "role", claims["role"])
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func AdminRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Context().Value("role")

		if role == "member" {
			http.Error(w, "forbidden access", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
