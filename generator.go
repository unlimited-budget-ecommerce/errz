package errorz

import "fmt"

// GenerateErrorsFromJSON validates, parses, and generates Go + Markdown error definitions
// from a given JSON file. Outputs are grouped by domain and written to respective folders.
func GenerateErrorsFromJSON(jsonPath string) error {
	domainGroups, err := GroupErrorsByDomain(jsonPath)
	if err != nil {
		return err
	}

	for domain, errors := range domainGroups {
		if err := WriteGoFile("errorz_go", domain, errors); err != nil {
			return fmt.Errorf("failed to write Go file for domain %s: %w", domain, err)
		}

		if err := WriteMarkdown("errorz_doc", domain, errors); err != nil {
			return fmt.Errorf("failed to write markdown for domain %s: %w", domain, err)
		}
	}

	return nil
}
