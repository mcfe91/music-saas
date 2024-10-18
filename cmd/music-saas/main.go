package main

import (
	"fmt"
	"log"
	"music-saas/internal/middleware"
	"music-saas/pkg/api"
	"music-saas/pkg/db"
	"music-saas/pkg/service"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"database/sql"
	_ "github.com/lib/pq"
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
	jwtKey := []byte(os.Getenv("JWT_SECRET"))

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=require", dbUser, dbPassword, dbHost, dbPort, dbName)

	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbConn.Close()

	// Check the connection
	err = dbConn.Ping()
	if err != nil {
		log.Fatalf("Unable to ping the database: %v\n", err)
	}

	// Repositories
	userRepo := db.NewPostgresUserRepository(dbConn)
	productRepo := db.NewPostgresProductRepository(dbConn)
	cartRepo := db.NewPostgresCartRepository(dbConn)
	cartItemRepo := db.NewPostgresCartItemRepository(dbConn)
	orderRepo := db.NewPostgresOrderRepository(dbConn)

	// Services
	authService := service.NewAuthService(userRepo, jwtKey)
	profileService := service.NewProfileService(userRepo)
	productService := service.NewProductService(productRepo)
	cartService := service.NewCartService(cartRepo, cartItemRepo)
	orderService := service.NewOrderService(orderRepo, productRepo)

	// API Handlers
	authAPI := api.NewAuthAPI(authService)
	productAPI := api.NewProductHandler(productService)
	profileAPI := api.NewProfileHandler(profileService)
	cartAPI := api.NewCartHandler(cartService)
	orderAPI := api.NewOrderHandler(orderService)

	// Router setup using Gorilla mux
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/signup", authAPI.Signup).Methods("POST")
	r.HandleFunc("/login", authAPI.Login).Methods("POST")
	r.HandleFunc("/products/search", productAPI.SearchProducts).Methods("GET")
	r.HandleFunc("/products/{id:[0-9]+}", productAPI.GetProduct).Methods("GET")
	r.HandleFunc("/products", productAPI.GetProducts).Methods("GET")

	// Protected routes
	protectedRouter := r.PathPrefix("/api").Subrouter()
	protectedRouter.Use(middleware.AuthMiddleware(authService, jwtKey))
	protectedRouter.HandleFunc("/profile", profileAPI.GetProfile).Methods("GET")
	protectedRouter.HandleFunc("/cart", cartAPI.AddToCart).Methods("POST")
	protectedRouter.HandleFunc("/cart", cartAPI.RemoveFromCart).Methods("DELETE")
	protectedRouter.HandleFunc("/orders", orderAPI.CreateOrder).Methods("POST")
	protectedRouter.HandleFunc("/orders/{id}/items", orderAPI.GetOrderItems).Methods("GET")

	// Admin routes
	adminRouter := r.PathPrefix("/api/admin").Subrouter()
	adminRouter.Use(middleware.AuthMiddleware(authService, jwtKey))
	adminRouter.Use(middleware.AdminMiddleware)
	adminRouter.HandleFunc("/products", productAPI.CreateProduct).Methods("POST")
	adminRouter.HandleFunc("/products/{id:[0-9]+}", productAPI.UpdateProduct).Methods("PUT")
	adminRouter.HandleFunc("/products/{id:[0-9]+}", productAPI.DeleteProduct).Methods("DELETE")

	// Start the server
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
