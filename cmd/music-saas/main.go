package main

import (
	"log"
	"music-saas/pkg/api"
	"net/http"
)

func main() {
	http.HandleFunc("/signup", api.Signup)
	http.HandleFunc("/login", api.Login)

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
