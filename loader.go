package errorz

import (
	"encoding/json"
	"fmt"
	"os"
)

// LoadErrors loads and validates multiple JSON error definition files.
// It accepts the JSON Schema path and a slice of JSON file paths to load.
// Returns a combined map of error code to ErrorDefinition.
// Returns error if validation fails, file read fails, JSON unmarshal fails,
// or if duplicate error codes are found across files.
func LoadErrors(schemaPath string, jsonPath []string) (map[string]ErrorDefinition, error) {
	allErrors := make(map[string]ErrorDefinition)

	for _, path := range jsonPath {
		// Validate JSON file against schema
		if err := ValidateJSON(schemaPath, path); err != nil {
			return nil, fmt.Errorf("validation JSON failed for %s: %w", path, err)
		}

		// Read JSON file content
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", path, err)
		}

		// Parse JSON into map[string]ErrorDefinition
		var errorsMap map[string]ErrorDefinition
		if err := json.Unmarshal(fileContent, &errorsMap); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON %s: %w", path, err)
		}

		// Check for duplicates and merge into allErrors
		for code, errDef := range errorsMap {
			if _, exists := allErrors[code]; exists {
				return nil, fmt.Errorf("duplicate error code %s in file %s", code, path)
			}
			allErrors[code] = errDef
		}
	}

	return allErrors, nil
}
