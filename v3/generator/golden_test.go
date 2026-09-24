package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestGolden regenerates the SDK from source.yaml and checks the output
// matches the committed files: a generator change must not alter them
// unless the committed files are regenerated in the same change.
func TestGolden(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, run("./source.yaml", dir, "v3"))

	for _, name := range []string{"schemas.go", "client.go", "operations.go"} {
		t.Run(name, func(t *testing.T) {
			want, err := os.ReadFile(filepath.Join("..", name))
			require.NoError(t, err)
			got, err := os.ReadFile(filepath.Join(dir, name))
			require.NoError(t, err)
			// Not require.Equal: a diff of these files is too large to be useful.
			if string(want) != string(got) {
				t.Errorf("generated %s differs from ../%s, run `go generate` and inspect the diff", name, name)
			}
		})
	}
}
