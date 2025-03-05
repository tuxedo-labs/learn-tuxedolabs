package main

import (
	"fmt"
	"learn-tuxedolabs/internal/handler"
	"learn-tuxedolabs/internal/middleware"
	"learn-tuxedolabs/pkg/database"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("PORT")
	if port == "" {
		panic("PORT environment variable is not set")
	}

	err = database.DBConnect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	r := mux.NewRouter()

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", handler.Login).Methods("POST")
	authRoutes.HandleFunc("/register", handler.Register).Methods("POST")
	authRoutes.HandleFunc("/{provider}/login", handler.OAuthLogin).Methods("GET")
	authRoutes.HandleFunc("/{provider}/callback", handler.OAuthCallback).Methods("GET")
	authRoutes.HandleFunc("/logout", handler.Logout).Methods("GET")

	userRoutes := r.PathPrefix("/user").Subrouter()
	userRoutes.HandleFunc("/profile", middleware.Auth(handler.UserProfile)).Methods("GET")
  // userRoutes.HandleFunc("/profile", middleware.Auth(handler.UpdateProfile)).Methods("PUT")

	logServiceStart(port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func logServiceStart(port string) {
	startTime := time.Now().Format(time.RFC1123)
	message := fmt.Sprintf("🚀 Service running on http://localhost:%s | Started at: %s", port, startTime)
	log.Println(message)
}
