// cmd/gen_errors/gen.go

package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/unlimited-budget-ecommerce/errz"
)

const (
	relativeSchemaPath      = "schema/error_schema.json"
	relativeDefinitionsPath = "definitions"
	outputFile              = "errz_gen.go"
	outputDir               = "docs"
)

func main() {
	rootDir, err := projectRoot()
	if err != nil {
		log.Fatalf("cannot determine project root: %v", err)
	}

	gen := errz.Generator{
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
