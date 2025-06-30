package errz

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// validateAllJSONFiles validates all JSON files in a directory against the schema.
func validateAllJSONFiles(schemaPath, dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %w", dir, err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		jsonPath := filepath.Join(dir, entry.Name())
		if err := validateJSON(schemaPath, jsonPath); err != nil {
			return fmt.Errorf("validation failed for %s: %w", entry.Name(), err)
		}
	}

	return nil
}

// validateJSON validates a JSON file against a JSON Schema located at schemaPath.
// If validation fails, it returns a detailed error message with all issues found.
func validateJSON(schemaPath, jsonPath string) error {
	schemaLoader, err := loadFileAsReferenceLoader(schemaPath)
	if err != nil {
		return fmt.Errorf("cannot load schema %s: %w", schemaPath, err)
	}

	documentLoader, err := loadFileAsReferenceLoader(jsonPath)
	if err != nil {
		return fmt.Errorf("cannot load JSON file %s: %w", jsonPath, err)
	}

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to run validation: %w", err)
	}

	if !result.Valid() {
		var builder strings.Builder
		builder.WriteString("JSON validation failed:\n")

		for _, e := range result.Errors() {
			builder.WriteString("- ")
			builder.WriteString(e.String())
			builder.WriteRune('\n')
		}

		return fmt.Errorf("%s", builder.String())
	}

	return nil
}

// loadFileAsReferenceLoader converts a file path to a gojsonschema JSONLoader with error handling.
func loadFileAsReferenceLoader(path string) (gojsonschema.JSONLoader, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get absolute path: %w", err)
	}

	if _, err := os.Stat(abs); err != nil {
		return nil, fmt.Errorf("file not found: %w", err)
	}

	return gojsonschema.NewReferenceLoader("file:///" + filepath.ToSlash(abs)), nil
}
