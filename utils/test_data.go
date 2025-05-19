// FOR UNIT TEST
package testutils

const SchemaJSON = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "type": "object",
  "patternProperties": {
    "^[A-Z]{2}\\d{4}$": {
      "type": "object",
      "required": [
        "domain", "code", "msg",
        "cause", "http_status", "category",
        "severity", "is_retryable"
      ],
      "properties": {
        "domain": { "type": "string", "minLength": 1 },
        "code": { "type": "string", "pattern": "^[A-Z]{2}\\d{4}$" },
        "msg": { "type": "string", "minLength": 1 },
        "cause": { "type": "string", "minLength": 1 },
        "http_status": { "type": "integer", "minimum": 100, "maximum": 599 },
        "category": { "type": "string", "enum": ["validation", "timeout", "business", "external", "internal"] },
        "severity": { "type": "string", "enum": ["low", "medium", "high", "critical"] },
        "solution": { "type": "string" },
        "is_retryable": { "type": "boolean" },
        "tags": {
          "type": "array",
          "items": { "type": "string" },
          "uniqueItems": true
        }
      },
      "additionalProperties": false
    }
  },
  "additionalProperties": false
}`

const ValidJSON = `{
  "PM0001": {
    "domain": "payment",
    "code": "PM0001",
    "msg": "payment failed",
    "cause": "insufficient balance",
    "http_status": 400,
    "category": "business",
    "severity": "medium",
    "solution": "ask user to top-up balance",
    "is_retryable": false,
    "tags": ["payment", "user"]
  }
}`

const ManyValidJSON = `{
   "PM0001": {
    "domain": "payment",
    "code": "PM0001",
    "msg": "payment failed",
    "cause": "insufficient balance",
    "http_status": 400,
    "category": "business",
    "severity": "medium",
    "solution": "ask user to top-up balance",
    "is_retryable": false,
    "tags": ["payment", "user"]
  },
  "PM0002": {
    "domain": "payment",
    "code": "PM0002",
    "msg": "timeout",
    "cause": "payment provider slow",
    "http_status": 504,
    "category": "timeout",
    "severity": "high",
    "solution": "retry request",
    "is_retryable": true,
    "tags": ["payment", "timeout"]
  },
  "AU0001": {
    "domain": "auth",
    "code": "AU0001",
    "msg": "invalid token",
    "cause": "token expired",
    "http_status": 401,
    "category": "internal",
    "severity": "high",
    "solution": "refresh token",
    "is_retryable": true,
    "tags": ["auth", "security"]
  },
  "OR0001": {
    "domain": "order",
    "code": "OR0001",
    "msg": "order not found",
    "cause": "missing ID",
    "http_status": 404,
    "category": "validation",
    "severity": "low",
    "solution": "check order ID",
    "is_retryable": false,
    "tags": ["order", "lookup"]
  }
}`

const InvalidJSON = `{
  "PM0001": {
    "domain": "payment",
    "code": "INVALID_CODE",
    "msg": "",
    "http_status": 700,
    "category": "invalid_category",
    "severity": "unknown",
    "is_retryable": "not_a_boolean"
  }
}`

const MalformedJSON = `{
  "PM0001": {
    "domain": "payment",
    "code": "PM0001",
    // missing closing brace
`
