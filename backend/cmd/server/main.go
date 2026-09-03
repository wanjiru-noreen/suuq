package main

import (
	"log"
	"net/http"

	"suuq/routes"
)

func main() {
	// Initialize all application routes
	router := routes.SetupRoutes()

	// Configure the HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	log.Println("🚀 Suuq backend running on http://localhost:8080")

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
