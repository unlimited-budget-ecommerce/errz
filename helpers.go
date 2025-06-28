//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"fmt"
	"os"
	"path/filepath"
)

func ProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get working directory: %w", err)
	}

	for !fileExists(filepath.Join(dir, "go.mod")) && dir != "/" {
		dir = filepath.Dir(dir)
	}

	if dir == "/" {
		return "", fmt.Errorf("project root not found (no go.mod)")
	}

	return dir, nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
