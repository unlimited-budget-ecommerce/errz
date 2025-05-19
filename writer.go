package errorz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var ErrEmptyPackageName = errors.New("package name cannot be empty")
var ErrEmptyDomain = errors.New("domain cannot be empty")
var ErrEmptyPath = errors.New("path cannot be empty")

// make it easy for unit test
var MkdirAll = os.MkdirAll
var WriteToFileFunc = WriteToFile
var GenerateGoContentFunc = GenerateGoContent
var GenerateMarkdownContentFunc = GenerateMarkdownContent

// WriteGoFile generates a Go source file that defines an Error struct,
// individual error variables, and a slice containing all of them.
func WriteGoFile(outputDir, packageName string, errors map[string]ErrorDefinition) error {
	if packageName == "" {
		return ErrEmptyPackageName
	}

	// Create subdirectory by package name
	packageDir := filepath.Join(outputDir, strings.ToLower(packageName))
	if err := MkdirAll(packageDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create package directory: %w", err)
	}

	content, err := GenerateGoContentFunc(packageName, errors)
	if err != nil {
		return fmt.Errorf("failed to generate Go content: %w", err)
	}

	outputPath := filepath.Join(packageDir, fmt.Sprintf("%s.go", packageName))
	if err := WriteToFileFunc(outputPath, content); err != nil {
		return fmt.Errorf("failed to write Go file: %w", err)
	}

	return nil
}

// WriteMarkdown generates a Markdown (.md) documentation file containing
// a table of error definitions for a given domain.
func WriteMarkdown(outputDir, domain string, errors map[string]ErrorDefinition) error {
	if domain == "" {
		return ErrEmptyDomain
	}

	// Create subdirectory by domain name
	domainDir := filepath.Join(outputDir, strings.ToLower(domain))
	if err := MkdirAll(domainDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create documentation directory: %w", err)
	}

	filename := filepath.Join(domainDir, fmt.Sprintf("%s.md", strings.ToLower(domain)))
	content, err := GenerateMarkdownContentFunc(domain, errors)
	if err != nil {
		return fmt.Errorf("failed to generate markdown content: %w", err)
	}

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
