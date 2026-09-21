package routes

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"suuq/database"
	"suuq/internal/auth"
)

func testRouter(t *testing.T) http.Handler {
	t.Helper()

	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	authService := auth.NewService(db, "test-secret")

	return SetupRoutes(authService)
}

func TestHealthRoute(t *testing.T) {
	router := testRouter(t)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/health",
		nil,
	)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("health status = %d, want %d", response.Code, http.StatusOK)
	}
}

func TestRegisterRoute(t *testing.T) {
	router := testRouter(t)

	request := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/register",
		strings.NewReader(`{"name":"Amina","email":"amina@example.com","password":"password123"}`),
	)
	request.Header.Set("Content-Type", "application/json")

	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", response.Code, http.StatusCreated)
	}
}

func TestLoginRoute(t *testing.T) {
	router := testRouter(t)

	register := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/register",
		strings.NewReader(`{"name":"Amina","email":"amina@example.com","password":"password123"}`),
	)
	register.Header.Set("Content-Type", "application/json")

	registerResponse := httptest.NewRecorder()
	router.ServeHTTP(registerResponse, register)

	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerResponse.Code, http.StatusCreated)
	}

	login := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/login",
		strings.NewReader(`{"email":"amina@example.com","password":"password123"}`),
	)
	login.Header.Set("Content-Type", "application/json")

	loginResponse := httptest.NewRecorder()
	router.ServeHTTP(loginResponse, login)

	if loginResponse.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loginResponse.Code, http.StatusOK)
	}
}

func TestMeRouteRequiresAuthentication(t *testing.T) {
	router := testRouter(t)

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/auth/me",
		nil,
	)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("me status = %d, want %d", response.Code, http.StatusUnauthorized)
	}
}

func TestMeRouteWithAuthentication(t *testing.T) {
	router := testRouter(t)

	register := httptest.NewRequest(
		http.MethodPost,
		"/api/auth/register",
		strings.NewReader(`{"name":"Amina","email":"amina@example.com","password":"password123"}`),
	)
	register.Header.Set("Content-Type", "application/json")

	registerResponse := httptest.NewRecorder()
	router.ServeHTTP(registerResponse, register)

	if registerResponse.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registerResponse.Code, http.StatusCreated)
	}

	var body struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(registerResponse.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}

	if body.Token == "" {
		t.Fatal("register did not return a token")
	}

	request := httptest.NewRequest(
		http.MethodGet,
		"/api/auth/me",
		nil,
	)
	request.Header.Set("Authorization", "Bearer "+body.Token)

	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("me status = %d, want %d", response.Code, http.StatusOK)
	}
}
