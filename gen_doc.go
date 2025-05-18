package errorz

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"unicode"
)

var (
	titleCache sync.Map // map[string]string
)

// GenerateMarkdownContent builds Markdown content for a given domain and its errors.
func GenerateMarkdownContent(domain string, errors map[string]ErrorDefinition) string {
	// Sort error codes alphabetically for consistent ordering
	var codes []string
	for code := range errors {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	var builder strings.Builder
	builder.WriteString(NormalizeMarkdownTitle(domain))

	// Write Markdown header
	builder.WriteString("| Code | Msg | HTTP | Category | Severity | Retryable |\n")
	builder.WriteString("|------|-----|------|----------|----------|-----------|\n")

	// Write each error row
	for _, code := range codes {
		errDef := errors[code]
		builder.WriteString(fmt.Sprintf(
			"| %s | %s | %d | %s | %s | %t |\n",
			errDef.Code,
			EscapeMarkdown(errDef.Msg),
			errDef.HTTPStatus,
			errDef.Category,
			errDef.Severity,
			errDef.IsRetryable,
		))
	}

	return builder.String()
}

// EscapeMarkdown escapes special characters for Markdown table rendering.
func EscapeMarkdown(text string) string {
	return strings.ReplaceAll(text, "|", "\\|")
}

// NormalizeMarkdownTitle formats the domain into a Markdown header with each
// hyphen-separated part capitalized (first letter uppercase, rest lowercase).
// Results are cached for improved performance and concurrency safety.
func NormalizeMarkdownTitle(domain string) string {
	if cached, ok := titleCache.Load(domain); ok {
		return cached.(string)
	}

	parts := strings.Split(domain, "-")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}

		p = strings.ToLower(p)
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}

	result := "# " + strings.Join(parts, "-") + " Errors\n\n"
	titleCache.Store(domain, result)

	return result
}
