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
