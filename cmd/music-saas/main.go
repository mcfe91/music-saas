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

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Database connection parameters
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)

	// Connect to the database
	dbConn, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()

	log.Println("Successfully connected to the database")

	// Repositories
	// userRepo := db.NewInMemoryUserRepository()
	userRepo := db.NewPostgresUserRepository(dbConn)
	productRepo := db.NewPostgresProductRepository(dbConn)

	// Services
	authService := service.NewAuthService(userRepo)
	profileService := service.NewProfileService(userRepo)
	productService := service.NewProductService(productRepo)

	// API Handlers
	authAPI := api.NewAuthAPI(authService)
	profileAPI := api.NewProfileAPI(profileService)
	productAPI := api.NewProductHandler(productService)

	// Router setup using Gorilla mux
	r := mux.NewRouter()

	// Authentication routes
	r.HandleFunc("/signup", authAPI.Signup).Methods("POST")
	r.HandleFunc("/login", authAPI.Login).Methods("POST")

	// Profile route (protected)
	protectedRoutes := mux.NewRouter()
	protectedRoutes.HandleFunc("/api/profile", profileAPI.Profile).Methods("GET")
	r.PathPrefix("/api/").Handler(middleware.AuthMiddleware(protectedRoutes))

	// Product routes
	r.HandleFunc("/products", productAPI.CreateProduct).Methods("POST")
	r.HandleFunc("/products/{id}", productAPI.GetProduct).Methods("GET")
	r.HandleFunc("/products/{id}", productAPI.UpdateProduct).Methods("PUT")
	r.HandleFunc("/products/{id}", productAPI.DeleteProduct).Methods("DELETE")
	r.HandleFunc("/products", productAPI.GetProducts).Methods("GET")

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
