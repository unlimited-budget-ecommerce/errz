package errorz

import (
	"encoding/json"
	"fmt"
	"os"
)

const defaultSchemaPath = "schema/error_schema.json"

// GroupErrorsByDomain reads an error definition JSON file, validates it against
// the default JSON schema, parses its content into structured data, and groups
// the errors by their domain.
func GroupErrorsByDomain(jsonPath string) (map[string]map[string]ErrorDefinition, error) {
	// Step 1: Validate
	if err := ValidateJSON(defaultSchemaPath, jsonPath); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Step 2: Read
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var errorMap map[string]ErrorDefinition
	if err := json.Unmarshal(data, &errorMap); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	if len(errorMap) == 0 {
		return nil, fmt.Errorf("no errors found in input")
	}

	// Step 3: Group by domain
	domainGroups := make(map[string]map[string]ErrorDefinition)
	for code, def := range errorMap {
		if domainGroups[def.Domain] == nil {
			domainGroups[def.Domain] = make(map[string]ErrorDefinition)
		}
		domainGroups[def.Domain][code] = def
	}

	return domainGroups, nil
}
