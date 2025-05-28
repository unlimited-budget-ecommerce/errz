//go:generate go run ./cmd/gen_errors/gen.go
package errz

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"golang.org/x/sync/errgroup"
)

// LoadErrorDefinitions loads all JSON files from a directory and returns combined error definitions map.
func LoadErrorDefinitions(dir string) (map[string]ErrorDefinition, error) {
	result := make(map[string]ErrorDefinition)
	var mu sync.Mutex

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var g errgroup.Group
	for _, entry := range entries {
		entry := entry
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(dir, entry.Name())

		g.Go(func() error {
			content, err := os.ReadFile(fullPath)
			if err != nil {
				return fmt.Errorf("read error at %s: %w", fullPath, err)
			}

			var defs map[string]ErrorDefinition
			if err := json.Unmarshal(content, &defs); err != nil {
				return fmt.Errorf("unmarshal error at %s: %w", fullPath, err)
			}

			if len(defs) == 0 {
				return fmt.Errorf("no errors found in %s", fullPath)
			}

			mu.Lock()
			defer mu.Unlock()
			for k, v := range defs {
				if _, exists := result[k]; exists {
					return fmt.Errorf("duplicate error code detected: %s in %s", k, fullPath)
				}
				result[k] = v
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return result, nil
}
