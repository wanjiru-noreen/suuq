package auth

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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
	for index, body := range []string{
		`{"name":"Amina","email":"amina@example.com","password":"password123"}`,
		`{"name":"Other","email":"AMINA@example.com","password":"password123"}`,
	} {
		response := httptest.NewRecorder()
		service.Register(response, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(body)))
		expected := http.StatusCreated
		if index == 1 {
			expected = http.StatusConflict
		}
		if response.Code != expected {
			t.Fatalf("unexpected registration status = %d", response.Code)
		}
	}
	weak := httptest.NewRecorder()
	service.Register(weak, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"name":"Short","email":"short@example.com","password":"short"}`)))
	if weak.Code != http.StatusBadRequest {
		t.Fatalf("weak password status = %d, want %d", weak.Code, http.StatusBadRequest)
	}
}

func TestRegisterRejectsInvalidInput(t *testing.T) {
	service := testService(t)
	cases := []struct {
		name, body string
		status     int
	}{
		{"invalid email", `{"name":"A","email":"a@@example.com","password":"password123"}`, 400},
		{"long password", `{"name":"A","email":"a@example.com","password":"` + strings.Repeat("a", 73) + `"}`, 400},
		{"unknown field", `{"name":"A","email":"a@example.com","password":"password123","admin":true}`, 400},
		{"trailing JSON", `{"name":"A","email":"a@example.com","password":"password123"} {}`, 400},
		{"oversized", `{"name":"` + strings.Repeat("a", 17000) + `"}`, 413},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			service.Register(response, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(tc.body)))
			if response.Code != tc.status {
				t.Fatalf("status = %d, want %d", response.Code, tc.status)
			}
		})
	}
	var count int
	if err := service.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Fatal("invalid registration created users")
	}
}

func TestMiddlewareRejectsInvalidTokens(t *testing.T) {
	service := testService(t)
	issuedAt := time.Now()
	service.now = func() time.Time { return issuedAt }
	token, err := service.issueToken(1)
	if err != nil {
		t.Fatal(err)
	}
	service.now = func() time.Time { return issuedAt.Add(24 * time.Hour) }
	for _, header := range []string{"", "Bearer malformed", "Bearer " + token, "Bearer " + token + "tampered"} {
		request := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
		request.Header.Set("Authorization", header)
		response := httptest.NewRecorder()
		service.Middleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { t.Fatal("invalid token reached handler") })).ServeHTTP(response, request)
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d", response.Code)
		}
	}
}

func TestLoginWrongPassword(t *testing.T) {
	service := testService(t)
	registered := httptest.NewRecorder()
	service.Register(registered, httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{"name":"A","email":"a@example.com","password":"password123"}`)))
	if registered.Code != http.StatusCreated {
		t.Fatal(registered.Body.String())
	}
	for _, email := range []string{"a@example.com", "missing@example.com"} {
		response := httptest.NewRecorder()
		service.Login(response, httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email":"`+email+`","password":"incorrect"}`)))
		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d", response.Code)
		}
		if response.Header().Get("Cache-Control") != "no-store" {
			t.Fatal("authentication response can be cached")
		}
	}
}
