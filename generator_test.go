package errz

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate_Success(t *testing.T) {
	tmpDir := t.TempDir()
	outputGoFile := filepath.Join(tmpDir, "errors_gen.go")

	errors := map[string]Error{
		"UR0001": {
			Domain: "user",
			Code:   "USER_NOT_FOUND",
			Msg:    "User not found",
			Cause:  "User ID missing",
		},
		"OR0001": {
			Domain: "order",
			Code:   "ORDER_FAILED",
			Msg:    "Order could not be completed",
			Cause:  "Payment issue",
		},
	}

	err := generate(outputGoFile, tmpDir, errors)
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
	err := generate("", t.TempDir(), map[string]Error{})
	if err == nil || err.Error() != "failed to write go content: output file path cannot be empty" {
		t.Errorf("Expected output file path error, got: %v", err)
	}
}

func TestGenerate_EmptyMarkdownOutputDir(t *testing.T) {
	err := generate(t.TempDir()+"/go.go", "", map[string]Error{
		"X": {Code: "X", Domain: "abc"},
	})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected markdown directory error, got: %v", err)
	}
}

func TestGenerate_EmptyDomainInError(t *testing.T) {
	err := generate(t.TempDir()+"/go.go", t.TempDir()+"/doc", map[string]Error{
		"X": {Code: "X", Domain: ""},
	})
	if err == nil || err.Error() == "" {
		t.Errorf("Expected domain error, got: %v", err)
	}
}

func TestGenerateGoContent_Basic(t *testing.T) {
	defs := map[string]Error{
		"TT0001": {
			Domain: "test",
			Code:   "TT0001",
			Msg:    "Something went wrong",
			Cause:  "Unknown",
		},
	}

	code, err := generateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, `type Error struct {`)
	assert.Contains(t, code, "func (e *Error) Error() string")
	assert.Contains(t, code, `TT0001 = &Error{`)
	assert.Contains(t, code, `Domain: "test"`)
	assert.Contains(t, code, `Code: "TT0001"`)
}

func TestGenerateGoContent_MultipleErrorsSorted(t *testing.T) {
	defs := map[string]Error{
		"ZE0001": {Code: "ZE0001", Msg: "z"},
		"AE0001": {Code: "AE0001", Msg: "a"},
	}

	code, err := generateGoContent(defs)
	assert.NoError(t, err)

	zIndex := strings.Index(code, "ZE0001 = &Error{")
	aIndex := strings.Index(code, "AE0001 = &Error{")
	assert.True(t, aIndex < zIndex, "AE0001 should come before ZE0001")
}

func TestGenerateGoContent_EscapeCharacters(t *testing.T) {
	defs := map[string]Error{
		"TT0001": {
			Code: "TT0001",
			Msg:  `quote " and newline \n`,
		},
	}

	code, err := generateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, `quote \" and newline \\n`)
}

func TestGenerateGoContent_EmptyInput(t *testing.T) {
	code, err := generateGoContent(map[string]Error{})
	assert.Error(t, err)
	assert.Empty(t, code)
	assert.EqualError(t, err, "no error definitions provided")
}

func TestGenerateGoContent_ErrorMethodIncluded(t *testing.T) {
	defs := map[string]Error{
		"XX0001": {
			Domain: "x",
			Code:   "XX0001",
			Msg:    "msg",
			Cause:  "cause",
		},
	}

	code, err := generateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, "func (e *Error) Error() string")
}

// titleCacheReset clears the title cache (for test use only)
func titleCacheReset() {
	titleCache = sync.Map{}
}

func TestGenerateMarkdownContent_ValidInput(t *testing.T) {
	titleCacheReset()
	defer titleCacheReset()

	errorsMap := map[string]Error{
		"ERR001": {
			Code:  "ERR001",
			Msg:   "Invalid input | bad format",
			Cause: "Input contained unexpected value",
		},
		"ERR002": {
			Code:  "ERR002",
			Msg:   "Timeout `network`",
			Cause: "Service did not respond",
		},
	}

	md, err := generateMarkdownContent("core-api", errorsMap)
	assert.NoError(t, err)
	assert.Contains(t, md, "# Core-Api Errors")
	assert.Contains(t, md, "| ERR001 | Invalid input \\| bad format |")
	assert.Contains(t, md, "- **Cause**: Input contained unexpected value")
	assert.Contains(t, md, "## ERR002")
	assert.Contains(t, md, "Timeout \\`network\\`")
}

func TestGenerateMarkdownContent_InvalidDomain(t *testing.T) {
	titleCacheReset()
	defer titleCacheReset()

	_, err := generateMarkdownContent("bad domain", map[string]Error{})
	assert.ErrorIs(t, err, errInvalidDomainName)

	_, err = generateMarkdownContent(" ", map[string]Error{})
	assert.ErrorIs(t, err, errInvalidDomainName)
}

func TestGenerateMarkdownContent_EmptyErrors(t *testing.T) {
	titleCacheReset()
	defer titleCacheReset()

	md, err := generateMarkdownContent("example", map[string]Error{})
	assert.Error(t, err)
	assert.Empty(t, md)
	assert.EqualError(t, err, "no error definitions provided for markdown generation")
}

func TestGenerateMarkdownContent_Sorting(t *testing.T) {
	titleCacheReset()
	defer titleCacheReset()

	errorsMap := map[string]Error{
		"B": {Code: "B"},
		"A": {Code: "A"},
	}

	md, err := generateMarkdownContent("domain", errorsMap)
	assert.NoError(t, err)
	firstIdx := strings.Index(md, "## A")
	secondIdx := strings.Index(md, "## B")
	assert.Less(t, firstIdx, secondIdx)
}

func TestNormalizeMarkdownTitle_CachesCorrectly(t *testing.T) {
	titleCacheReset()
	defer titleCacheReset()

	title1 := normalizeMarkdownTitle("my-service")
	title2 := normalizeMarkdownTitle("my-service")

	assert.Equal(t, "# My-Service Errors\n\n", title1)
	assert.Equal(t, title1, title2) // from cache
}

func TestWriteToFile_Success(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "output.txt")
	content := "hello test"

	err := writeToFile(tmpFile, content)
	require.NoError(t, err)

	// Check file content
	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Equal(t, content, string(data))
}

func TestWriteToFile_EmptyPath(t *testing.T) {
	err := writeToFile("", "data")
	require.ErrorIs(t, err, errEmptyPath)
}

func TestWriteGoFile_Success(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "errors.go")

	err := writeGoFile(tmpFile, map[string]Error{
		"TEST_CODE": {
			Code:  "TEST_CODE",
			Msg:   "This is a test error",
			Cause: "Just testing",
		},
	})

	require.NoError(t, err)

	data, err := os.ReadFile(tmpFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "TEST_CODE")
}

func TestWriteGoFile_EmptyPath(t *testing.T) {
	err := writeGoFile("", nil)
	require.ErrorIs(t, err, errEmptyFile)
}

func TestWriteMarkdownFile_Success(t *testing.T) {
	tmpDir := t.TempDir()
	domain := "test-domain"

	err := writeMarkdownFile(tmpDir, domain, map[string]Error{
		"TEST_MARKDOWN": {
			Code:  "TEST_MARKDOWN",
			Msg:   "Markdown message",
			Cause: "Some cause",
		},
	})

	require.NoError(t, err)

	mdFile := filepath.Join(tmpDir, strings.ToLower(domain), strings.ToLower(domain)+".md")
	data, err := os.ReadFile(mdFile)
	require.NoError(t, err)
	require.Contains(t, string(data), "Markdown message")
	require.Contains(t, string(data), "# Test-Domain Errors")
}

func TestWriteMarkdownFile_EmptyDir(t *testing.T) {
	err := writeMarkdownFile("", "domain", nil)
	require.ErrorIs(t, err, errEmptyDir)
}
