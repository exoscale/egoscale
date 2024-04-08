package operations

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/exoscale/egoscale/v3/generator/schemas"
)

// Generate go requests from OpenAPI spec paths operations into a go file.
func Generate(doc libopenapi.Document, path, packageName string) error {
	model, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	output.WriteString(fmt.Sprintf(`package %s
	
import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"
)
`, packageName))

	// Iterate over all paths.
	err := helpers.ForEachMapSorted(model.Model.Paths.PathItems, func(path string, item any) error {
		pathItems := item.(*v3.PathItem)
		// For each path, render each operations (GET, POST, PUT...etc) schemas and requests.
		return helpers.ForEachMapSorted(pathItems.GetOperations(), func(opName string, op any) error {
			operation := op.(*v3.Operation)
			funcName := helpers.ToCamel(operation.OperationId)
			if funcName == "" {
				funcName = helpers.ToCamel(path)
			}

			schemaResponses, err := renderResponseSchema(funcName, operation)
			if err != nil {
				return err
			}
			output.Write(schemaResponses)

			schemaRequest, err := renderRequestSchema(funcName, operation)
			if err != nil {
				return err
			}
			output.Write(schemaRequest)

			schemaParameters, err := renderRequestParametersSchema(funcName, operation)
			if err != nil {
				return err
			}
			output.Write(schemaParameters)

			request, err := serializeRequest(path, opName, funcName, operation)
			if err != nil {
				return err
			}

			if request == nil {
				return nil
			}

			m, err := renderRequest(request)
			if err != nil {
				return err
			}
			output.Write(m)

			return nil
		})
	})
	if err != nil {
		return err
	}

	if os.Getenv("GENERATOR_DEBUG") == "operations" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

// renderResponseSchema renders all schemas for every HTTP code response.
func renderResponseSchema(name string, op *v3.Operation) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	err := helpers.ForEachMapSorted(op.Responses.Codes, func(code string, v any) error {
		response := v.(*v3.Response)

		// TODO support other content type from spec.
		media, ok := response.Content["application/json"]
		if !ok {
			return nil
		}

		findable, err := renderFindable(name, media.Schema)
		if err != nil {
			return err
		}

		// Skip on reference.
		if media.Schema.IsReference() {
			return nil
		}

		// Skip on array referencing a schema.
		_, ok = isArrayReference(media.Schema)
		if ok {
			return nil
		}

		schemaResp, err := schemas.RenderSchema(name+"Response", media.Schema)
		if err != nil {
			return err
		}
		output.Write(schemaResp)
		if findable != nil {
			output.Write(findable)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// renderRequestSchema renders request body schemas, mostly for HTTP POST and PUT.
// It returns a nil output if there is no request schema to render for a given operation.
func renderRequestSchema(name string, op *v3.Operation) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	if op.RequestBody == nil {
		return nil, nil
	}

	// TODO support other content type from spec.
	media, ok := op.RequestBody.Content["application/json"]
	if !ok {
		return nil, nil
	}

	// Skip on reference.
	if media.Schema.IsReference() {
		return nil, nil
	}

	// Skip on array referencing a schema.
	_, ok = isArrayReference(media.Schema)
	if ok {
		return nil, nil
	}

	schemaResp, err := schemas.RenderSchema(name+"Request", media.Schema)
	if err != nil {
		return nil, err
	}
	output.Write(schemaResp)

	return output.Bytes(), nil
}

const queryParamTemplate = `
func {{ .FuncName }}({{ .ParamName }} {{ .ParamType }}) {{ .FuncReturn }} {
	return func(q url.Values) {
		q.Add("{{ .ParamName }}", fmt.Sprint({{ .ParamName }}))
	}
}
`

type QueryParam struct {
	FuncName   string
	ParamName  string
	ParamType  string
	FuncReturn string
}

// renderRequestParametersSchema renders the schemas for optional query params and path params.
func renderRequestParametersSchema(name string, op *v3.Operation) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	query := bytes.NewBuffer([]byte{})

	someQueryParam := false
	for _, p := range op.Parameters {
		s := p.Schema.Schema()
		if s == nil {
			continue
		}

		ParamTypeName := name + helpers.ToCamel(p.Name)

		if p.In == "query" {
			t, err := template.New("queryParam").Parse(queryParamTemplate)
			if err != nil {
				return nil, err
			}

			typ := ParamTypeName
			if schemas.IsSimpleSchema(s) && len(s.Enum) == 0 {
				typ = schemas.RenderSimpleType(s)
			}
			if err := t.Execute(query, QueryParam{
				FuncName:   name + "With" + helpers.ToCamel(p.Name),
				ParamName:  helpers.ToLowerCamel(p.Name),
				ParamType:  typ,
				FuncReturn: name + "Opt",
			}); err != nil {
				return nil, err
			}
			someQueryParam = true
		}

		// Skip simple types, not needed to be rendered.
		if schemas.IsSimpleSchema(s) && len(s.Enum) == 0 {
			continue
		}

		// As long as an HTTP query param and path param not using objects or arrays in our spec,
		// this code path is called only for string enum types.
		// TODO: To support array or object, add a .String() method to those types for marshalling like described here:
		// https://swagger.io/docs/specification/describing-parameters/#path-parameters
		// https://swagger.io/docs/specification/describing-parameters/#query-parameters

		if len(s.Enum) == 0 {
			slog.Warn(
				"object/array as query/path params are not implemented",
				slog.String("request", name),
				slog.String("param", ParamTypeName),
			)
		}

		sc, err := schemas.RenderSchema(ParamTypeName, p.Schema)
		if err != nil {
			return nil, err
		}
		output.Write(sc)
	}

	if someQueryParam {
		q := append([]byte(fmt.Sprintf("type %sOpt func(url.Values)\n", name)), query.Bytes()...)
		output.Write(q)
	}

	return output.Bytes(), nil
}

const findableTemplate = `
// Find{{ .TypeName }} attempts to find an {{ .TypeName }} by name or ID.
func (l {{ .ListTypeName }}) Find{{ .TypeName }}(nameOrID string) ({{ .TypeName }}, error) {
	for i, elem := range l.{{ .ListFieldName }} {
		if elem.Name == nameOrID || elem.ID.String() == nameOrID {
			return l.{{ .ListFieldName }}[i], nil
		}
	}

	return {{ .TypeName }}{}, fmt.Errorf("%q not found in {{ .ListTypeName }}: %w", nameOrID, ErrNotFound)
}
`

type Findable struct {
	TypeName      string
	ListTypeName  string
	ListFieldName string
}

// renderFindable renders a find method on listable resource.
// this find method get the resource by name or id if available.
// returns nil on non listable resources.
func renderFindable(funcName string, s *base.SchemaProxy) ([]byte, error) {
	sc, err := s.BuildSchema()
	if err != nil {
		return nil, err
	}
	schemas.InferType(sc)

	// Check if listable response.
	if !strings.HasPrefix(strings.ToLower(funcName), "list") {
		return nil, nil
	}

	if len(sc.Type) > 0 && sc.Type[0] != "object" {
		return nil, nil
	}

	for propName, propSc := range sc.Properties {
		prop, err := propSc.BuildSchema()
		if err != nil {
			return nil, err
		}
		schemas.InferType(prop)

		if len(prop.Type) > 0 && prop.Type[0] != "array" {
			continue
		}

		if prop.Items == nil {
			continue
		}
		if !prop.Items.IsA() {
			continue
		}

		item, err := prop.Items.A.BuildSchema()
		if err != nil {
			return nil, err
		}

		typeName := funcName + "Response" + helpers.ToCamel(propName)
		if prop.Items.A.IsReference() {
			typeName = helpers.RenderReference(prop.Items.A.GetReference())
		}

		_, hasName := item.Properties["name"]
		_, hasID := item.Properties["id"]
		if hasName && hasID {
			output := bytes.NewBuffer([]byte{})
			t, err := template.New("Findable").Parse(findableTemplate)
			if err != nil {
				return nil, err
			}
			if err := t.Execute(output, Findable{
				ListTypeName:  funcName + "Response",
				ListFieldName: helpers.ToCamel(propName),
				TypeName:      typeName,
			}); err != nil {
				return nil, err
			}

			return output.Bytes(), nil
		}
	}

	return nil, nil
}

type RequestTmpl struct {
	Comment        string
	Name           string
	OperationId    string
	Params         string
	ValueReturn    string
	URLPathBuilder string
	HTTPMethod     string
	BodyRequest    bool
	BodyRespType   string
	ContentType    string
	QueryParams    map[string]string
}

// serializeRequest serializes the openAPI spec into the request template.
func serializeRequest(path, httpMethod, funcName string, op *v3.Operation) (*RequestTmpl, error) {
	p := RequestTmpl{
		Name:        funcName,
		OperationId: op.OperationId,
		HTTPMethod:  strings.ToUpper(httpMethod),
	}
	p.Comment = renderDoc(op)
	params := getParameters(op, funcName)
	p.Params = strings.Join(params, ", ")
	valuesReturn := getValuesReturn(op, funcName)
	if len(valuesReturn) == 2 {
		p.BodyRespType = valuesReturn[0]
		if !strings.HasPrefix(valuesReturn[0], "[]") {
			p.BodyRespType = "&" + p.BodyRespType[1:]
		}
	}
	p.ValueReturn = fmt.Sprintf("(%s)", strings.Join(valuesReturn, ", "))
	// This should never happen in our Exoscale API spec.
	// This is here as a reminder the day we add such a behavior in the OpenAPI spec.
	if p.ValueReturn == "(error)" {
		slog.Error(
			"single error value return not implemented",
			slog.String("path", path),
			slog.String("operation", funcName),
		)
		return nil, nil
	}
	p.URLPathBuilder = renderURLPathBuilder(path, op)

	if op.RequestBody != nil {
		p.BodyRequest = true
		//TODO: manage other content type from spec.
		p.ContentType = "application/json"
	}

	p.QueryParams = getQueryParams(op)

	return &p, nil
}

// renderRequest using the request.tmpl.
func renderRequest(m *RequestTmpl) ([]byte, error) {
	t, err := template.New("request.tmpl").ParseFiles("./operations/request.tmpl")
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func getParameters(op *v3.Operation, funcName string) (params []string) {
	params = []string{"ctx context.Context"}

	someQueryParam := false
	if len(op.Parameters) > 0 {
		for _, param := range op.Parameters {
			s := param.Schema.Schema()
			if s == nil {
				continue
			}

			// Only concat other parameters first,
			// Since query param is set as vaarg parameter.
			if param.In == "query" {
				someQueryParam = true
				continue
			}

			name := helpers.ToLowerCamel(param.Name)

			// https://github.com/exoscale/entities/commit/dda7e9f52ded1879e509d465555023b5a16d0155
			if strings.Contains(name, "*") {
				slog.Warn(
					"parameter name contains '*' in spec",
					slog.String("operation", funcName),
					slog.String("param", name),
				)
				name = strings.Trim(name, "*")
			}

			ptr := ""
			if !param.Required {
				slog.Warn(
					"path parameter not required in spec",
					slog.String("operation", funcName),
					slog.String("param", name),
				)
				// XXX: we should never handle this case in our spec
				// since optional param are query param and path param are always required.
				// https://swagger.io/docs/specification/describing-parameters/#path-parameters
				ptr = "*"
			}

			if !schemas.IsSimpleSchema(s) || len(s.Enum) > 0 {
				params = append(params, name+" "+ptr+funcName+helpers.ToCamel(name))
				continue
			}

			params = append(params, name+" "+ptr+schemas.RenderSimpleType(s))
		}
	}

	// Add variadic arguments to the end
	if someQueryParam {
		defer func() {
			params = append(params, fmt.Sprintf("opts ...%sOpt", funcName))
		}()
	}

	if op.RequestBody == nil {
		return
	}

	// TODO support other content type from OpenAPI spec.
	media, ok := op.RequestBody.Content["application/json"]
	if !ok {
		return
	}
	if media.Schema.IsReference() {
		params = append(params, "req "+helpers.RenderReference(media.Schema.GetReference()))
		return
	}

	params = append(params, "req "+funcName+"Request")

	return params
}

func getValuesReturn(op *v3.Operation, funcName string) (values []string) {
	values = []string{}
	defer func() {
		values = append(values, "error")
	}()

	if len(op.Responses.Codes) == 0 {
		return values
	}

	for k, v := range op.Responses.Codes {
		// We support only 200 return as body reply in our OpenAPI spec.
		// Skip other HTTP response code.
		if k != "200" {
			continue
		}

		media, ok := v.Content["application/json"]
		if !ok {
			continue
		}
		if media.Schema.IsReference() {
			values = append(values, "*"+helpers.RenderReference(media.Schema.GetReference()))
			return values
		}

		a, ok := isArrayReference(media.Schema)
		if ok {
			values = append(values, a)
			return values
		}
	}

	values = append(values, "*"+funcName+"Response")
	return values
}

func renderDoc(op *v3.Operation) string {
	doc := op.Description
	if doc == "" {
		doc = op.Summary
	}

	return helpers.RenderDoc(doc)
}

// renderURLPathBuilder renders the sprintf code used to build the path request in request template.
func renderURLPathBuilder(rawPath string, op *v3.Operation) string {
	if len(op.Parameters) == 0 {
		return fmt.Sprintf("%q", rawPath)
	}
	path := rawPath
	sprintfParam := []string{}
	for _, p := range op.Parameters {
		if p.In != "path" {
			continue
		}

		path = strings.Replace(path, "{"+p.Name+"}", "%v", 1)
		sprintfParam = append(
			sprintfParam,
			// https://github.com/exoscale/entities/commit/dda7e9f52ded1879e509d465555023b5a16d0155
			helpers.ToLowerCamel(strings.Trim(p.Name, "*")),
		)
	}
	if path == rawPath {
		return fmt.Sprintf("%q", rawPath)
	}

	return `fmt.Sprintf("` + path + `", ` + strings.Join(sprintfParam, ", ") + ")"
}

func getQueryParams(op *v3.Operation) map[string]string {
	if len(op.Parameters) == 0 {
		return nil
	}

	result := make(map[string]string)
	for _, p := range op.Parameters {
		if p.In != "query" {
			continue
		}
		result[p.Name] = helpers.ToLowerCamel(p.Name)
	}
	if len(result) == 0 {
		return nil
	}

	return result
}

// isArrayReference returns true if it's an array pointing on a schema reference.
// Returns a formatted type corresponding to it on true.
func isArrayReference(sp *base.SchemaProxy) (string, bool) {
	s := sp.Schema()
	if s == nil {
		return "", false
	}

	if s.Type[0] == "array" {
		if s.Items == nil {
			return "", false
		}
		if !s.Items.IsA() {
			return "", false
		}

		isReference := s.Items.A.IsReference()
		if isReference {
			return "[]" + helpers.RenderReference(s.Items.A.GetReference()), true
		}
	}

	return "", false
}
