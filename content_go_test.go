//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGoContent_Basic(t *testing.T) {
	defs := map[string]ErrorDefinition{
		"TT0001": {
			Domain: "test",
			Code:   "TT0001",
			Msg:    "Something went wrong",
			Cause:  "Unknown",
		},
	}

	code, err := GenerateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, `type Error struct {`)
	assert.Contains(t, code, "func (e *Error) Error() string")
	assert.Contains(t, code, `TT0001 = &Error{`)
	assert.Contains(t, code, `Domain: "test"`)
	assert.Contains(t, code, `Code: "TT0001"`)
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

func TestGenerateGoContent_ErrorMethodIncluded(t *testing.T) {
	defs := map[string]ErrorDefinition{
		"XX0001": {
			Domain: "x",
			Code:   "XX0001",
			Msg:    "msg",
			Cause:  "cause",
		},
	}

	code, err := GenerateGoContent(defs)
	assert.NoError(t, err)
	assert.Contains(t, code, "func (e *Error) Error() string")
}
