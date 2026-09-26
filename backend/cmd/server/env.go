package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// loadDevelopmentEnv loads the nearest .env from the working directory or its
// parent (for go run from backend/). Injected environment variables always win.
// A missing file is allowed so deployments can use only injected configuration.
func loadDevelopmentEnv(directory string) error {
	for _, path := range []string{filepath.Join(directory, ".env"), filepath.Join(directory, "..", ".env")} {
		err := godotenv.Load(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			// Parser errors can contain secret values; never include them in logs.
			return fmt.Errorf("could not load %s: check file permissions and dotenv syntax", path)
		}
		return nil
	}
	return nil
}
