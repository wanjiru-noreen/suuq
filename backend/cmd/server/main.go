package main

import (
	"log"
	"net/http"
	"os"

	"suuq/database"
	"suuq/internal/auth"
	"suuq/routes"
)

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/suuq.db"
	}

	db, err := database.Open(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := database.Migrate(db); err != nil {
		log.Fatal(err)
	}

	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		secret = "development-only-change-me"
	}
	authService := auth.NewService(db, secret)

	// Initialize all application routes
	router := routes.SetupRoutes(authService)

	// Configure the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("Suuq backend running on http://localhost:8080")

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
