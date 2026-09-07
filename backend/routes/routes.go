package routes

import (
	"encoding/json"
	"net/http"

	"suuq/internal/auth"
)

// SetupRoutes registers all API endpoints and returns the router.
func SetupRoutes(authService *auth.Service) http.Handler {
	mux := http.NewServeMux()

	// Health check route
	mux.HandleFunc("/api/health", healthHandler)
	mux.HandleFunc("POST /api/auth/register", authService.Register)
	mux.HandleFunc("POST /api/auth/login", authService.Login)
	mux.Handle("GET /api/auth/me", authService.Middleware(http.HandlerFunc(authService.Me)))

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
