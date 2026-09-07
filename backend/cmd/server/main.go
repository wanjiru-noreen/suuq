package main

import (
	"log"
	"net/http"
	"os"

	"suuq/database"
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

	// Initialize all application routes
	router := routes.SetupRoutes()

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
