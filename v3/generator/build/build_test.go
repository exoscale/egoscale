package build

import (
	"slices"
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
		if slices.Contains(typeNames(d), name) {
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
		names = append(names, typeNames(d)...)
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

func TestEmitRejectsDuplicateAliases(t *testing.T) {
	model, err := Load([]byte(schemasSpec(`    foo:
      type: string
`)))
	require.NoError(t, err)
	b := New(model, config.Config{Aliases: []config.Alias{{Name: "Foo", Target: "string"}}})

	_, err = b.Schemas("v3")
	require.EqualError(t, err, "type Foo from config is already declared from #/components/schemas/foo")
}

func TestCheckRejectsUnusedRefOverrides(t *testing.T) {
	model, err := Load([]byte(schemasSpec(`    foo:
      type: object
      properties:
        bar:
          $ref: '#/components/schemas/bar-ref'
    bar-ref:
      type: string
`)))
	require.NoError(t, err)
	b := New(model, config.Config{RefOverrides: map[string]map[string]string{
		"Foo": {"#/components/schemas/bar-ref": "Bar"},
		"Baz": {"#/components/schemas/bar-ref": "Bar"},
	}})

	f, err := b.Schemas("v3")
	require.NoError(t, err)
	require.Equal(t, ir.Named("Bar"), decl(t, f, "Foo").(ir.Struct).Fields[0].Type)
	require.EqualError(t, b.Check(), "ref override Baz #/components/schemas/bar-ref: no such reference in the spec")
}
