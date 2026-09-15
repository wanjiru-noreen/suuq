package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func clearEnv(t *testing.T, key string) {
	t.Helper()
	t.Setenv(key, "")
	if err := os.Unsetenv(key); err != nil {
		t.Fatal(err)
	}
}

func writeEnv(t *testing.T, directory, content string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(directory, ".env"), []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDevelopmentEnv(t *testing.T) {
	for _, fromBackend := range []bool{false, true} {
		name := "root"
		if fromBackend {
			name = "backend"
		}
		t.Run(name, func(t *testing.T) {
			clearEnv(t, "AUTH_SECRET")
			clearEnv(t, "DB_PATH")
			root := t.TempDir()
			writeEnv(t, root, "# development settings\nAUTH_SECRET=\"development test secret\"\nDB_PATH=./data/test.db\n")
			directory := root
			if fromBackend {
				directory = filepath.Join(root, "backend")
				if err := os.Mkdir(directory, 0700); err != nil {
					t.Fatal(err)
				}
			}
			if err := loadDevelopmentEnv(directory); err != nil {
				t.Fatal(err)
			}
			if os.Getenv("AUTH_SECRET") != "development test secret" || os.Getenv("DB_PATH") != "./data/test.db" {
				t.Fatal("configuration was not loaded")
			}
		})
	}
}

func TestLoadDevelopmentEnvPreservesInjectedValues(t *testing.T) {
	for _, value := range []string{"injected", ""} {
		t.Run("value="+value, func(t *testing.T) {
			t.Setenv("AUTH_SECRET", value)
			root := t.TempDir()
			writeEnv(t, root, "AUTH_SECRET=file-value\n")
			if err := loadDevelopmentEnv(root); err != nil {
				t.Fatal(err)
			}
			if os.Getenv("AUTH_SECRET") != value {
				t.Fatal("file overwrote injected environment")
			}
		})
	}
}

func TestLoadDevelopmentEnvPrefersCurrentDirectory(t *testing.T) {
	clearEnv(t, "AUTH_SECRET")
	root := t.TempDir()
	child := filepath.Join(root, "backend")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}
	writeEnv(t, root, "AUTH_SECRET=parent\n")
	writeEnv(t, child, "AUTH_SECRET=current\n")
	if err := loadDevelopmentEnv(child); err != nil {
		t.Fatal(err)
	}
	if os.Getenv("AUTH_SECRET") != "current" {
		t.Fatal("did not prefer current directory")
	}
}

func TestLoadDevelopmentEnvMissingFile(t *testing.T) {
	root := t.TempDir()
	child := filepath.Join(root, "backend")
	if err := os.Mkdir(child, 0700); err != nil {
		t.Fatal(err)
	}
	if err := loadDevelopmentEnv(child); err != nil {
		t.Fatal(err)
	}
}

func TestLoadDevelopmentEnvMalformedFileDoesNotLeakSecrets(t *testing.T) {
	root := t.TempDir()
	writeEnv(t, root, "AUTH_SECRET=\"private-value\n")
	err := loadDevelopmentEnv(root)
	if err == nil {
		t.Fatal("accepted malformed file")
	}
	if strings.Contains(err.Error(), "private-value") {
		t.Fatal("error leaked secret")
	}
}
