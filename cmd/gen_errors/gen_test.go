package main

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

	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to chdir: %v", err)
	}

	root, err := projectRoot()
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	expected, _ := filepath.EvalSymlinks(tmpDir)
	actual, _ := filepath.EvalSymlinks(root)
	if actual != expected {
		t.Errorf("expected %s, got %s", expected, actual)
	}
}

func TestProjectRoot_NotFound(t *testing.T) {
	tmpDir := t.TempDir()
	oldDir, _ := os.Getwd()
	defer os.Chdir(oldDir)
	os.Chdir(tmpDir)

	_, err := projectRoot()
	if err == nil || err.Error() != "project root not found (no go.mod)" {
		t.Errorf("expected not found error, got: %v", err)
	}
}
