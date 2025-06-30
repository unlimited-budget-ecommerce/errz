package errz

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateJSON_Valid(t *testing.T) {
	schema := "testdata/error_schema.json"
	validFile := "testdata/valid/common.json"

	err := validateJSON(schema, validFile)
	assert.NoError(t, err)
}

func TestValidateJSON_Invalid(t *testing.T) {
	schema := "testdata/error_schema.json"
	invalidFile := "testdata/invalid/invalid_missing_required.json"

	err := validateJSON(schema, invalidFile)
	assert.Error(t, err)
	assert.True(t, strings.Contains(err.Error(), "JSON validation failed"))
}

func TestValidateJSON_FileNotFound(t *testing.T) {
	schema := "testdata/error_schema.json"
	missingFile := "testdata/missing.json"

	err := validateJSON(schema, missingFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file not found")
}

func TestValidateJSON_InvalidSchema(t *testing.T) {
	schema := "testdata/invalid_schema.json" // malformed schema
	validFile := "testdata/valid/common.json"

	err := validateJSON(schema, validFile)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to run validation")
}

func TestValidateAllJSONFiles_AllValid(t *testing.T) {
	schema := "testdata/error_schema.json"
	dir := "testdata/valid"

	err := validateAllJSONFiles(schema, dir)
	assert.NoError(t, err)
}

func TestValidateAllJSONFiles_HasInvalid(t *testing.T) {
	schema := "testdata/error_schema.json"
	dir := "testdata/mixed" // includes both valid and invalid JSONs

	err := validateAllJSONFiles(schema, dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed for")
}

func TestValidateAllJSONFiles_DirNotFound(t *testing.T) {
	schema := "testdata/error_schema.json"
	dir := "testdata/notfound"

	err := validateAllJSONFiles(schema, dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read directory")
}

func TestValidateAllJSONFiles_SkipNonJSON(t *testing.T) {
	schema := "testdata/error_schema.json"
	dir := "testdata/mixed_with_nonjson"

	err := validateAllJSONFiles(schema, dir)
	assert.NoError(t, err)
}

func TestValidateAllJSONFiles_EmptyDirectory(t *testing.T) {
	schema := "testdata/error_schema.json"
	dir := "testdata/empty"

	err := validateAllJSONFiles(schema, dir)
	assert.Error(t, err)
}

func TestValidateJSON_ExtraFields(t *testing.T) {
	schema := "testdata/error_schema.json"
	file := "testdata/invalid/invalid_extra_field.json"

	err := validateJSON(schema, file)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Additional property")
}
