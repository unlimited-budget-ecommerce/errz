# errorz - Centralize error library

`errz` is a centralized error code management and generation tool for Go projects. It reads structured error definitions from JSON files, validates them against a JSON Schema, and generates Go source code and Markdown documentation.

## Features

- JSON Schema validation
- Code generation for:
  - Go: structured error variables
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

| Field          |  Type   | Required | Description                         |
| :------------- | :-----: | :------: | :---------------------------------- |
| `domain`       | string  |    ✅    | Logical domain (e.g. `"auth"`)      |
| `code`         | string  |    ✅    | Unique code, like `"PM0001"`        |
| `msg`          | string  |    ✅    | User-friendly message               |
| `cause`        | string  |    ✅    | Root cause of the error             |
| `severity`     | string  |    ✅    | `low`, `medium`, `high`, `critical` |
| `is_retryable` | boolean |    ✅    | Whether it's safe to retry          |

Example error definition JSON:

```json
{
  "PM0001": {
    "domain": "payment",
    "code": "PM0001",
    "msg": "insufficient balance",
    "cause": "user has not enough balance",
    "severity": "medium",
    "is_retryable": false
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
  outputFile              = "output/errz_gen.go"
  outputDir               = "output/docs"
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
}
```

Or step-by-step (if preferred):

```go
errors := errorz.LoadErrorDefinitions("error_definitions")
errorz.ValidateAllJSONFiles("schema/error_schema.json", "error_definitions")
errorz.Generator("output/errors_gen.go", "output/docs", errors)
```

## Usage and Output

### Error code catalog

You can get a quick overview of all error codes and their meaning in `errorz_code_catalog.md`

### Go generation contains (Already Generated – Ready to Use)

- Error struct (implements Go's built-in `error` interface)
- Global variables (e.g., PM0001):
- How to Use the Generated Errors (4 Ways)

  - Use as error directly

  ```go
  return errz.PM0001
  ```

  - Pretty-print with fmt.Println or log

  ```go
  fmt.Println(errz.PM0001)
  // Output:
  // [Domain: payment] [Code: PM0001] Msg: insufficient balance | Cause: user has not enough balance | Severity: medium | Retryable: false
  ```

  - Type assertion for accessing fields

  ```go
  var err error = errz.PM0001
  if e, ok := err.(*errz.Error); ok {
    fmt.Println("Error code:", e.Code)
  }
  ```

  - Conditionally act on metadata

  ```go
  if errz.PM0002.IsRetryable {
    retry()
  } else {
    alertAdmin(errz.PM0002.Msg)
  }

  ```

> **Note:**
>
> ✅ Each generated error variable implements Go's built-in `error` interface, so you can use them directly with `fmt.Println`, `return`, or any function expecting an `error`.  
> ✅ No need to generate anything yourself. This package already includes the generated Go code in `output/errz_gen.go`.  
> 👉 Just import and use the variables directly!

### Markdown generation contains

- Generated in `output/docs` (or configured output directory), grouped by domain and including all metadata.

> **Note:**
>
> 👉 You can view human-readable error definitions in the `docs` directory.

## Example Error Struct

```go
type Error struct {
  Domain      string
  Code        string
  Msg         string
  Cause       string
  Severity    string
  IsRetryable bool
}
```

## Validations

- JSON is validated using **[xeipuuv/gojsonschema](https://github.com/xeipuuv/gojsonschema.git)**

## Tips

- Keep your domain files (e.g., auth.json, payment.json) separate for clarity.
