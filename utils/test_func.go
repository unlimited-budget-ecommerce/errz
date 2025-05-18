// JUST FOR UNIT TEST
package utils_test

import (
	"os"
	"path/filepath"
	"testing"
)

func WriteTempFile(t *testing.T, name, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	return path
}
