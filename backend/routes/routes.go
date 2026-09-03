package routes

import (
	"encoding/json"
	"net/http"
)

// SetupRoutes registers all API endpoints and returns the router.
func SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Health check route
	mux.HandleFunc("/api/health", healthHandler)

	return mux
}

// healthHandler confirms the backend is running.
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	response := map[string]string{
		"status":  "ok",
		"message": "Suuq API is running",
	}

	json.NewEncoder(w).Encode(response)
}
