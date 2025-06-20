//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TitleCacheReset clears the title cache (for test use only)
func TitleCacheReset() {
	titleCache = sync.Map{}
}

func TestGenerateMarkdownContent_ValidInput(t *testing.T) {
	TitleCacheReset()
	defer TitleCacheReset()

	errorsMap := map[string]ErrorDefinition{
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

	md, err := GenerateMarkdownContent("core-api", errorsMap)
	assert.NoError(t, err)
	assert.Contains(t, md, "# Core-Api Errors")
	assert.Contains(t, md, "| ERR001 | Invalid input \\| bad format |")
	assert.Contains(t, md, "- **Cause**: Input contained unexpected value")
	assert.Contains(t, md, "## ERR002")
	assert.Contains(t, md, "Timeout \\`network\\`")
}

func TestGenerateMarkdownContent_InvalidDomain(t *testing.T) {
	TitleCacheReset()
	defer TitleCacheReset()

	_, err := GenerateMarkdownContent("bad domain", map[string]ErrorDefinition{})
	assert.ErrorIs(t, err, ErrInvalidDomainName)

	_, err = GenerateMarkdownContent(" ", map[string]ErrorDefinition{})
	assert.ErrorIs(t, err, ErrInvalidDomainName)
}

func TestGenerateMarkdownContent_EmptyErrors(t *testing.T) {
	TitleCacheReset()
	defer TitleCacheReset()

	md, err := GenerateMarkdownContent("example", map[string]ErrorDefinition{})
	assert.Error(t, err)
	assert.Empty(t, md)
	assert.EqualError(t, err, "no error definitions provided")
}

func TestGenerateMarkdownContent_Sorting(t *testing.T) {
	TitleCacheReset()
	defer TitleCacheReset()

	errorsMap := map[string]ErrorDefinition{
		"B": {Code: "B"},
		"A": {Code: "A"},
	}

	md, err := GenerateMarkdownContent("domain", errorsMap)
	assert.NoError(t, err)
	firstIdx := strings.Index(md, "## A")
	secondIdx := strings.Index(md, "## B")
	assert.Less(t, firstIdx, secondIdx)
}

func TestNormalizeMarkdownTitle_CachesCorrectly(t *testing.T) {
	TitleCacheReset()
	defer TitleCacheReset()

	title1 := NormalizeMarkdownTitle("my-service")
	title2 := NormalizeMarkdownTitle("my-service")

	assert.Equal(t, "# My-Service Errors\n\n", title1)
	assert.Equal(t, title1, title2) // from cache
}
