package build

import (
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/exoscale/egoscale/v3/generator/ir"
	"github.com/exoscale/egoscale/v3/generator/naming"
)

// Operations builds the file of the spec path operations: a Client method per operation,
// preceded by its inline response, request and parameter types.
func (b *Builder) Operations(packageName string) (ir.File, error) {
	if b.model.Paths == nil || orderedmap.Len(b.model.Paths.PathItems) == 0 {
		return ir.File{}, errors.New("no path items defined in the spec")
	}

	for pair := orderedmap.SortAlpha(b.model.Paths.PathItems).First(); pair != nil; pair = pair.Next() {
		path, item := pair.Key(), pair.Value()
		for pair := orderedmap.SortAlpha(item.GetOperations()).First(); pair != nil; pair = pair.Next() {
			method, op := pair.Key(), pair.Value()
			if err := b.operation(path, method, op); err != nil {
				return ir.File{}, fmt.Errorf("%s %s: %w", strings.ToUpper(method), path, err)
			}
		}
	}

	return b.file(packageName, "context", "fmt", "net", "net/http", "net/url", "time"), nil
}

func (b *Builder) operation(path, method string, op *v3.Operation) error {
	funcName := b.names.Camel(op.OperationId)
	if funcName == "" {
		funcName = b.names.Camel(path)
	}
	origin := "operation " + op.OperationId

	if err := b.responseSchemas(funcName, op, origin); err != nil {
		return err
	}
	if err := b.requestSchema(funcName, op, origin); err != nil {
		return err
	}
	queryOpts, err := b.parameterSchemas(funcName, op, origin)
	if err != nil {
		return err
	}

	params, err := b.parameters(funcName, op)
	if err != nil {
		return err
	}
	o := ir.Operation{
		Doc:            operationDoc(op),
		Name:           funcName,
		OperationID:    op.OperationId,
		Method:         strings.ToUpper(method),
		Params:         params,
		Body:           b.requestBody(funcName, op),
		QueryOpts:      queryOpts,
		Return:         b.returnType(funcName, op),
		ErrorResponses: b.errorResponses(funcName, op),
		SkipAuth:       b.skipAuth(op.OperationId),
	}
	o.Path, o.PathArgs = b.urlPath(path, op)

	return b.emit(o, origin)
}

// responseSchemas emits the inline response schemas of every HTTP code,
// and the Find method of a list response.
func (b *Builder) responseSchemas(funcName string, op *v3.Operation, origin string) error {
	if op.Responses == nil {
		return nil
	}

	for pair := op.Responses.Codes.First(); pair != nil; pair = pair.Next() {
		code, response := pair.Key(), pair.Value()
		if _, err := strconv.Atoi(code); err != nil {
			slog.Warn(
				"non numeric HTTP response code not implemented",
				slog.String("operation", funcName),
				slog.String("code", code),
			)
			continue
		}

		media, ok := jsonMedia(responseContent(response))
		if !ok {
			continue
		}

		// Findable is only rendered on the HTTP 200 list response.
		var findable *ir.Findable
		if code == "200" {
			var err error
			if findable, err = b.findable(funcName, media.Schema); err != nil {
				return err
			}
		}

		switch {
		case media.Schema.IsReference():
			// A $ref schema name ending with "-response" is rendered as funcName+"Response"
			// (e.g. list-ai-api-keys-response -> ListAIAPIKeysResponse): findable is still valid.
			// Otherwise (e.g. dbaas-clickhouse-roles), no funcName+"Response" type exists.
			if !strings.HasSuffix(strings.ToLower(media.Schema.GetReference()), "-response") {
				continue
			}
		case b.arrayReference(media.Schema) != nil:
			continue
		default:
			s, err := media.Schema.BuildSchema()
			if err != nil {
				return fmt.Errorf("response %s: %w", code, err)
			}
			if err := b.schema(responseTypeName(funcName, code), s, origin+" response "+code); err != nil {
				return err
			}
		}

		if findable != nil {
			if err := b.emit(*findable, origin); err != nil {
				return err
			}
		}
	}

	return nil
}

// responseTypeName returns the go type name of an inline response schema for an HTTP code.
// HTTP 200 keeps the <funcName>Response name, other codes are suffixed with the code
// (e.g. <funcName>Response400).
func responseTypeName(funcName, code string) string {
	if code == "200" {
		return funcName + "Response"
	}

	return funcName + "Response" + code
}

// responseType returns the go type of a response schema for an HTTP code:
// the referenced type for a $ref, a slice for an array of $ref,
// or the inline response schema type (see responseTypeName).
func (b *Builder) responseType(funcName, code string, media *v3.MediaType) ir.Type {
	if media.Schema.IsReference() {
		return ir.Pointer{Elem: b.ref(media.Schema.GetReference(), "")}
	}
	if slice := b.arrayReference(media.Schema); slice != nil {
		return slice
	}

	return ir.Pointer{Elem: ir.Named(responseTypeName(funcName, code))}
}

// returnType returns the type of the HTTP 200 response body if present,
// otherwise of the lowest other 2xx HTTP code with a body. Nil without 2xx body.
func (b *Builder) returnType(funcName string, op *v3.Operation) ir.Type {
	if op.Responses == nil {
		return nil
	}

	successCode := 0
	var successMedia *v3.MediaType
	for pair := op.Responses.Codes.First(); pair != nil; pair = pair.Next() {
		code, err := strconv.Atoi(pair.Key())
		if err != nil || code < 200 || code > 299 {
			continue
		}

		media, ok := jsonMedia(responseContent(pair.Value()))
		if !ok {
			continue
		}

		if successMedia != nil {
			slog.Warn(
				"multiple 2xx HTTP response bodies, only one is returned",
				slog.String("operation", funcName),
				slog.Int("code", code),
			)
		}

		if successMedia == nil || code == 200 || (successCode != 200 && code < successCode) {
			successCode, successMedia = code, media
		}
	}

	if successMedia == nil {
		return nil
	}

	return b.responseType(funcName, strconv.Itoa(successCode), successMedia)
}

// errorResponses returns the body types of the 4xx and 5xx HTTP responses,
// used to decode the typed body of an APIError. Sorted by HTTP code.
func (b *Builder) errorResponses(funcName string, op *v3.Operation) []ir.ErrorResponse {
	if op.Responses == nil {
		return nil
	}

	var result []ir.ErrorResponse
	for pair := op.Responses.Codes.First(); pair != nil; pair = pair.Next() {
		code, err := strconv.Atoi(pair.Key())
		if err != nil || code < 400 || code > 599 {
			continue
		}

		media, ok := jsonMedia(responseContent(pair.Value()))
		if !ok {
			continue
		}

		typ := b.responseType(funcName, pair.Key(), media)
		if p, ok := typ.(ir.Pointer); ok {
			typ = p.Elem
		}
		result = append(result, ir.ErrorResponse{Code: code, Type: typ})
	}

	sort.Slice(result, func(i, j int) bool { return result[i].Code < result[j].Code })

	return result
}

// requestSchema emits the inline request body schema, mostly for HTTP POST and PUT.
func (b *Builder) requestSchema(funcName string, op *v3.Operation, origin string) error {
	if op.RequestBody == nil {
		return nil
	}

	media, ok := jsonMedia(op.RequestBody.Content)
	if !ok || media.Schema.IsReference() || b.arrayReference(media.Schema) != nil {
		return nil
	}

	s, err := media.Schema.BuildSchema()
	if err != nil {
		return fmt.Errorf("request body: %w", err)
	}

	return b.schema(funcName+"Request", s, origin+" request body")
}

// requestBody returns the request body type, nil if none.
func (b *Builder) requestBody(funcName string, op *v3.Operation) ir.Type {
	if op.RequestBody == nil {
		return nil
	}

	// TODO support other content type from OpenAPI spec.
	media, ok := jsonMedia(op.RequestBody.Content)
	if !ok {
		return nil
	}
	if media.Schema.IsReference() {
		return b.ref(media.Schema.GetReference(), "")
	}

	return ir.Named(funcName + "Request")
}

// parameterSchemas emits the enum types of the query and path parameters,
// and the query options. It returns the query options type name, "" if none.
func (b *Builder) parameterSchemas(funcName string, op *v3.Operation, origin string) (string, error) {
	queryOpts := ir.QueryOpts{Type: funcName + "Opt"}

	for _, p := range op.Parameters {
		s := p.Schema.Schema()
		if s == nil {
			continue
		}

		paramTypeName := funcName + b.names.Camel(p.Name)
		simple := isSimple(s) && len(s.Enum) == 0

		if p.In == "query" {
			var typ ir.Type = ir.Named(paramTypeName)
			if simple {
				var err error
				if typ, err = simpleType(s); err != nil {
					return "", fmt.Errorf("query param %s: %w", p.Name, err)
				}
			}
			queryOpts.Funcs = append(queryOpts.Funcs, ir.QueryOpt{
				Func:    funcName + "With" + b.names.Camel(p.Name),
				Arg:     b.names.LowerCamel(p.Name),
				ArgType: typ,
				Param:   p.Name,
			})
		}

		// Simple types have no declaration.
		if simple {
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
				slog.String("request", funcName),
				slog.String("param", paramTypeName),
			)
		}

		if err := b.schema(paramTypeName, s, origin+" param "+p.Name); err != nil {
			return "", err
		}
	}

	if len(queryOpts.Funcs) == 0 {
		return "", nil
	}

	return queryOpts.Type, b.emit(queryOpts, origin)
}

// parameters returns the method parameters of the path parameters.
func (b *Builder) parameters(funcName string, op *v3.Operation) ([]ir.Param, error) {
	var params []ir.Param
	for _, p := range op.Parameters {
		s := p.Schema.Schema()
		if s == nil || p.In == "query" {
			continue
		}

		name := b.names.LowerCamel(p.Name)
		// https://github.com/exoscale/entities/commit/dda7e9f52ded1879e509d465555023b5a16d0155
		if strings.Contains(name, "*") {
			slog.Warn(
				"parameter name contains '*' in spec",
				slog.String("operation", funcName),
				slog.String("param", name),
			)
			name = strings.Trim(name, "*")
		}

		var typ ir.Type = ir.Named(funcName + b.names.Camel(name))
		if isSimple(s) && len(s.Enum) == 0 {
			var err error
			if typ, err = simpleType(s); err != nil {
				return nil, fmt.Errorf("param %s: %w", p.Name, err)
			}
		}

		if p.Required != nil && !*p.Required {
			// XXX: we should never handle this case in our spec
			// since optional param are query param and path param are always required.
			// https://swagger.io/docs/specification/describing-parameters/#path-parameters
			slog.Warn(
				"path parameter not required in spec",
				slog.String("operation", funcName),
				slog.String("param", name),
			)
			typ = ir.Pointer{Elem: typ}
		}

		params = append(params, ir.Param{Name: name, Type: typ})
	}

	return params, nil
}

// urlPath returns the URL path with %v in place of path parameters, and the parameters.
func (b *Builder) urlPath(rawPath string, op *v3.Operation) (string, []string) {
	path := rawPath
	var args []string
	for _, p := range op.Parameters {
		if p.In != "path" {
			continue
		}

		path = strings.Replace(path, "{"+p.Name+"}", "%v", 1)
		// https://github.com/exoscale/entities/commit/dda7e9f52ded1879e509d465555023b5a16d0155
		args = append(args, b.names.LowerCamel(strings.Trim(p.Name, "*")))
	}
	if path == rawPath {
		return rawPath, nil
	}

	return path, args
}

// findable returns the Find method of a list response schema, nil if not listable:
// the first array property of objects having a name, id or x-go-findable property.
func (b *Builder) findable(funcName string, proxy *base.SchemaProxy) (*ir.Findable, error) {
	if !strings.HasPrefix(strings.ToLower(funcName), "list") {
		return nil, nil
	}

	s, err := proxy.BuildSchema()
	if err != nil {
		return nil, err
	}
	if k := kindOf(s); k != "" && k != kindObject {
		return nil, nil
	}
	if orderedmap.Len(s.Properties) == 0 {
		return nil, nil
	}

	for pair := s.Properties.First(); pair != nil; pair = pair.Next() {
		propName, propProxy := pair.Key(), pair.Value()
		prop, err := propProxy.BuildSchema()
		if err != nil {
			return nil, err
		}
		if k := kindOf(prop); k != "" && k != kindArray {
			continue
		}
		if prop.Items == nil || !prop.Items.IsA() {
			continue
		}

		item, err := prop.Items.A.BuildSchema()
		if err != nil {
			return nil, err
		}
		if item.Properties == nil {
			continue
		}

		itemType := funcName + "Response" + b.names.Camel(propName)
		if prop.Items.A.IsReference() {
			itemType = b.ref(prop.Items.A.GetReference(), "").String()
		}

		// x-go-findable "1" and "2" override the name and id fields.
		var field1, field2 string
		if _, ok := item.Properties.Get("name"); ok {
			field1 = "name"
		}
		if _, ok := item.Properties.Get("id"); ok {
			field2 = "id"
		}
		for pair := item.Properties.First(); pair != nil; pair = pair.Next() {
			sc := pair.Value().Schema()
			if sc == nil || sc.Extensions == nil {
				continue
			}
			switch val, _ := sc.Extensions.Get("x-go-findable"); {
			case val == nil:
			case val.Value == "1":
				field1 = pair.Key()
			case val.Value == "2":
				field2 = pair.Key()
			}
		}

		f := &ir.Findable{
			ListType:  funcName + "Response",
			ListField: b.names.Camel(propName),
			ItemType:  itemType,
		}
		switch {
		case field1 != "" && field2 != "":
			f.Param = b.names.LowerCamel(field1) + "Or" + b.names.Camel(field2)
			f.Fields = []string{b.names.Camel(field1), b.names.Camel(field2)}
		case field1 != "":
			f.Param = b.names.LowerCamel(field1)
			f.Fields = []string{b.names.Camel(field1)}
		case field2 != "":
			f.Param = b.names.LowerCamel(field2)
			f.Fields = []string{b.names.Camel(field2)}
		default:
			continue
		}

		return f, nil
	}

	return nil, nil
}

// arrayReference returns the slice type of an array of schema reference, nil otherwise.
func (b *Builder) arrayReference(proxy *base.SchemaProxy) ir.Type {
	s := proxy.Schema()
	if s == nil || kindOf(s) != kindArray {
		return nil
	}
	if s.Items == nil || !s.Items.IsA() || !s.Items.A.IsReference() {
		return nil
	}

	return ir.Slice{Elem: b.ref(s.Items.A.GetReference(), "")}
}

// jsonMedia returns the application/json media type of a content if any.
func jsonMedia(content *orderedmap.Map[string, *v3.MediaType]) (*v3.MediaType, bool) {
	if content == nil {
		return nil, false
	}
	// TODO support other content type from spec.
	media, ok := content.Get("application/json")
	if !ok || media == nil || media.Schema == nil {
		return nil, false
	}

	return media, true
}

func responseContent(r *v3.Response) *orderedmap.Map[string, *v3.MediaType] {
	if r == nil {
		return nil
	}

	return r.Content
}

func operationDoc(op *v3.Operation) string {
	doc := op.Description
	if doc == "" {
		doc = op.Summary
	}

	return naming.Doc(doc)
}
