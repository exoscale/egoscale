package build

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

const operationsSpec = `openapi: 3.1.0
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
  /foo/{id}:
    delete:
      operationId: delete-foo
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
            format: uuid
      responses:
        "204":
          description: No Content
  /zones:
    get:
      operationId: list-zones
      parameters:
        - name: state
          in: query
          schema:
            type: string
            enum: [up, down]
        - name: limit
          in: query
          schema:
            type: integer
      responses:
        "200":
          description: OK
          content:
            application/json:
              schema:
                type: object
                properties:
                  zones:
                    type: array
                    items:
                      $ref: '#/components/schemas/zone'
components:
  schemas:
    rate-limited:
      type: object
      properties:
        retry_after:
          type: number
    zone:
      type: object
      properties:
        name:
          type: string
`

func buildOperations(t *testing.T) ir.File {
	t.Helper()

	b := newTestBuilder(t, operationsSpec)
	b.cfg.SkipAuthOperations = []string{"list-zones"}
	f, err := b.Operations("v3")
	require.NoError(t, err)

	return f
}

func TestOperationsResponseSchemas(t *testing.T) {
	names := declNames(buildOperations(t))
	require.Contains(t, names, "GetFooResponse")
	require.Contains(t, names, "GetFooResponse400")
	require.Contains(t, names, "CreateFooResponse201")
	require.NotContains(t, names, "GetFooResponse429")
	require.NotContains(t, names, "GetFooResponsedefault")
}

func TestOperationReturnAndErrors(t *testing.T) {
	f := buildOperations(t)

	get := decl(t, f, "GetFoo").(ir.Operation)
	require.Equal(t, ir.Pointer{Elem: ir.Named("GetFooResponse")}, get.Return)
	require.Equal(t, []ir.ErrorResponse{
		{Code: 400, Type: ir.Named("GetFooResponse400")},
		{Code: 429, Type: ir.Named("RateLimited")},
	}, get.ErrorResponses)

	create := decl(t, f, "CreateFoo").(ir.Operation)
	require.Equal(t, ir.Pointer{Elem: ir.Named("CreateFooResponse201")}, create.Return)

	del := decl(t, f, "DeleteFoo").(ir.Operation)
	require.Nil(t, del.Return)
	require.Nil(t, del.ErrorResponses)
	require.Equal(t, []ir.Param{{Name: "id", Type: ir.Named("UUID")}}, del.Params)
	require.Equal(t, "/foo/%v", del.Path)
	require.Equal(t, []string{"id"}, del.PathArgs)
	require.Equal(t, "DELETE", del.Method)
}

func TestOperationQueryOptsAndFindable(t *testing.T) {
	f := buildOperations(t)

	list := decl(t, f, "ListZones").(ir.Operation)
	require.True(t, list.SkipAuth)
	require.Equal(t, "ListZonesOpt", list.QueryOpts)
	require.Equal(t, "/zones", list.Path)
	require.Nil(t, list.PathArgs)

	require.Equal(t, ir.QueryOpts{
		Type: "ListZonesOpt",
		Funcs: []ir.QueryOpt{
			{Func: "ListZonesWithState", Arg: "state", ArgType: ir.Named("ListZonesState"), Param: "state"},
			{Func: "ListZonesWithLimit", Arg: "limit", ArgType: ir.Named("int"), Param: "limit"},
		},
	}, decl(t, f, "ListZonesOpt"))

	var findable *ir.Findable
	for _, d := range f.Decls {
		if fd, ok := d.(ir.Findable); ok {
			findable = &fd
		}
	}
	require.Equal(t, &ir.Findable{
		ListType:  "ListZonesResponse",
		ListField: "Zones",
		ItemType:  "Zone",
		Param:     "name",
		Fields:    []string{"Name"},
	}, findable)
}
