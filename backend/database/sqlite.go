package database

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/001_users.sql
var usersSchema string

// Open creates the database directory and opens the SQLite database.
func Open(path string) (*sql.DB, error) {
	if path == ":memory:" {
		path = "file:suuq?mode=memory&cache=shared"
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create database directory: %w", err)
	}

	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	return db, nil
}

// Migrate creates the tables required by the application.
func Migrate(db *sql.DB) error {
	for _, statement := range strings.Split(usersSchema, ";") {
		if strings.TrimSpace(statement) == "" {
			continue
		}
		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("run database migration: %w", err)
		}
	}
	return nil
}
