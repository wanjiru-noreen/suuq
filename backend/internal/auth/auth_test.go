package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"suuq/database"
)

func testService(t *testing.T) *Service {
	t.Helper()
	db, err := database.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	if err := database.Migrate(db); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return NewService(db, "test-secret")
}

func TestRegisterLoginAndMe(t *testing.T) {
	service := testService(t)
	register := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"name":"Amina","email":"AMINA@example.com","password":"password123"}`))
	registered := httptest.NewRecorder()
	service.Register(registered, register)
	if registered.Code != http.StatusCreated {
		t.Fatalf("register status = %d, want %d", registered.Code, http.StatusCreated)
	}

	var response struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(registered.Body).Decode(&response); err != nil || response.Token == "" {
		t.Fatal("register did not return a token")
	}

	login := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"amina@example.com","password":"password123"}`))
	loggedIn := httptest.NewRecorder()
	service.Login(loggedIn, login)
	if loggedIn.Code != http.StatusOK {
		t.Fatalf("login status = %d, want %d", loggedIn.Code, http.StatusOK)
	}

	me := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	me.Header.Set("Authorization", "Bearer "+response.Token)
	meResponse := httptest.NewRecorder()
	service.Middleware(http.HandlerFunc(service.Me)).ServeHTTP(meResponse, me)
	if meResponse.Code != http.StatusOK {
		t.Fatalf("me status = %d, want %d", meResponse.Code, http.StatusOK)
	}
}

func TestRegisterRejectsDuplicateAndWeakPassword(t *testing.T) {
	service := testService(t)
	for _, body := range []string{
		`{"name":"Amina","email":"amina@example.com","password":"password123"}`,
		`{"name":"Other","email":"AMINA@example.com","password":"password123"}`,
	} {
		response := httptest.NewRecorder()
		service.Register(response, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body)))
		if response.Code != http.StatusCreated && response.Code != http.StatusConflict {
			t.Fatalf("unexpected registration status = %d", response.Code)
		}
	}
	weak := httptest.NewRecorder()
	service.Register(weak, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"name":"Short","email":"short@example.com","password":"short"}`)))
	if weak.Code != http.StatusBadRequest {
		t.Fatalf("weak password status = %d, want %d", weak.Code, http.StatusBadRequest)
	}
}
