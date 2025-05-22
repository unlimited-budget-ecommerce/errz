//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectRoot_Success(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)
	if err != nil {
		t.Fatalf("failed to create go.mod: %v", err)
	}

	// Change to subdir
	subDir := filepath.Join(tmpDir, "sub")
	if err := os.Mkdir(subDir, 0755); err != nil {
		t.Fatalf("failed to create subdir: %v", err)
	}
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	if err := os.Chdir(subDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	root, err := ProjectRoot()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if root != tmpDir {
		t.Errorf("expected %s, got %s", tmpDir, root)
	}
}

func TestProjectRoot_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	_, err := ProjectRoot()
	if err == nil || err.Error() != "project root not found (no go.mod)" {
		t.Errorf("expected not found error, got: %v", err)
	}
}
