//go:generate go run ./cmd/gen_errors/gen.go
package errorz

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadErrorDefinitions_Valid(t *testing.T) {
	defs, err := LoadErrorDefinitions("testdata/valid")
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(defs), 1)
	assert.Contains(t, defs, "CM0000")
}

func TestLoadErrorDefinitions_DirNotFound(t *testing.T) {
	_, err := LoadErrorDefinitions("testdata/does_not_exist")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
}

func TestLoadErrorDefinitions_EmptyFile(t *testing.T) {
	_, err := LoadErrorDefinitions("testdata/empty_file")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal error at")
}

func TestLoadErrorDefinitions_InvalidJSON(t *testing.T) {
	_, err := LoadErrorDefinitions("testdata/invalid_json")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unmarshal error")
}

func TestLoadErrorDefinitions_SkipNonJSONFiles(t *testing.T) {
	defs, err := LoadErrorDefinitions("testdata/mixed_with_nonjson")
	assert.NoError(t, err)
	assert.Contains(t, defs, "CM0002")
	assert.NotContains(t, defs, "FAKECODE")
}

func TestLoadErrorDefinitions_UnreadableFile(t *testing.T) {
	// Create a file and lock the permissions.
	dir := t.TempDir()
	file := filepath.Join(dir, "bad.json")
	err := os.WriteFile(file, []byte(`{}`), 0000)
	assert.NoError(t, err)

	defer os.Remove(file)

	_, err = LoadErrorDefinitions(dir)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "read error")
}

func TestLoadErrorDefinitions_DuplicateKey(t *testing.T) {
	_, err := LoadErrorDefinitions("testdata/dup_key")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate error code detected")
}
