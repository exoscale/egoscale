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
			got, err := renderTestSchema(t, types, "")
			require.NoError(t, err)
			require.Contains(t, got, "Value *int")
		})
	}
}

func TestRenderSchemaRejectsUnsupportedNullableTypes(t *testing.T) {
	for _, types := range []string{`[integer, string, "null"]`, `["null"]`} {
		t.Run(types, func(t *testing.T) {
			_, err := renderTestSchema(t, types, "")
			require.EqualError(t, err, `property "value": nullable type must contain exactly one non-null type`)
		})
	}
}

const nlbSelectorSpec = `openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths: {}
components:
  schemas:
    labels:
      type: object
      additionalProperties:
        type: string
    load-balancer-service-label-selector:
      type: object
      properties:
        labels:
          $ref: '#/components/schemas/labels'
      required:
        - labels
      description: Load Balancer Service label selector
    load-balancer-service:
      type: object
      properties:
        label-selector:
          $ref: '#/components/schemas/load-balancer-service-label-selector'
        instance-pool:
          type: object
          properties:
            id:
              type: string
              format: uuid
`

func TestRenderSchemaOptionalObjectReferenceKeepsOmitempty(t *testing.T) {
	got, err := renderTestSchema(t, nlbSelectorSpec, "load-balancer-service")
	require.NoError(t, err)

	// A non-required reference to an object is a pointer with "omitempty" and no
	// validate tag at all (renderValidation drops a lone "omitempty").
	require.Contains(t, got, "LabelSelector *LoadBalancerServiceLabelSelector `json:\"label-selector,omitempty\"`")

	// instance-pool is not in the required list either: it must lose the
	// "required" validate tag rather than turn into "validate:\"omitempty\"".
	require.Contains(t, got, "InstancePool *LoadBalancerServiceInstancePool `json:\"instance-pool,omitempty\"`")
}

func TestRenderSchemaRequiredMapPropertyHasNoOmitempty(t *testing.T) {
	got, err := renderTestSchema(t, nlbSelectorSpec, "load-balancer-service-label-selector")
	require.NoError(t, err)

	// Required properties never get "omitempty", which is what allows an empty
	// map to be marshalled as `{"labels":{}}` instead of being dropped.
	require.Contains(t, got, "Labels Labels `json:\"labels\" validate:\"required\"`")
	require.NotContains(t, got, `json:"labels,omitempty"`)

	// Referencing the shared labels schema reuses the Labels type instead of
	// minting a per-schema map alias.
	require.NotContains(t, got, "LoadBalancerServiceLabelSelectorLabels")
}

func renderTestSchema(t *testing.T, spec, schemaName string) (string, error) {
	t.Helper()

	if schemaName == "" {
		schemaName = "test"
		spec = fmt.Sprintf(`openapi: 3.1.0
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
`, spec)
	}

	doc, err := libopenapi.NewDocument([]byte(spec))
	require.NoError(t, err)

	model, errs := doc.BuildV3Model()
	require.Empty(t, errs)

	result, err := RenderSchema(schemaName, model.Model.Components.Schemas.GetOrZero(schemaName))
	return string(result), err
}
