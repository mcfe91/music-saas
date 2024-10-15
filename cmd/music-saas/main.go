package main

import (
	// "database/sql"
	"log"
	"music-saas/internal/middleware"
	"music-saas/pkg/api"
	"music-saas/pkg/db"
	"music-saas/pkg/service"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// Uncomment this to use PostgreSQL
	// connStr := "user=username dbname=mydb sslmode=disable"
	// db, err := sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()
	// userRepo := service.NewPostgresUserRepository(db)

	userRepo := db.NewInMemoryUserRepository()
	authService := service.NewAuthService(userRepo)
	apiHandler := api.NewAPIHandler(authService)

	mux := http.NewServeMux()
	mux.HandleFunc("/signup", apiHandler.Signup)
	mux.HandleFunc("/login", apiHandler.Login)

	protectedRoutes := http.NewServeMux()
	protectedRoutes.HandleFunc("/api/profile", apiHandler.Profile)

	mux.Handle("/api/", middleware.AuthMiddleware(protectedRoutes))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
