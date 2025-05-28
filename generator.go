//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Generate(outputPath, outputDirPath string, errors map[string]ErrorDefinition) error {
	path := strings.ToLower(outputPath)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	// Write the errors_gen.go file.
	err := WriteGoFile(outputPath, errors)
	if err != nil {
		return fmt.Errorf("failed to write go content: %w", err)
	}

	domainGroups := make(map[string]map[string]ErrorDefinition)
	for code, def := range errors {
		domain := def.Domain
		if domain == "" {
			return fmt.Errorf("error code %q has empty domain", code)
		}

		if _, ok := domainGroups[domain]; !ok {
			domainGroups[domain] = make(map[string]ErrorDefinition)
		}

		domainGroups[domain][code] = def
	}

	for domain, group := range domainGroups {
		if err := WriteMarkdownFile(outputDirPath, domain, group); err != nil {
			return fmt.Errorf("failed to write markdown for domain %q: %w", domain, err)
		}
	}

	return nil
}
