package auth

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	errInvalidCredentials = errors.New("invalid credentials")
	errEmailExists        = errors.New("email already registered")
)

type Service struct {
	db     *sql.DB
	secret []byte
	now    func() time.Time
}

type user struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type credentials struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type tokenClaims struct {
	UserID    int64 `json:"sub"`
	ExpiresAt int64 `json:"exp"`
}

func NewService(db *sql.DB, secret string) *Service {
	return &Service{db: db, secret: []byte(secret), now: time.Now}
}

func (s *Service) Register(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	input.Name = strings.TrimSpace(input.Name)
	if input.Name == "" || !strings.Contains(input.Email, "@") || len(input.Password) < 8 {
		writeError(w, http.StatusBadRequest, "name, a valid email, and a password of at least 8 characters are required")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not secure password")
		return
	}
	result, err := s.db.Exec("INSERT INTO users (name, email, password_hash) VALUES (?, ?, ?)", input.Name, input.Email, string(hash))
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			writeError(w, http.StatusConflict, errEmailExists.Error())
			return
		}
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}
	userID, err := result.LastInsertId()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create user")
		return
	}

	createdUser := user{ID: userID, Name: input.Name, Email: input.Email}
	token, err := s.issueToken(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"token": token, "user": createdUser})
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	var input credentials
	if !decodeJSON(w, r, &input) {
		return
	}
	input.Email = strings.ToLower(strings.TrimSpace(input.Email))
	var userID int64
	var passwordHash, name string
	err := s.db.QueryRow("SELECT id, name, password_hash FROM users WHERE email = ?", input.Email).Scan(&userID, &name, &passwordHash)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(input.Password)) != nil {
		writeError(w, http.StatusUnauthorized, errInvalidCredentials.Error())
		return
	}

	token, err := s.issueToken(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not create session")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"token": token, "user": user{ID: userID, Name: name, Email: input.Email}})
}

func (s *Service) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "authentication required")
		return
	}
	var current user
	err := s.db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userID).Scan(&current.ID, &current.Name, &current.Email)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusUnauthorized, "user not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "could not load user")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"user": current})
}

func (s *Service) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Fields(r.Header.Get("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			writeError(w, http.StatusUnauthorized, "authentication required")
			return
		}
		claims, err := s.parseToken(parts[1])
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid or expired token")
			return
		}
		next.ServeHTTP(w, r.WithContext(withUserID(r.Context(), claims.UserID)))
	})
}

func (s *Service) issueToken(userID int64) (string, error) {
	claims := tokenClaims{UserID: userID, ExpiresAt: s.now().Add(24 * time.Hour).Unix()}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	return encodedPayload + "." + base64.RawURLEncoding.EncodeToString(s.sign([]byte(encodedPayload))), nil
}

func (s *Service) parseToken(token string) (tokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return tokenClaims{}, errInvalidCredentials
	}
	expected := s.sign([]byte(parts[0]))
	provided, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil || subtle.ConstantTimeCompare(expected, provided) != 1 {
		return tokenClaims{}, errInvalidCredentials
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return tokenClaims{}, errInvalidCredentials
	}
	var claims tokenClaims
	if json.Unmarshal(payload, &claims) != nil || claims.UserID < 1 || claims.ExpiresAt <= s.now().Unix() {
		return tokenClaims{}, errInvalidCredentials
	}
	return claims, nil
}

func (s *Service) sign(value []byte) []byte {
	hash := hmac.New(sha256.New, s.secret)
	hash.Write(value)
	return hash.Sum(nil)
}

type contextKey string

const userIDKey contextKey = "userID"

func withUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func UserID(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(userIDKey).(int64)
	return userID, ok
}

func decodeJSON(w http.ResponseWriter, r *http.Request, value any) bool {
	if json.NewDecoder(r.Body).Decode(value) != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return false
	}
	return true
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
