package operations

import (
	"bytes"
	"fmt"
	"go/format"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/exoscale/egoscale/v3/generator/schemas"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"context"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/go-playground/validator/v10"
)
`, packageName))

	methods := []Method{}

	err := helpers.ForEachMapSorted(model.Model.Paths.PathItems, func(path string, item any) error {
		pathItems := item.(*v3.PathItem)
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

			schemaParameters, err := renderMethodParameterSchema(funcName, operation)
			if err != nil {
				return err
			}
			output.Write(schemaParameters)

			method, err := serializeMethod(path, opName, funcName, operation)
			if err != nil {
				return err
			}

			methods = append(methods, *method)

			m, err := renderMethod(method)
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

	// Add extra methods here from static generated (operationUtils).
	// TODO put this static list somewhere.
	methods = append(methods, Method{Name: "Wait", Params: "context.Context, *Operation, ...OperationState", ValueReturn: "(*Operation, error)"})
	methods = append(methods, Method{Name: "WithZone", Params: "APIZone", ValueReturn: "Client"})
	methods = append(methods, Method{Name: "WithContext", Params: "context.Context", ValueReturn: "Client"})
	methods = append(methods, Method{Name: "WithHttpClient", Params: " *http.Client", ValueReturn: "Client"})
	methods = append(methods, Method{Name: "WithRequestMiddleware", Params: "RequestMiddlewareFn", ValueReturn: "Client"})

	iface, err := renderInterface(methods)
	if err != nil {
		return err
	}
	output.Write(iface)
	output.WriteString(operationUtils)

	if os.Getenv("GENERATOR_DEBUG") == "operations" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

func renderResponseSchema(name string, op *v3.Operation) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})
	err := helpers.ForEachMapSorted(op.Responses.Codes, func(code string, v any) error {
		response := v.(*v3.Response)

		// TODO support other content type from spec.
		media, ok := response.Content["application/json"]
		if !ok {
			return nil
		}
		if media.Schema.IsReference() {
			return nil
		}

		ok, _ = isArrayReference(media.Schema)
		if ok {
			return nil
		}

		schemaResp, err := schemas.RenderSchema(name+"Response", media.Schema)
		if err != nil {
			return err
		}
		output.Write(schemaResp)

		return nil
	})
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

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

	if media.Schema.IsReference() {
		return nil, nil
	}

	ok, _ = isArrayReference(media.Schema)
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
func {{ .MethodName }}({{ .ParamName }} {{ .ParamType }}) {{ .MethodReturn }} {
	return func(q url.Values) {
		q.Add("{{ .ParamName }}", fmt.Sprint({{ .ParamName }}))
	}
}
`

type QueryParam struct {
	MethodName   string
	ParamName    string
	ParamType    string
	MethodReturn string
}

func renderMethodParameterSchema(name string, op *v3.Operation) ([]byte, error) {
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
				MethodName:   name + "With" + helpers.ToCamel(p.Name),
				ParamName:    helpers.ToLowerCamel(p.Name),
				ParamType:    typ,
				MethodReturn: name + "Opt",
			}); err != nil {
				return nil, err
			}
			someQueryParam = true
		}

		if schemas.IsSimpleSchema(s) && len(s.Enum) == 0 {
			continue
		}

		//TODO: build String() methone for not simple type.
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

func serializeMethod(path, httpMethode, funcName string, op *v3.Operation) (*Method, error) {
	p := Method{
		// Generate code only for prod for now
		APIName:    "API",
		Name:       funcName,
		HTTPMethod: strings.ToUpper(httpMethode),
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
	// This should never happen in our Exoscale API.
	// This panic is here as a reminder the day we add such a behavior in the OpenAPI spec.
	if p.ValueReturn == "(error)" {
		panic("single error value return not implemented")
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

func renderMethod(m *Method) ([]byte, error) {
	t, err := template.New("Method").Parse(methodTemplate)
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

			if param.In == "query" {
				someQueryParam = true
				continue
			}

			name := helpers.ToLowerCamel(param.Name)

			// This is WIP to be removed in:
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
				// since optional param are query param and path param are alwayse required.
				ptr = "*"
			}

			if !schemas.IsSimpleSchema(s) || len(s.Enum) > 0 {
				params = append(params, name+" "+ptr+funcName+helpers.ToCamel(name))
				continue
			}

			params = append(params, name+" "+ptr+schemas.RenderSimpleType(s))
		}
	}
	// Add vaargs to the end
	if someQueryParam {
		defer func() {
			params = append(params, fmt.Sprintf("opts ...%sOpt", funcName))
		}()
	}

	if op.RequestBody == nil {
		return
	}

	// TODO support other content type from spec.
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

		ok, a := isArrayReference(media.Schema)
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
	if doc == "null" {
		return ""
	}

	docs := strings.Split(doc, "\n")
	r := []string{}
	for i, d := range docs {
		if d == "" {
			docs = append(docs[:i], docs[i+1:]...)
			continue
		}
		r = append(r, "// "+strings.TrimSpace(d))
	}

	if len(r) == 0 {
		return ""
	}

	return strings.Join(r, "\n")
}

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
			// This is WIP to be removed in:
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

type Method struct {
	Comment        string
	APIName        string
	Name           string
	Params         string
	ValueReturn    string
	URLPathBuilder string
	HTTPMethod     string
	BodyRequest    bool
	BodyRespType   string
	ContentType    string
	QueryParams    map[string]string
}

const methodTemplate = `
{{ .Comment }}
func (c Client{{ .APIName }}) {{ .Name }}({{ .Params }}) {{ .ValueReturn }} {
	path := {{ .URLPathBuilder }}

	{{ if .BodyRequest }}
	body, err := prepareJsonBody(req)
	if err != nil {
		return nil, fmt.Errorf("{{ .Name }}: prepare Json body: %w", err)
	}
	{{ end }}

	request, err := http.NewRequestWithContext(ctx, "{{ .HTTPMethod }}", c.serverURL + path, {{ if .BodyRequest }}body{{else}}nil{{end}})
	if err != nil {
		return nil, fmt.Errorf("{{ .Name }}: new request: %w", err)
	}

	{{ if ne .QueryParams nil }}if len(opts) > 0 {
		q := request.URL.Query()
		for _, opt := range opts {
			opt(q)
		}
		request.URL.RawQuery = q.Encode()
	}{{ end }}

	{{ if ne .ContentType "" }}
	request.Header.Add("Content-Type", "{{ .ContentType }}")
	{{ end }}

	if err := registerRequestMiddlewares(&c, ctx, request); err != nil {
		return nil, fmt.Errorf("{{ .Name }}: register middlewares: %w", err)
	}

	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("{{ .Name }}: http client do: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return nil, fmt.Errorf("request {{ .Name }} returned %d code", resp.StatusCode)
	}

	bodyresp := {{ .BodyRespType }}{}
	if err := prepareJsonResponse(resp, bodyresp); err != nil {
		return nil, fmt.Errorf("{{ .Name }}: prepare Json response: %w", err)
	}

	return bodyresp, nil
}
`

type InterfaceModel struct {
	Methods []Method
}

func renderInterface(methods []Method) ([]byte, error) {
	t, err := template.New("Interface").Parse(interfaceTemplate)
	if err != nil {
		return nil, err
	}

	buf := bytes.NewBuffer([]byte{})
	if err := t.Execute(buf, InterfaceModel{Methods: methods}); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

const interfaceTemplate = `
// Client represents an API interface representation for ClientAPI and mock test.
type Client interface {
	{{ range .Methods }}{{ .Comment }}
	{{ .Name }}({{ .Params }}) {{ .ValueReturn }}
	{{end}}
}
`

const operationUtils = `
// Wait is a helper that waits for async operation to reach the final state.
// Final states are one of: failure, success, timeout.
// If states argument are given, returns an error if the final state not match on of those.
func (c ClientAPI) Wait(ctx context.Context, op *Operation, states ...OperationState) (*Operation, error) {
	if op == nil {
		return nil, fmt.Errorf("operation is nil")
	}

	ticker := time.NewTicker(c.pollingInterval)
	defer ticker.Stop()

	if op.State != OperationStatePending {
		return op, nil
	}

	var operation *Operation
polling:
	for {
		select {
		case <-ticker.C:
			o, err := c.GetOperation(ctx, op.ID)
			if err != nil {
				return nil, err
			}
			if o.State == OperationStatePending {
				continue
			}

			operation = o
			break polling
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}

	if len(states) == 0 {
		return operation, nil
	}

	for _, st := range states {
		if operation.State == st {
			return operation, nil
		}
	}

	var ref OperationReference
	if operation.Reference != nil {
		ref = *operation.Reference
	}

	return nil,
		fmt.Errorf("operation: %q %v, state: %s, reason: %q, message: %q",
			operation.ID,
			ref,
			operation.State,
			operation.Reason,
			operation.Message,
		)
}

func String(s string) *string {
	return &s
}

func Int64(i int64) *int64 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

// Validate any struct from schema or request
func Validate(r any) error {
	return validator.New().Struct(r)
}

func prepareJsonBody(body any) (*bytes.Reader, error) {
	buf, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf), nil
}

func prepareJsonResponse(resp *http.Response, v any) error {
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(buf, v); err != nil {
		return err
	}

	return nil
}

type UUID string

func (u UUID) String() string {
	return string(u)
}

func ParseUUID(s string) (UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return "", err
	}

	return UUID(id.String()), nil
}
`

func isArrayReference(sp *base.SchemaProxy) (bool, string) {
	s := sp.Schema()
	if s == nil {
		return false, ""
	}

	if s.Type[0] == "array" {
		if s.Items == nil {
			return false, ""
		}
		if !s.Items.IsA() {
			return false, ""
		}

		isReference := s.Items.A.IsReference()
		if isReference {
			return true, "[]" + helpers.RenderReference(s.Items.A.GetReference())
		}
	}

	return false, ""
}
