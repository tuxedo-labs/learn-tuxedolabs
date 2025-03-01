package main

import (
	"fmt"
	handlers "learn-tuxedolabs/internal/handler"
	"learn-tuxedolabs/internal/middleware"
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

	r := mux.NewRouter()

	authRoutes := r.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/login", handlers.Login).Methods("POST")
	authRoutes.HandleFunc("/register", handlers.Register).Methods("POST")

	r.HandleFunc("/health", middleware.Auth(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})).Methods("GET")

	logServiceStart(port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}

func logServiceStart(port string) {
	startTime := time.Now().Format(time.RFC1123)
	message := fmt.Sprintf("🚀 Service running on http://localhost:%s | Started at: %s", port, startTime)
	log.Println(message)
}
