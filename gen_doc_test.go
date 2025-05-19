package errorz_test

import (
	"strings"
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
)

func TestGenerateMarkdownContent_SingleError(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM0001": {
			Code:        "PM0001",
			Msg:         "payment failed",
			HTTPStatus:  400,
			Category:    "business",
			Severity:    "medium",
			IsRetryable: false,
		},
	}

	md, err := errorz.GenerateMarkdownContent("payment", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "# Payment Errors") {
		t.Errorf("missing title header: %s", md)
	}
	if !strings.Contains(md, "| PM0001 | payment failed | 400 | business | medium | false |") {
		t.Errorf("missing or incorrect table row: %s", md)
	}
}

func TestGenerateMarkdownContent_MultipleErrors_Sorted(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM0002": {
			Code:        "PM0002",
			Msg:         "timeout",
			HTTPStatus:  408,
			Category:    "timeout",
			Severity:    "low",
			IsRetryable: true,
		},
		"PM0001": {
			Code:        "PM0001",
			Msg:         "payment failed",
			HTTPStatus:  400,
			Category:    "business",
			Severity:    "medium",
			IsRetryable: false,
		},
	}

	md, err := errorz.GenerateMarkdownContent("payment", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	lines := strings.Split(md, "\n")
	// The first row of data should be PM0001 (sorted)
	found := false
	for _, line := range lines {
		if strings.HasPrefix(line, "| PM0001 ") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected sorted output with PM0001 first:\n%s", md)
	}
}

func TestGenerateMarkdownContent_WithPipesInMessage(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM0003": {
			Code:        "PM0003",
			Msg:         "invalid | character",
			HTTPStatus:  422,
			Category:    "validation",
			Severity:    "high",
			IsRetryable: false,
		},
	}

	md, err := errorz.GenerateMarkdownContent("payment", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "invalid \\| character") {
		t.Errorf("expected pipe character to be escaped:\n%s", md)
	}
}

func TestGenerateMarkdownContent_EmptyMap(t *testing.T) {
	md, err := errorz.GenerateMarkdownContent("billing", map[string]errorz.ErrorDefinition{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "# Billing Errors") {
		t.Errorf("missing header for empty error list:\n%s", md)
	}

	if !strings.Contains(md, "| Code |") {
		t.Errorf("missing markdown table header:\n%s", md)
	}
}

func TestGenerateMarkdownContent_DomainWithUnicodeAndCase(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM9999": {
			Code:        "PM9999",
			Msg:         "something went wrong",
			HTTPStatus:  500,
			Category:    "internal",
			Severity:    "critical",
			IsRetryable: false,
		},
	}

	md, err := errorz.GenerateMarkdownContent("üSêr-Dømãïn", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedHeader := "# Üsêr-Dømãïn Errors"
	if !strings.Contains(md, expectedHeader) {
		t.Errorf("unexpected markdown header casing: got\n%s\nexpected header:\n%s", md, expectedHeader)
	}
}

func TestGenerateMarkdownContent_MissingOptionalFields(t *testing.T) {
	errors := map[string]errorz.ErrorDefinition{
		"PM1234": {
			Code:        "PM1234",
			Msg:         "some error occurred",
			HTTPStatus:  500,
			Category:    "internal",
			Severity:    "high",
			IsRetryable: true,
			// Missing: Solution, Tags
		},
	}

	md, err := errorz.GenerateMarkdownContent("core", errors)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(md, "| PM1234 | some error occurred | 500 | internal | high | true |") {
		t.Errorf("expected table row with missing optional fields to be valid: %s", md)
	}
}

func TestEscapeMarkdown(t *testing.T) {
	input := "value with | pipe"
	expected := "value with \\| pipe"

	if out := errorz.EscapeMarkdown(input); out != expected {
		t.Errorf("expected '%s', got '%s'", expected, out)
	}
}

func TestNormalizeMarkdownTitle(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		want   string
	}{
		{
			name:   "single lowercase word",
			domain: "payment",
			want:   "# Payment Errors\n\n",
		},
		{
			name:   "single uppercase word",
			domain: "PAYMENT",
			want:   "# Payment Errors\n\n",
		},
		{
			name:   "hyphenated lowercase words",
			domain: "user-domain",
			want:   "# User-Domain Errors\n\n",
		},
		{
			name:   "hyphenated mixed case words",
			domain: "User-Domain",
			want:   "# User-Domain Errors\n\n",
		},
		{
			name:   "empty domain",
			domain: "",
			want:   "#  Errors\n\n",
		},
		{
			name:   "domain with empty parts",
			domain: "user--domain",
			want:   "# User--Domain Errors\n\n",
		},
		{
			name:   "unicode characters",
			domain: "üsêr-Dømãïn",
			want:   "# Üsêr-Dømãïn Errors\n\n",
		},
		{
			name:   "single character parts",
			domain: "a-b-c",
			want:   "# A-B-C Errors\n\n",
		},
		{
			name:   "part with empty string",
			domain: "-domain",
			want:   "# -Domain Errors\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := errorz.NormalizeMarkdownTitle(tt.domain)
			if got != tt.want {
				t.Errorf("normalizeMarkdownTitle(%q) = %q; want %q", tt.domain, got, tt.want)
			}
		})
	}
}
