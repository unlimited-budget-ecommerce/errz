//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import "golang.org/x/sync/errgroup"

type ErrorDefinition struct {
	Domain      string `json:"domain"`
	Code        string `json:"code"`
	Msg         string `json:"msg"`
	Cause       string `json:"cause"`
	Severity    string `json:"severity"`
	IsRetryable bool   `json:"is_retryable"`
}

type Generator struct {
	SchemaPath     string
	DefinitionsDir string
	OutputPath     string
	OutputDocDir   string
}

func (g *Generator) Run() error {
	var errors map[string]ErrorDefinition

	var eg errgroup.Group

	eg.Go(func() error {
		return ValidateAllJSONFiles(g.SchemaPath, g.DefinitionsDir)
	})

	eg.Go(func() error {
		var err error
		errors, err = LoadErrorDefinitions(g.DefinitionsDir)
		return err
	})

	if err := eg.Wait(); err != nil {
		return err
	}

	// Generate code content
	return Generate(g.OutputPath, g.OutputDocDir, errors)
}
