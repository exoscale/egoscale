package operations

import (
	"testing"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/stretchr/testify/require"
)

const testSpec = `openapi: 3.1.0
info:
  title: test
  version: 1.0.0
paths:
  /foo:
    get:
      operationId: get-foo
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  name:
                    type: string
        "400":
          description: Bad Request
          content:
            application/json:
              schema:
                type: object
                properties:
                  reason:
                    type: string
        "429":
          description: Rate limited
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/rate-limited'
        default:
          description: Unexpected
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/rate-limited'
    post:
      operationId: create-foo
      responses:
        "201":
          description: Created
          content:
            application/json:
              schema:
                type: object
                properties:
                  id:
                    type: string
    delete:
      operationId: delete-foo
      responses:
        "204":
          description: No Content
components:
  schemas:
    rate-limited:
      type: object
      properties:
        retry_after:
          type: number
`

func testOperation(t *testing.T, method string) *v3.Operation {
	t.Helper()

	doc, err := libopenapi.NewDocument([]byte(testSpec))
	require.NoError(t, err)
	model, errs := doc.BuildV3Model()
	require.Empty(t, errs)

	item, ok := model.Model.Paths.PathItems.Get("/foo")
	require.True(t, ok)
	op, ok := item.GetOperations().Get(method)
	require.True(t, ok)

	return op
}

func TestRenderResponseSchemaNamesNon200WithCode(t *testing.T) {
	got, err := renderResponseSchema("GetFoo", testOperation(t, "get"))
	require.NoError(t, err)
	require.Contains(t, string(got), "type GetFooResponse struct")
	require.Contains(t, string(got), "type GetFooResponse400 struct")
	require.NotContains(t, string(got), "GetFooResponse429")
	require.NotContains(t, string(got), "GetFooResponsedefault")
}

func TestGetValuesReturn(t *testing.T) {
	require.Equal(t, []string{"*GetFooResponse", "error"}, getValuesReturn(testOperation(t, "get"), "GetFoo"))
	require.Equal(t, []string{"*CreateFooResponse201", "error"}, getValuesReturn(testOperation(t, "post"), "CreateFoo"))
	require.Equal(t, []string{"error"}, getValuesReturn(testOperation(t, "delete"), "DeleteFoo"))
}

func TestGetErrorResponses(t *testing.T) {
	require.Equal(t, []ErrorResponseTmpl{
		{Code: 400, Constructor: "new(GetFooResponse400)"},
		{Code: 429, Constructor: "new(RateLimited)"},
	}, getErrorResponses(testOperation(t, "get"), "GetFoo"))
	require.Nil(t, getErrorResponses(testOperation(t, "delete"), "DeleteFoo"))
}
