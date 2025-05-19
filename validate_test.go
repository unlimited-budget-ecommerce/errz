package errorz_test

import (
	"strings"
	"testing"

	"github.com/unlimited-budget-ecommerce/errorz"
	testutils "github.com/unlimited-budget-ecommerce/errorz/utils"
)

func TestValidateJSON_Valid(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	valid := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.ValidateJSON(schema, valid)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateJSON_InvalidJSONFile(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	invalid := testutils.WriteTempFile(t, "invalid.json", testutils.InvalidJSON)

	err := errorz.ValidateJSON(schema, invalid)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected validation error, got: %v", err)
	}
}

func TestValidateJSON_InvalidSchemaPath(t *testing.T) {
	err := errorz.ValidateJSON("/non/existing/schema.json", "/dev/null")
	if err == nil || !strings.Contains(err.Error(), "failed to run validation") {
		t.Errorf("expected schema path error, got: %v", err)
	}
}

func TestValidateJSON_InvalidJSONPath(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	err := errorz.ValidateJSON(schema, "/non/existing/file.json")
	if err == nil || !strings.Contains(err.Error(), "failed to run validation") {
		t.Errorf("expected json path error, got: %v", err)
	}
}

func TestValidateJSON_InvalidJSONSyntax(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	badJSON := testutils.WriteTempFile(t, "bad.json", testutils.MalformedJSON)

	err := errorz.ValidateJSON(schema, badJSON)
	if err == nil || !strings.Contains(err.Error(), "failed to run validation") {
		t.Errorf("expected JSON syntax error, got: %v", err)
	}
}

func TestValidateJSON_InvalidSchemaSyntax(t *testing.T) {
	badSchema := testutils.WriteTempFile(t, "bad-schema.json", `{`)
	data := testutils.WriteTempFile(t, "valid.json", testutils.ValidJSON)

	err := errorz.ValidateJSON(badSchema, data)
	if err == nil || !strings.Contains(err.Error(), "failed to run validation") {
		t.Errorf("expected schema syntax error, got: %v", err)
	}
}

func TestValidateJSON_AdditionalProperty(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	extraJSON := testutils.WriteTempFile(t, "extra.json", `{
	  "PM0001": {
		"domain": "payment", "code": "PM0001", "msg": "OK",
		"cause": "none", "http_status": 200, "category": "business",
		"severity": "low", "is_retryable": false, "extra": "not allowed"
	  }
	}`)

	err := errorz.ValidateJSON(schema, extraJSON)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected additionalProperty error, got: %v", err)
	}
}

func TestValidateJSON_MissingRequiredField(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	missingJSON := testutils.WriteTempFile(t, "missing.json", `{
	  "PM0001": {
		"code": "PM0001", "msg": "OK",
		"cause": "none", "http_status": 200,
		"category": "business", "severity": "low", "is_retryable": false
	  }
	}`)

	err := errorz.ValidateJSON(schema, missingJSON)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected required field error, got: %v", err)
	}
}

func TestValidateJSON_InvalidFieldType(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	badTypeJSON := testutils.WriteTempFile(t, "badtype.json", `{
	  "PM0001": {
		"domain": "payment", "code": "PM0001", "msg": "OK",
		"cause": "none", "http_status": "not-number",
		"category": "business", "severity": "low", "is_retryable": false
	  }
	}`)

	err := errorz.ValidateJSON(schema, badTypeJSON)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected field type error, got: %v", err)
	}
}

func TestValidateJSON_InvalidFieldEnum(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	enumFailJSON := testutils.WriteTempFile(t, "enum.json", `{
	  "PM0001": {
		"domain": "payment", "code": "PM0001", "msg": "OK",
		"cause": "none", "http_status": 200,
		"category": "unknown", "severity": "low", "is_retryable": false
	  }
	}`)

	err := errorz.ValidateJSON(schema, enumFailJSON)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected enum validation error, got: %v", err)
	}
}

func TestValidateJSON_InvalidPattern(t *testing.T) {
	schema := testutils.WriteTempFile(t, "schema.json", testutils.SchemaJSON)
	badCodeJSON := testutils.WriteTempFile(t, "pattern.json", `{
	  "PMXXXX": {
		"domain": "payment", "code": "PMXXXX", "msg": "OK",
		"cause": "none", "http_status": 200,
		"category": "business", "severity": "low", "is_retryable": false
	  }
	}`)

	err := errorz.ValidateJSON(schema, badCodeJSON)
	if err == nil || !strings.Contains(err.Error(), "JSON validation failed:") {
		t.Errorf("expected pattern error, got: %v", err)
	}
}
