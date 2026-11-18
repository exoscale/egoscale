package schemas

import (
	"fmt"
	"testing"

	"github.com/pb33f/libopenapi"
	"github.com/stretchr/testify/require"
)

func TestRenderSchemaNullableTypeArray(t *testing.T) {
	for _, types := range []string{`[integer, "null"]`, `["null", integer]`} {
		t.Run(types, func(t *testing.T) {
			got, err := renderTestSchema(t, types)
			require.NoError(t, err)
			require.Contains(t, got, "Value *int")
		})
	}
}

func TestRenderSchemaRejectsUnsupportedNullableTypes(t *testing.T) {
	for _, types := range []string{`[integer, string, "null"]`, `["null"]`} {
		t.Run(types, func(t *testing.T) {
			_, err := renderTestSchema(t, types)
			require.EqualError(t, err, `property "value": nullable type must contain exactly one non-null type`)
		})
	}
}

func TestRenderSchemaFormats(t *testing.T) {
	tests := map[string]string{
		"date":      "string",
		"date-time": "time.Time",
		"uuid":      "UUID",
	}
	for format, expected := range tests {
		t.Run(format, func(t *testing.T) {
			got, err := renderTestSchema(t, "string\n          format: "+format)
			require.NoError(t, err)
			require.Contains(t, got, "Value "+expected)
		})
	}
}

func TestRenderSchemaDateVariants(t *testing.T) {
	tests := map[string]struct {
		schema   string
		expected string
	}{
		"required": {
			schema: `type: object
      required: [value]
      properties:
        value:
          type: string
          format: date`,
			expected: "Value string `json:\"value\" validate:\"required\"`",
		},
		"nullable": {
			schema: `type: object
      properties:
        value:
          type: [string, "null"]
          format: date`,
			expected: "Value *string",
		},
		"array": {
			schema: `type: object
      properties:
        value:
          type: array
          items:
            type: string
            format: date`,
			expected: "Value []string",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			doc, err := libopenapi.NewDocument([]byte(`openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths: {}
components:
  schemas:
    test:
      ` + tc.schema + "\n"))
			require.NoError(t, err)
			model, errs := doc.BuildV3Model()
			require.Empty(t, errs)
			got, err := RenderSchema("test", model.Model.Components.Schemas.GetOrZero("test"))
			require.NoError(t, err)
			require.Contains(t, string(got), tc.expected)
		})
	}
}

func renderTestSchema(t *testing.T, types string) (string, error) {
	t.Helper()

	doc, err := libopenapi.NewDocument([]byte(fmt.Sprintf(`openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths: {}
components:
  schemas:
    test:
      type: object
      properties:
        value:
          type: %s
`, types)))
	require.NoError(t, err)

	model, errs := doc.BuildV3Model()
	require.Empty(t, errs)

	result, err := RenderSchema("test", model.Model.Components.Schemas.GetOrZero("test"))
	return string(result), err
}
