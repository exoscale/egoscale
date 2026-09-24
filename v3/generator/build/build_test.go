package build

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/ir"
)

func newTestBuilder(t *testing.T, spec string) *Builder {
	t.Helper()

	model, err := Load([]byte(spec))
	require.NoError(t, err)

	return New(model, config.Config{})
}

// decl returns the declaration named name.
func decl(t *testing.T, f ir.File, name string) ir.Decl {
	t.Helper()

	for _, d := range f.Decls {
		if declName(d) == name {
			return d
		}
		if op, ok := d.(ir.Operation); ok && op.Name == name {
			return d
		}
	}
	require.Failf(t, "declaration not found", name)

	return nil
}

func declNames(f ir.File) []string {
	var names []string
	for _, d := range f.Decls {
		if name := declName(d); name != "" {
			names = append(names, name)
		}
	}

	return names
}

func TestEmitRejectsDuplicateTypes(t *testing.T) {
	b := newTestBuilder(t, `openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths: {}
components:
  schemas:
    snapshot:
      type: object
      properties:
        export:
          type: object
          properties:
            url:
              type: string
    snapshot-export:
      type: object
      properties:
        url:
          type: string
`)
	_, err := b.Schemas("v3")
	require.EqualError(t, err, "schema snapshot-export: type SnapshotExport from #/components/schemas/snapshot-export"+
		" is already declared from #/components/schemas/snapshot")
}
