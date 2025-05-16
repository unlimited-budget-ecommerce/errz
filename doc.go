package errorz

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func WriteMarkdown(outputDir string, domain string, errors map[string]ErrorDefinition) error {
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create doc dir: %w", err)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("%s.md", strings.ToLower(domain)))
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create markdown file: %w", err)
	}
	defer file.Close()

	// Sort by code for consistency
	codes := make([]string, 0, len(errors))
	for code := range errors {
		codes = append(codes, code)
	}
	sort.Strings(codes)

	titleCaser := cases.Title(language.English)

	builder := &strings.Builder{}
	builder.WriteString(fmt.Sprintf("# %s Errors\n\n", titleCaser.String(domain)))
	builder.WriteString("| Code | Msg | HTTP | Category | Severity | Retryable |\n")
	builder.WriteString("|------|-----|------|----------|----------|-----------|\n")

	for _, code := range codes {
		errDef := errors[code]
		builder.WriteString(fmt.Sprintf("| %s | %s | %d | %s | %s | %t |\n",
			errDef.Code,
			escapeMarkdown(errDef.Msg),
			errDef.HTTPStatus,
			errDef.Category,
			errDef.Severity,
			errDef.IsRetryable,
		))
	}

	if _, err := file.WriteString(builder.String()); err != nil {
		return fmt.Errorf("failed to write markdown: %w", err)
	}

	fmt.Printf("Generated: %s\n", filename)
	return nil
}

func escapeMarkdown(text string) string {
	// Escape pipe character used in table columns
	return strings.ReplaceAll(text, "|", "\\|")
}
