package errorz_test

import (
	"strings"
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
	utils_test "github.com/unlimited-budget-ecommerce/errorz/utils"
)

func TestLoadAndValidateJSON_Success(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	jsonPath := utils_test.WriteTempFile(t, "valid.json", utils_test.ValidJSON)

	result, err := errorz.LoadAndValidateJSON(schemaPath, jsonPath)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if _, ok := result["PM0001"]; !ok {
		t.Errorf("expected PM0001 in result")
	}
}

func TestLoadAndValidateJSON_InvalidSchema(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	jsonPath := utils_test.WriteTempFile(t, "invalid.json", utils_test.InvalidJSON)

	_, err := errorz.LoadAndValidateJSON(schemaPath, jsonPath)
	if err == nil || !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected schema validation error, got: %v", err)
	}
}

func TestLoadAndValidateJSON_MalformedJSON(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	jsonPath := utils_test.WriteTempFile(t, "malformed.json", utils_test.MalformedJSON)

	_, err := errorz.LoadAndValidateJSON(schemaPath, jsonPath)
	if err == nil || !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected parse error, got: %v", err)
	}
}

func TestLoadAndValidateJSON_FileNotFound(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)

	_, err := errorz.LoadAndValidateJSON(schemaPath, "not_exists.json")
	if err == nil || !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected file read error, got: %v", err)
	}
}

func TestLoadErrors_SuccessMultipleFiles(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	jsonPath1 := utils_test.WriteTempFile(t, "file1.json", utils_test.ValidJSON)
	jsonPath2 := utils_test.WriteTempFile(t, "file2.json", strings.Replace(utils_test.ValidJSON, "PM0001", "PM0002", 1))

	errors, err := errorz.LoadErrors(schemaPath, []string{jsonPath1, jsonPath2})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if len(errors) != 2 {
		t.Errorf("expected 2 errors, got: %d", len(errors))
	}
}

func TestLoadErrors_DuplicateCode(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	jsonPath1 := utils_test.WriteTempFile(t, "file1.json", utils_test.ValidJSON)
	jsonPath2 := utils_test.WriteTempFile(t, "file2.json", utils_test.ValidJSON) // same code: PM0001

	_, err := errorz.LoadErrors(schemaPath, []string{jsonPath1, jsonPath2})
	if err == nil || !strings.Contains(err.Error(), "duplicate error code PM0001") {
		t.Errorf("expected duplicate error code, got: %v", err)
	}
}

func TestLoadErrors_WithInvalidFile(t *testing.T) {
	schemaPath := utils_test.WriteTempFile(t, "schema.json", utils_test.SchemaJSON)
	invalidPath := utils_test.WriteTempFile(t, "bad.json", utils_test.InvalidJSON)

	_, err := errorz.LoadErrors(schemaPath, []string{invalidPath})
	if err == nil || !strings.Contains(err.Error(), "validation failed") {
		t.Errorf("expected validation failure, got: %v", err)
	}
}
