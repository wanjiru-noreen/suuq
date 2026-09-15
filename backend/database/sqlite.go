package database

import (
	"database/sql"
	"embed"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations/*.sql
var migrations embed.FS

// Open creates the database directory and opens the SQLite database.
func Open(path string) (*sql.DB, error) {
	var dsn string
	if path == ":memory:" {
		// Each Open owns its database; one connection preserves its lifetime.
		dsn = ":memory:?_foreign_keys=on&_busy_timeout=5000"
	} else {
		absolute, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("resolve database path: %w", err)
		}
		if err := os.MkdirAll(filepath.Dir(absolute), 0o700); err != nil {
			return nil, fmt.Errorf("create database directory: %w", err)
		}
		uri := url.URL{Scheme: "file", Path: absolute}
		dsn = uri.String() + "?_foreign_keys=on&_busy_timeout=5000"
	}
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	// Serialize SQLite writes and keep in-memory databases alive.
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
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

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin migrations: %w", err)
	}
	defer tx.Rollback()
	for _, file := range files {
		data, err := migrations.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", file, err)
		}

		for _, statement := range strings.Split(string(data), ";") {
			if strings.TrimSpace(statement) == "" {
				continue
			}

			if _, err := tx.Exec(statement); err != nil {
				return fmt.Errorf("run migration %s: %w", file, err)
			}
		}
	}

	return tx.Commit()
}
