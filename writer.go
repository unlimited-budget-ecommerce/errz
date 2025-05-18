package errorz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrEmptyPath = errors.New("path cannot be empty")

// make it easy for unit test
var MkdirAll = os.MkdirAll
var WriteToFileFunc = WriteToFile
var GenerateGoContentFunc = GenerateGoContent
var GenerateMarkdownContentFunc = GenerateMarkdownContent

// WriteGoFile generates a Go source file that defines an Error struct,
// individual error variables, and a slice containing all of them.
func WriteGoFile(outputDir, packageName string, errors map[string]ErrorDefinition) error {
	if err := MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create output dir: %w", err)
	}

	content, err := GenerateGoContentFunc(packageName, errors)
	if err != nil {
		return fmt.Errorf("failed to generate Go content: %w", err)
	}

	outputPath := filepath.Join(outputDir, fmt.Sprintf("%s.go", packageName))
	if err := WriteToFileFunc(outputPath, content); err != nil {
		return fmt.Errorf("failed to write Go file: %w", err)
	}

	return nil
}

// WriteMarkdown generates a Markdown (.md) documentation file containing
// a table of error definitions for a given domain.
func WriteMarkdown(outputDir, domain string, errors map[string]ErrorDefinition) error {
	if err := MkdirAll(outputDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create documentation directory: %w", err)
	}

	filename := filepath.Join(outputDir, fmt.Sprintf("%s.md", strings.ToLower(domain)))
	content := GenerateMarkdownContentFunc(domain, errors)

	if err := WriteToFileFunc(filename, content); err != nil {
		return fmt.Errorf("failed to write markdown file: %w", err)
	}

	return nil
}

// writeToFile writes content to the specified file path.
func WriteToFile(path, content string) error {
	if path == "" {
		return ErrEmptyPath
	}

	return os.WriteFile(path, []byte(content), 0644)
}
