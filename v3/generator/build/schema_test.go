package build

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

func schemasSpec(schemas string) string {
	return `openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths: {}
components:
  schemas:
` + schemas
}

func TestSchemaNullableTypeArray(t *testing.T) {
	for _, types := range []string{`[integer, "null"]`, `["null", integer]`} {
		t.Run(types, func(t *testing.T) {
			b := newTestBuilder(t, schemasSpec(fmt.Sprintf(`    test:
      type: object
      properties:
        value:
          type: %s
`, types)))
			f, err := b.Schemas("v3")
			require.NoError(t, err)
			st := decl(t, f, "Test").(ir.Struct)
			require.Equal(t, ir.Pointer{Elem: ir.Named("int")}, st.Fields[0].Type)
		})
	}
}

func TestSchemaRejectsUnsupportedNullableTypes(t *testing.T) {
	for _, types := range []string{`[integer, string, "null"]`, `["null"]`} {
		t.Run(types, func(t *testing.T) {
			b := newTestBuilder(t, schemasSpec(fmt.Sprintf(`    test:
      type: object
      properties:
        value:
          type: %s
`, types)))
			_, err := b.Schemas("v3")
			require.EqualError(t, err, `schema test: property "value": nullable type must contain exactly one non-null type`)
		})
	}
}

func TestSchemaRejectsUnknownFormat(t *testing.T) {
	b := newTestBuilder(t, schemasSpec(`    test:
      type: string
      format: email
`))
	_, err := b.Schemas("v3")
	require.EqualError(t, err, `schema test: Test: schema format "email" not implemented`)
}

func TestSchemaObject(t *testing.T) {
	b := newTestBuilder(t, schemasSpec(`    zone:
      type: object
      required: [name]
      properties:
        name:
          type: string
          minLength: 1
        enabled:
          type: boolean
        size:
          type: integer
          format: int64
          minimum: 10
        state:
          type: string
          enum: [running, stopped]
        labels:
          type: object
          additionalProperties:
            type: string
        spec:
          type: object
          properties:
            a:
              type: string
        instance:
          $ref: '#/components/schemas/instance'
    instance:
      type: object
      properties:
        id:
          type: string
          format: uuid
`))
	f, err := b.Schemas("v3")
	require.NoError(t, err)

	// Inline sub schemas are declared before their parent.
	require.Equal(t, []string{"Instance", "ZoneSpec", "ZoneState", "Zone"}, declNames(f))
	require.Equal(t, ir.Struct{
		Name: "Zone",
		Fields: []ir.Field{
			{Name: "Enabled", Type: ir.Pointer{Elem: ir.Named("bool")}, JSON: "enabled", OmitEmpty: true},
			{Name: "Instance", Type: ir.Pointer{Elem: ir.Named("Instance")}, JSON: "instance", OmitEmpty: true},
			{Name: "Labels", Type: ir.Map{Elem: ir.Named("string")}, JSON: "labels", OmitEmpty: true},
			{Name: "Name", Type: ir.Named("string"), JSON: "name", Validate: []string{"required", "gte=1"}},
			{Name: "Size", Type: ir.Named("int64"), JSON: "size", OmitEmpty: true, Validate: []string{"omitempty", "gte=10"}},
			{Name: "Spec", Type: ir.Pointer{Elem: ir.Named("ZoneSpec")}, JSON: "spec", OmitEmpty: true},
			{Name: "State", Type: ir.Named("ZoneState"), JSON: "state", OmitEmpty: true},
		},
	}, decl(t, f, "Zone"))
	require.Equal(t, ir.Enum{
		Name: "ZoneState",
		Base: ir.Named("string"),
		Values: []ir.EnumValue{
			{Name: "ZoneStateRunning", Literal: `"running"`},
			{Name: "ZoneStateStopped", Literal: `"stopped"`},
		},
	}, decl(t, f, "ZoneState"))
}
