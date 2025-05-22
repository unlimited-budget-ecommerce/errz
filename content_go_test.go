//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGoContent_Basic(t *testing.T) {
	defs := map[string]ErrorDefinition{
		"TT0001": {
			Code:        "TT0001",
			Msg:         "Something went wrong",
			Cause:       "Unknown",
			HTTPStatus:  500,
			Category:    "Internal",
			Severity:    "High",
			IsRetryable: false,
			Solution:    "Try again later",
			Tags:        []string{"internal", "retry"},
		},
	}

	code, err := GenerateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, `var ErrorMap = map[string]*Error`)
	assert.Contains(t, code, `TT0001 = &Error{`)
	assert.Contains(t, code, `Code: "TT0001"`)
	assert.Contains(t, code, `Tags: []string{`)
}

func TestGenerateGoContent_MultipleErrorsSorted(t *testing.T) {
	defs := map[string]ErrorDefinition{
		"ZE0001": {Code: "ZE0001", Msg: "z"},
		"AE0001": {Code: "AE0001", Msg: "a"},
	}

	code, err := GenerateGoContent(defs)
	assert.NoError(t, err)

	zIndex := strings.Index(code, "ZE0001 = &Error{")
	aIndex := strings.Index(code, "AE0001 = &Error{")
	assert.True(t, aIndex < zIndex, "AE0001 should come before ZE0001")
}

func TestGenerateGoContent_EscapeCharacters(t *testing.T) {
	defs := map[string]ErrorDefinition{
		"TT0001": {
			Code: "TT0001",
			Msg:  `quote " and newline \n`,
		},
	}

	code, err := GenerateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, `quote \" and newline \\n`)
}

func TestGenerateGoContent_EmptyInput(t *testing.T) {
	code, err := GenerateGoContent(map[string]ErrorDefinition{})
	assert.Error(t, err)
	assert.Empty(t, code)
	assert.EqualError(t, err, "no error definitions provided")
}
