// cmd/gen_errors/gen.go

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
