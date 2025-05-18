package errorz_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
)

// Helpers to patch functions inside errorz package (using closures)
func resetPatches() {
	errorz.MkdirAll = os.MkdirAll
	errorz.WriteToFileFunc = errorz.WriteToFile
	errorz.GenerateGoContentFunc = errorz.GenerateGoContent
	errorz.GenerateMarkdownContentFunc = errorz.GenerateMarkdownContent
}

func TestWriteGoFile(t *testing.T) {
	defer resetPatches()

	tempDir := t.TempDir() // use isolated temp directory

	// Case: MkdirAll error
	errorz.MkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("mkdir error")
	}
	err := errorz.WriteGoFile(tempDir, "pkg", nil)
	if err == nil || err.Error() != "failed to create output dir: mkdir error" {
		t.Fatalf("expected mkdir error, got %v", err)
	}

	// Case: GenerateGoContent error
	errorz.MkdirAll = os.MkdirAll
	errorz.GenerateGoContentFunc = func(pkg string, errs map[string]errorz.ErrorDefinition) (string, error) {
		return "", errors.New("gen content error")
	}
	err = errorz.WriteGoFile(tempDir, "pkg", nil)
	if err == nil || err.Error() != "failed to generate Go content: gen content error" {
		t.Fatalf("expected generate content error, got %v", err)
	}

	// Case: WriteToFile error
	errorz.GenerateGoContentFunc = func(pkg string, errs map[string]errorz.ErrorDefinition) (string, error) {
		return "go content", nil
	}
	errorz.WriteToFileFunc = func(path, content string) error {
		return errors.New("write file error")
	}
	err = errorz.WriteGoFile(tempDir, "pkg", nil)
	if err == nil || err.Error() != "failed to write Go file: write file error" {
		t.Fatalf("expected write file error, got %v", err)
	}

	// Case: Success
	errorz.WriteToFileFunc = func(path, content string) error {
		return nil
	}
	err = errorz.WriteGoFile(tempDir, "pkg", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteMarkdown(t *testing.T) {
	defer resetPatches()

	tempDir := t.TempDir()

	// Case: MkdirAll returns an error
	errorz.MkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("mkdir error")
	}
	err := errorz.WriteMarkdown(tempDir, "domain", nil)
	if err == nil || err.Error() != "failed to create documentation directory: mkdir error" {
		t.Fatalf("expected mkdir error, got %v", err)
	}

	// Case: WriteToFile returns an error
	errorz.MkdirAll = os.MkdirAll
	errorz.GenerateMarkdownContentFunc = func(domain string, errs map[string]errorz.ErrorDefinition) string {
		return "markdown content"
	}
	errorz.WriteToFileFunc = func(path, content string) error {
		return errors.New("write file error")
	}
	err = errorz.WriteMarkdown(tempDir, "domain", nil)
	if err == nil || err.Error() != "failed to write markdown file: write file error" {
		t.Fatalf("expected write file error, got %v", err)
	}

	// Case: Success path
	errorz.WriteToFileFunc = func(path, content string) error {
		return nil
	}
	err = errorz.WriteMarkdown(tempDir, "domain", nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWriteToFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	err := errorz.WriteToFile(filePath, "hello world")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected content 'hello world', got %s", string(data))
	}
}

func TestWriteToFile_EmptyPath(t *testing.T) {
	err := errorz.WriteToFile("", "content")
	if err != errorz.ErrEmptyPath {
		t.Errorf("expected ErrEmptyPath, got %v", err)
	}
}

func TestWriteToFile_NoPermission(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping on CI where root fs might be writable")
	}

	err := errorz.WriteToFile("/root/forbidden.txt", "test")
	if err == nil {
		t.Fatal("expected error when writing to /root, got nil")
	}
}

func TestWriteToFile_PathIsDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	err := errorz.WriteToFile(tmpDir, "should fail")
	if err == nil {
		t.Fatal("expected error when writing to a directory, got nil")
	}
}
