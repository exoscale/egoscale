package render

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

func renderFile(t *testing.T, decls ...ir.Decl) string {
	t.Helper()

	src, err := File(ir.File{Package: "v3", Imports: []string{"context", "fmt", "net/http", "net/url", "time"}, Decls: decls})
	require.NoError(t, err)

	return string(src)
}

func TestStruct(t *testing.T) {
	got := renderFile(t, ir.Struct{
		Doc:  "// Zone",
		Name: "Zone",
		Fields: []ir.Field{
			{Doc: "// Zone name", Name: "Name", Type: ir.Named("string"), JSON: "name", Validate: []string{"required"}},
			{Name: "Size", Type: ir.Pointer{Elem: ir.Named("int")}, JSON: "size", OmitEmpty: true, Validate: []string{"omitempty", "gte=1"}},
		},
	})
	require.Contains(t, got, "// Zone\ntype Zone struct {\n"+
		"\t// Zone name\n"+
		"\tName string `json:\"name\" validate:\"required\"`\n"+
		"\tSize *int   `json:\"size,omitempty\" validate:\"omitempty,gte=1\"`\n"+
		"}\n")
}

func TestQueryOpts(t *testing.T) {
	got := renderFile(t, ir.QueryOpts{
		Type: "ListOpt",
		Funcs: []ir.QueryOpt{
			{Func: "ListWithSince", Arg: "since", ArgType: ir.Named("time.Time"), Param: "since"},
			{Func: "ListWithLimit", Arg: "limit", ArgType: ir.Named("int"), Param: "limit"},
		},
	})
	require.Contains(t, got, "type ListOpt func(url.Values)\n")
	require.Contains(t, got, `q.Add("since", since.Format(time.RFC3339))`)
	require.Contains(t, got, `q.Add("limit", fmt.Sprint(limit))`)
}

func TestOperation(t *testing.T) {
	got := renderFile(t,
		ir.Operation{
			Name:        "GetFoo",
			OperationID: "get-foo",
			Method:      "GET",
			Path:        "/foo/%v",
			PathArgs:    []string{"id"},
			Params:      []ir.Param{{Name: "id", Type: ir.Named("UUID")}},
			QueryOpts:   "GetFooOpt",
			Return:      ir.Pointer{Elem: ir.Named("Foo")},
			ErrorResponses: []ir.ErrorResponse{
				{Code: 400, Type: ir.Named("GetFooResponse400")},
			},
		},
		ir.Operation{
			Name:     "ListFoos",
			Method:   "GET",
			Path:     "/foo",
			Return:   ir.Slice{Elem: ir.Named("Foo")},
			SkipAuth: true,
		},
		ir.Operation{
			Name:   "DeleteFoo",
			Method: "DELETE",
			Path:   "/foo",
		},
	)

	require.Contains(t, got, "func (c Client) GetFoo(ctx context.Context, id UUID, opts ...GetFooOpt) (*Foo, error) {")
	require.Contains(t, got, `path := fmt.Sprintf("/foo/%v", id)`)
	require.Contains(t, got, "400: func() any { return new(GetFooResponse400) },")
	require.Contains(t, got, "bodyresp := new(Foo)\n\tif err := prepareJSONResponse(response, bodyresp)")

	require.Contains(t, got, "func (c Client) ListFoos(ctx context.Context) ([]Foo, error) {")
	require.Contains(t, got, "bodyresp := []Foo{}\n\tif err := prepareJSONResponse(response, &bodyresp)")
	require.Equal(t, 2, strings.Count(got, "c.signRequest(request)"), "SkipAuth operation must not be signed")

	require.Contains(t, got, "func (c Client) DeleteFoo(ctx context.Context) error {")
	require.Contains(t, got, "_ = response.Body.Close()\n\treturn nil\n}")
}
