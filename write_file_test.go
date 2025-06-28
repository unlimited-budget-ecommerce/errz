package errz

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

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

	err := writeGoFile(tmpFile, map[string]ErrorDefinition{
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

	err := writeMarkdownFile(tmpDir, domain, map[string]ErrorDefinition{
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
