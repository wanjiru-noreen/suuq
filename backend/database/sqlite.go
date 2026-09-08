package database

import (
	"database/sql"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrations embed.FS

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

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}

	return db, nil
}

// Migrate creates the tables required by the application.
func Migrate(db *sql.DB) error {
	files := []string{
		"migrations/001_users.sql",
		"migrations/002_businesses.sql",
		"migrations/003_products.sql",
		"migrations/004_stock_updates.sql",
		"migrations/005_debtors.sql",
		"migrations/006_creditors.sql",
		"migrations/007_transactions.sql",
		"migrations/008_transaction_items.sql",
	}

	for _, file := range files {
		data, err := migrations.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		for _, statement := range strings.Split(string(data), ";") {
			if strings.TrimSpace(statement) == "" {
				continue
			}

			if _, err := db.Exec(statement); err != nil {
				return fmt.Errorf("run migration %s: %w", file, err)
			}
		}
	}

	return nil
}
