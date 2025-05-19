package errorz_test

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/unlimited-budget-ecommerce/errorz"
	testutils "github.com/unlimited-budget-ecommerce/errorz/utils"
)

var ErrGoGenerationFails = errors.New("go content fail")
var ErrMdGenerationFails = errors.New("md fail")

// mock versions of dependencies
func mockGenerateGoContent(packageName string, errors map[string]errorz.ErrorDefinition) (string, error) {
	return "// mock Go content", nil
}

func mockGenerateMarkdownContent(domain string, errors map[string]errorz.ErrorDefinition) (string, error) {
	return "# mock Markdown content", nil
}

func mockWriteToFile(path, content string) error {
	return nil
}

func mockMkdirAll(path string, perm os.FileMode) error {
	return nil
}

func setupMocks() {
	errorz.GenerateGoContentFunc = mockGenerateGoContent
	errorz.GenerateMarkdownContentFunc = mockGenerateMarkdownContent
	errorz.WriteToFileFunc = mockWriteToFile
	errorz.MkdirAll = mockMkdirAll
}

func TestGenerateErrorsFromJSON_ValidSingle(t *testing.T) {
	setupMocks()
	path := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.NoError(t, err)
}

func TestGenerateErrorsFromJSON_ValidMany(t *testing.T) {
	setupMocks()
	path := testutils.WriteTempFile(t, "many.json", testutils.ManyValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.NoError(t, err)
}

func TestGenerateErrorsFromJSON_InvalidJSONSchema(t *testing.T) {
	setupMocks()
	path := testutils.WriteTempFile(t, "invalid.json", testutils.InvalidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "JSON validation failed")
}

func TestGenerateErrorsFromJSON_MalformedJSON(t *testing.T) {
	setupMocks()
	path := testutils.WriteTempFile(t, "malformed.json", testutils.MalformedJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestGenerateErrorsFromJSON_GoGenerationFails(t *testing.T) {
	// override GenerateGoContent to return error
	errorz.GenerateGoContentFunc = func(packageName string, errors map[string]errorz.ErrorDefinition) (string, error) {
		return "", ErrGoGenerationFails
	}
	errorz.GenerateMarkdownContentFunc = mockGenerateMarkdownContent
	errorz.WriteToFileFunc = mockWriteToFile
	errorz.MkdirAll = mockMkdirAll

	path := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "go content fail")
}

func TestGenerateErrorsFromJSON_MarkdownGenerationFails(t *testing.T) {
	errorz.GenerateGoContentFunc = mockGenerateGoContent
	errorz.GenerateMarkdownContentFunc = func(domain string, errors map[string]errorz.ErrorDefinition) (string, error) {
		return "", ErrMdGenerationFails
	}
	errorz.WriteToFileFunc = mockWriteToFile
	errorz.MkdirAll = mockMkdirAll

	path := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "md fail")
}

func TestGenerateErrorsFromJSON_WriteToFileFails(t *testing.T) {
	errorz.GenerateGoContentFunc = mockGenerateGoContent
	errorz.GenerateMarkdownContentFunc = mockGenerateMarkdownContent
	errorz.WriteToFileFunc = func(path, content string) error {
		return errors.New("write fail")
	}
	errorz.MkdirAll = mockMkdirAll

	path := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "write fail")
}

func TestGenerateErrorsFromJSON_MkdirFails(t *testing.T) {
	errorz.GenerateGoContentFunc = mockGenerateGoContent
	errorz.GenerateMarkdownContentFunc = mockGenerateMarkdownContent
	errorz.WriteToFileFunc = mockWriteToFile
	errorz.MkdirAll = func(path string, perm os.FileMode) error {
		return errors.New("mkdir fail")
	}

	path := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.GenerateErrorsFromJSON(path)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "mkdir fail")
}
