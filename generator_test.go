//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	outputGoFile := filepath.Join(tmpDir, "errors_gen.go")

	errors := map[string]ErrorDefinition{
		"UR0001": {
			Domain:      "user",
			Code:        "USER_NOT_FOUND",
			Msg:         "User not found",
			Cause:       "User ID missing",
			Severity:    "low",
			IsRetryable: false,
		},
		"OR0001": {
			Domain:      "order",
			Code:        "ORDER_FAILED",
			Msg:         "Order could not be completed",
			Cause:       "Payment issue",
			Severity:    "critical",
			IsRetryable: true,
		},
	}

	err := Generate(outputGoFile, tmpDir, errors)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Check Go file created
	if _, err := os.Stat(outputGoFile); err != nil {
		t.Errorf("Go file not created: %v", err)
	}

	// Check Markdown files created
	expectedFiles := []string{
		filepath.Join(tmpDir, "user", "user.md"),
		filepath.Join(tmpDir, "order", "order.md"),
	}
	for _, path := range expectedFiles {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("Markdown file missing: %s", path)
		}
	}
}

func TestGenerate_EmptyOutputPath(t *testing.T) {
	err := Generate("", t.TempDir(), map[string]ErrorDefinition{})
	if err == nil || err.Error() != "failed to write go content: output file path cannot be empty" {
		t.Errorf("Expected output file path error, got: %v", err)
	}
}

func TestGenerate_EmptyMarkdownOutputDir(t *testing.T) {
	err := Generate(t.TempDir()+"/go.go", "", map[string]ErrorDefinition{
		"X": {Code: "X", Domain: "abc"},
	})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected markdown directory error, got: %v", err)
	}
}

func TestGenerate_EmptyDomainInError(t *testing.T) {
	err := Generate(t.TempDir()+"/go.go", t.TempDir()+"/doc", map[string]ErrorDefinition{
		"X": {Code: "X", Domain: ""},
	})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected domain error, got: %v", err)
	}
}
