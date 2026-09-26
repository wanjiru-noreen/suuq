package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"suuq/database"
	"suuq/internal/auth"
	"suuq/routes"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func validateSecret(secret string) error {
	if len(strings.TrimSpace(secret)) < 32 {
		return errors.New("AUTH_SECRET must contain at least 32 bytes; generate it with openssl rand -hex 32")
	}
	if secret == "replace-with-a-random-secret-of-at-least-32-bytes" {
		return errors.New("AUTH_SECRET is the .env.example placeholder; generate a real secret with openssl rand -hex 32")
	}
	return nil
}

func run() error {
	if err := loadDevelopmentEnv("."); err != nil {
		return err
	}
	secret := os.Getenv("AUTH_SECRET")
	if err := validateSecret(secret); err != nil {
		return err
	}
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/suuq.db"
	}
	db, err := database.Open(dbPath)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := database.Migrate(db); err != nil {
		return err
	}
	server := &http.Server{
		Addr:              ":8080",
		Handler:           routes.SetupRoutes(auth.NewService(db, secret)),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    16 << 10,
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	result := make(chan error, 1)
	go func() { result <- server.ListenAndServe() }()
	log.Println("Suuq backend listening on :8080")
	select {
	case err := <-result:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		stop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			server.Close()
			return fmt.Errorf("shutdown server: %w", err)
		}
		if err := <-result; !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	}
}
