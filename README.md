# errorz - Centralize error library

`errorz` is a centralized error code management and generation tool for Go projects. It reads structured error definitions from JSON files, validates them against a JSON Schema, and generates Go source code and Markdown documentation.

## Features

- JSON Schema validation
- Concurrent loading of multiple JSON error files
- Code generation for:
  - Go: structured error variables and slice
  - Markdown: human-readable documentation grouped by domain
- Duplicate error code detection
- Caching of Markdown headers for performance

## Installation

```bash
go get github.com/unlimited-budget-ecommerce/errorz
```

## Configuration

This project uses JSON files to define error definitions, validated against a JSON Schema to ensure correct format.

- Schema JSON file: `/schema/error_schema.json`
- The JSON error definitions must be an object with error codes as keys (error codes must follow the pattern: 2 uppercase letters followed by 4 digits, e.g. `AB1234`).
- Each error definition must include the following fields:
  - `domain` (string) — error domain/category (e.g., "payment", "user-auth")
  - `code` (string) — error code matching the defined pattern
  - `msg` (string) — error message to display
  - `cause` (string) — cause of the error
  - `http_status` (integer) — related HTTP status code (100-599)
  - `category` (string) — error category (`validation`, `timeout`, `business`, `external`, `internal`)
  - `severity` (string) — severity level (`low`, `medium`, `high`, `critical`)
  - `solution` (string, optional) — suggested solution to fix the error
  - `is_retryable` (boolean) — indicates if the error is retryable
  - `tags` (array of strings, optional) — tags for additional grouping

Example error definition JSON:

```json
{
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
}
```

> **Note:**  
> You should create your own JSON file containing error definitions, for example at `errors/example.json`, which will be used as input to generate Go code and Markdown documentation.

## Usage

- Validate and load error definitions
- Generate Go source file
- Generate Markdown documentation

### Usage Example

```bash
package main

import (
    "encoding/json"
    "fmt"
    "os"

    "github.com/unlimited-budget-ecommerce/errorz"
)

func main() {
    // jsonPath should point to your own error definitions JSON file.
    // Replace this with the actual path to your input file.
    jsonPath := "errors/example.json"

    // Generate Go code and Markdown
    err := errorz.GenerateErrorsFromJSON(jsonPath)
    if err != nil {
      fmt.Printf("Could not generate error form JSON: %v", err)
      os.Exit(1)
    }

    fmt.Println("Generation completed successfully.")
}
```

## Output

- Go generation contains:
  - Error struct
  - Error variables like ER0001, ER0002, ...
  - Errors slice
- Markdown generation contains:
  - Markdown table of all user domain errors

## JSON Schema

Located at: `/schema/error_schema.json`

Supports fields such as:

- code, msg, cause, domain, http_status
- category: validation, timeout, business, external, internal
- severity: low, medium, high, critical
- is_retryable: boolean
- Optional: solution, tags

## Example Error Struct

```bash
type Error struct {
    Code       string
    Msg        string
    HTTPStatus int
}
```

## Validations

- JSON is validated using **[xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema.git)**
