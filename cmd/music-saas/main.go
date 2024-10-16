package main

import (
	"context"
	"fmt"
	"log"
	"music-saas/internal/middleware"
	"music-saas/pkg/api"
	"music-saas/pkg/db"
	"music-saas/pkg/service"
	"net/http"
	"os"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()

	log.Println("Successfully connected to the database")

	// userRepo := db.NewInMemoryUserRepository()
	userRepo := db.NewPostgresUserRepository(dbConn)
	authService := service.NewAuthService(userRepo)
	profileService := service.NewProfileService(userRepo)
	authAPI := api.NewAuthAPI(authService)
	profileAPI := api.NewProfileAPI(profileService)

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", authAPI.Signup)
	mux.HandleFunc("/login", authAPI.Login)

	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("/api/profile", profileAPI.Profile)

	mux.Handle("/api/", middleware.AuthMiddleware(protectedRoutes))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
