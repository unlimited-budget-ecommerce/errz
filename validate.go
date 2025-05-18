package errorz

import (
	"fmt"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSON validates a JSON file against a JSON Schema located at schemaPath.
// If validation fails, it returns a detailed error message with all issues found.
func ValidateJSON(schemaPath, jsonPath string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + jsonPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("failed to run validation: %w", err)
	}

	if !result.Valid() {
		var builder strings.Builder
		builder.WriteString("JSON validation failed:\n")
		for _, e := range result.Errors() {
			builder.WriteString(fmt.Sprintf("- %s\n", e.String()))
		}
		return fmt.Errorf("%s", builder.String())
	}

	return nil
}
