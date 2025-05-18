package errorz_test

import (
	"strings"
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
)

func TestGenerateGoContent_Success(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM0001": {Code: "PM0001", Msg: "Payment failed", HTTPStatus: 400},
		"PM0002": {Code: "PM0002", Msg: "Timeout", HTTPStatus: 504},
	}

	code, err := errorz.GenerateGoContent("payment", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(code, "type Error struct") {
		t.Errorf("expected type definition in output")
	}

	if !strings.Contains(code, `PM0001 = Error{`) || !strings.Contains(code, `PM0002 = Error{`) {
		t.Errorf("expected error constants in output")
	}

	if !strings.Contains(code, "var Errors = []Error{") {
		t.Errorf("expected Errors slice in output")
	}

	if strings.Index(code, "PM0001") > strings.Index(code, "PM0002") {
		t.Errorf("error codes not sorted alphabetically")
	}
}

func TestGenerateGoContent_EmptyMap(t *testing.T) {
	code, err := errorz.GenerateGoContent("errorz", map[string]errorz.ErrorDefinition{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(code, "var Errors = []Error{") {
		t.Errorf("expected Errors slice even when empty")
	}
}

func TestGenerateGoContent_InvalidPackage(t *testing.T) {
	_, err := errorz.GenerateGoContent(" ", map[string]errorz.ErrorDefinition{})
	if err != errorz.ErrInvalidPackageName {
		t.Errorf("expected ErrInvalidPackageName, got %v", err)
	}
}

func TestGenerateGoContent_EscapedString(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM0003": {
			Code:       `PM0003"`,
			Msg:        "multi\nline \"message\"",
			HTTPStatus: 500,
		},
	}

	code, err := errorz.GenerateGoContent("payment", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(code, `Code: "PM0003\""`) {
		t.Errorf("expected escaped quote in Code")
	}
	if !strings.Contains(code, `Msg: "multi\nline \"message\""`) {
		t.Errorf("expected escaped content in Msg:\n%s", code)
	}
}
