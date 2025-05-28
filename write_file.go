//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyPath = errors.New("path cannot be empty")
	ErrEmptyFile = errors.New("output file path cannot be empty")
	ErrEmptyDir  = errors.New("directory path cannot be empty")
)

func WriteGoFile(outputPath string, errors map[string]ErrorDefinition) error {
	if strings.TrimSpace(outputPath) == "" {
		return ErrEmptyFile
	}

	content, err := GenerateGoContent(errors)
	if err != nil {
		return fmt.Errorf("failed to generate Go content: %w", err)
	}

	if err := WriteToFile(outputPath, content); err != nil {
		return fmt.Errorf("failed to write Go file: %w", err)
	}

	return nil
}

func WriteMarkdownFile(outputDirPath, domain string, errors map[string]ErrorDefinition) error {
	if strings.TrimSpace(outputDirPath) == "" {
		return ErrEmptyDir
	}

	domainLower := strings.ToLower(domain)
	domainDir := filepath.Join(outputDirPath, domainLower)

	if err := os.MkdirAll(domainDir, 0755); err != nil {
		return fmt.Errorf("failed to create documentation sub directory: %w", err)
	}

	filename := filepath.Join(domainDir, fmt.Sprintf("%s.md", domainLower))
	content, err := GenerateMarkdownContent(domain, errors)
	if err != nil {
		return fmt.Errorf("failed to generate markdown content: %w", err)
	}

	if err := WriteToFile(filename, content); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}

// writeToFile writes content to the specified file path.
func WriteToFile(path, content string) error {
	if strings.TrimSpace(path) == "" {
		return ErrEmptyPath
	}

	return os.WriteFile(path, []byte(content), 0644)
}
