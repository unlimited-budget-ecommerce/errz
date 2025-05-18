package errorz

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

// LoadAndValidateJSON loads a JSON file, validates it against the provided schema,
// and unmarshals it into a map of ErrorDefinition keyed by error code.
func LoadAndValidateJSON(schemaPath, jsonPath string) (map[string]ErrorDefinition, error) {
	// Validate JSON file against schema
	if err := ValidateJSON(schemaPath, jsonPath); err != nil {
		return nil, fmt.Errorf("validation failed for %s: %w", jsonPath, err)
	}

	// Read JSON file content
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %w", jsonPath, err)
	}

	// Unmarshal JSON content into map[string]ErrorDefinition
	var errorsMap map[string]ErrorDefinition
	if err := json.Unmarshal(data, &errorsMap); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %w", jsonPath, err)
	}
	return errorsMap, nil
}

// LoadErrors concurrently loads multiple JSON error definition files,
// validates them, and combines the results into a single map.
// It returns an error if validation fails or if duplicate error codes are detected.
func LoadErrors(schemaPath string, jsonPaths []string) (map[string]ErrorDefinition, error) {
	var (
		wg       sync.WaitGroup // WaitGroup to wait for all goroutines to finish
		errOnce  sync.Once      // Ensures only the first error is recorded
		firstErr error          // Stores the first error encountered
		mu       sync.Mutex     // Mutex to protect concurrent map writes
	)

	allErrors := make(map[string]ErrorDefinition) // Combined error definitions

	wg.Add(len(jsonPaths))

	for _, path := range jsonPaths {
		go func(p string) {
			defer wg.Done()

			// Load and validate the JSON file
			errDefs, err := LoadAndValidateJSON(schemaPath, p)
			if err != nil {
				// Capture the first error encountered and return early from this goroutine
				errOnce.Do(func() { firstErr = err })
				return
			}

			mu.Lock()
			defer mu.Unlock()

			// Merge loaded error definitions into the combined map
			for code, errDef := range errDefs {
				if _, exists := allErrors[code]; exists {
					// Detect duplicate error codes and record error
					errOnce.Do(func() {
						firstErr = fmt.Errorf("duplicate error code %s found in %s", code, p)
					})
					return
				}
				allErrors[code] = errDef
			}
		}(path)
	}

	wg.Wait()

	// Return any error that occured during loading or merging
	if firstErr != nil {
		return nil, firstErr
	}

	return allErrors, nil
}
