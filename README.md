# errorz - Centralize error library

`errorz` is a centralized error code management and generation tool for Go projects. It reads structured error definitions from JSON files, validates them against a JSON Schema, and generates Go source code and Markdown documentation.

## Features

- JSON Schema validation
- Code generation for:
  - Go: structured error variables and `ErrorsMap` for fast lookup
  - Markdown: human-readable documentation grouped by domain

## Installation

```bash
go get github.com/unlimited-budget-ecommerce/errorz
```

## Configuration

This project uses JSON files to define error definitions, validated against a JSON Schema to ensure correct format.

- Schema JSON file: `/schema/error_schema.json`
- The JSON error definitions must be an object with error codes as keys.(Error codes must follow the pattern: 2 uppercase letters followed by 4 digits, e.g. `PM0001`.)
- Each error definition must include the following fields:

| Field          |   Type    |   Required    | Description                         |
| :------------- | :-------: | :-----------: | :---------------------------------- |
| `domain`       |  string   |      ✅       | Logical domain (e.g. `"auth"`)      |
| `code`         |  string   |      ✅       | Unique code, like `"PM0001"`        |
| `msg`          |  string   |      ✅       | User-friendly message               |
| `cause`        |  string   |      ✅       | Root cause of the error             |
| `http_status`  |  integer  |      ✅       | HTTP status code (100–599)          |
| `category`     |  string   |      ✅       | `validation`, `business`, etc.      |
| `severity`     |  string   |      ✅       | `low`, `medium`, `high`, `critical` |
| `solution`     |  string   | ❌ (optional) | Suggested fix (if available)        |
| `is_retryable` |  boolean  |      ✅       | Whether it's safe to retry          |
| `tags`         | \[]string | ❌ (optional) | Optional grouping keywords          |

Example error definition JSON:

```json
{
  "PM0001": {
    "domain": "payment",
    "code": "PM0001",
    "msg": "insufficient balance",
    "cause": "user has not enough balance",
    "http_status": 402,
    "category": "business",
    "severity": "medium",
    "is_retryable": false,
    "solution": "ask user to top-up or choose another method",
    "tags": ["payment", "balance"]
  }
}
```

## Generator

### Generator Pattern

Use `Generator()` for unified generation:

```go
package main

import (
  "fmt"
  "log"
  "path/filepath"

  "github.com/unlimited-budget-ecommerce/errorz"
)

const (
  relativeSchemaPath      = "schema/error_schema.json"
  relativeDefinitionsPath = "error_definitions"
  outputFile              = "errors_gen.go"
  outputDir               = "docs"
)

func main() {
  rootDir, err := errorz.ProjectRoot()
    if err != nil {
      log.Fatalf("cannot determine project root: %v", err)
    }

  gen := errorz.Generator{
    SchemaPath:     filepath.Join(rootDir, relativeSchemaPath),
    DefinitionsDir: filepath.Join(rootDir, relativeDefinitionsPath),
    OutputPath:     filepath.Join(rootDir, outputFile),
    OutputDocDir:   filepath.Join(rootDir, outputDir),
  }

  if err := gen.Run(); err != nil {
    log.Fatalf("generate failed: %v", err)
  }

  fmt.Println("Generated", outputFile)
}
```

Or step-by-step (if preferred):

```go
errors := errorz.LoadErrorDefinitions("error_definitions")
errorz.ValidateAllJSONFiles("schema/error_schema.json", "error_definitions")
errorz.Generator("errors_gen.go", "docs", errors)
```

## Usage and Output

### Error code catalog

You can get a quick overview of all error codes and their meaning in `errorz_code_catalog.md`

### Go generation contains (Already Generated – Ready to Use)

- Error struct
- Global variables (e.g., PM0001)
- ErrorsMap map for fast string-based lookup:

```go
err := errorz.ErrorsMap["PM0001"] // preferred for performance
```

> **Note:**
>
> ✅ No need to generate anything yourself. This package already includes the generated Go code.  
> 👉 Just import and use the variables or ErrorsMap directly!
>
> - ErrorsMap["CODE"] is recommended for dynamic lookups.
> - Use errorz.PM0001 for static compile-time usage.

### Markdown generation contains

- Generated in `/docs` (or configured output directory), grouped by domain and including all metadata.

> **Note:**
>
> 👉 You can view human-readable error definitions in the `docs` directory.

## Example Error Struct

```go
type Error struct {
  Code        string
  Msg         string
  Cause       string
  HTTPStatus  int
  Category    string
  Severity    string
  IsRetryable bool
  Solution    string
  Tags        []string
}
```

## Validations

- JSON is validated using **[xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema.git)**

## Tips

- Use `ErrorsMap["CODE"]` when lookup is based on string (e.g., from logs or API).
- Keep your domain files (e.g., auth.json, payment.json) separate for clarity.
