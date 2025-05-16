package errorz

import (
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSON validates a JSON file using the provided schema file path.
func ValidateJSON(schemaPath, jsonPath string) error {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + jsonPath)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("validation error: %w", err)
	}

	if !result.Valid() {
		msg := "JSON validation failed:\n"
		for _, e := range result.Errors() {
			msg += fmt.Sprintf("- %s\n", e.String())
		}
		return fmt.Errorf("%s", msg)
	}

	return nil
}
