// Package ir is the generator intermediate representation: the Go declarations
// to generate, built from the OpenAPI spec by the build package and turned
// into Go source by the render package. It does not depend on libopenapi.
package ir

// File is a generated Go source file.
type File struct {
	Package string
	Imports []string
	// Decls are rendered in order.
	Decls []Decl
}

// Type is a Go type expression.
type Type interface {
	String() string
}

// Named is a declared or predeclared type, e.g. "string", "time.Time" or "Instance".
type Named string

func (n Named) String() string { return string(n) }

// Any is the empty interface type.
const Any = Named("any")

// Pointer is *Elem.
type Pointer struct{ Elem Type }

func (p Pointer) String() string { return "*" + p.Elem.String() }

// Slice is []Elem.
type Slice struct{ Elem Type }

func (s Slice) String() string { return "[]" + s.Elem.String() }

// Map is map[string]Elem.
type Map struct{ Elem Type }

func (m Map) String() string { return "map[string]" + m.Elem.String() }

// Decl is a top level Go declaration.
type Decl interface {
	decl()
}

// TypeDef is: type Name Type.
type TypeDef struct {
	Doc  string
	Name string
	Type Type
}

// Enum is a type definition with its constant values.
type Enum struct {
	Doc    string
	Name   string
	Base   Type
	Values []EnumValue
}

// EnumValue is a constant of an Enum.
type EnumValue struct {
	Name string
	// Literal is the Go literal of the value, e.g. `"running"` or `42`.
	Literal string
}

// Struct is a struct type definition.
type Struct struct {
	Doc    string
	Name   string
	Fields []Field
}

// Field is a JSON struct field.
type Field struct {
	Doc  string
	Name string
	Type Type
	// JSON is the JSON property name.
	JSON      string
	OmitEmpty bool
	// Validate are the go-playground/validator rules, e.g. "required", "gte=1".
	Validate []string
}

// Aliases is a block of type aliases: type Name = Target.
type Aliases struct {
	Aliases []Alias
}

// Alias is: type Name = Target.
type Alias struct {
	Name   string
	Target string
}

// QueryOpts is an operation query options type and its constructors.
type QueryOpts struct {
	// Type is the option type name, e.g. ListInstancesOpt.
	Type  string
	Funcs []QueryOpt
}

// QueryOpt is a query option constructor: func Func(Arg ArgType) Type.
type QueryOpt struct {
	Func    string
	Arg     string
	ArgType Type
	// Param is the query parameter name.
	Param string
}

// Findable is a Find method on a list response type,
// matching list items on one or two string fields.
type Findable struct {
	// ListType is the response type holding the list.
	ListType string
	// ListField is the ListType field holding the items.
	ListField string
	// ItemType is the list items type.
	ItemType string
	// Param is the method argument name.
	Param string
	// Fields are the item fields compared to Param (1 or 2).
	Fields []string
}

// Operation is a Client method calling an API operation.
type Operation struct {
	Doc         string
	Name        string
	OperationID string
	Method      string
	// Path is the URL path, with %v in place of path parameters when PathArgs is set.
	Path     string
	PathArgs []string
	// Params are the method parameters between ctx and req.
	Params []Param
	// Body is the request body type, nil if none.
	Body Type
	// QueryOpts is the query options type, "" if the operation has no query parameter.
	QueryOpts string
	// Return is the decoded response body type, nil if the method only returns an error.
	Return Type
	// ErrorResponses are the typed bodies of 4xx/5xx HTTP responses, sorted by code.
	ErrorResponses []ErrorResponse
	// SkipAuth is true for operations called without credentials.
	SkipAuth bool
}

// Param is a method parameter.
type Param struct {
	Name string
	Type Type
}

// ErrorResponse is the typed body of an error HTTP response.
type ErrorResponse struct {
	Code int
	Type Type
}

// Client is the API client with its zone endpoints.
type Client struct {
	Endpoints []Endpoint
}

// Endpoint is a zone endpoint constant.
type Endpoint struct {
	Name string
	URL  string
}

func (TypeDef) decl()   {}
func (Enum) decl()      {}
func (Struct) decl()    {}
func (Aliases) decl()   {}
func (QueryOpts) decl() {}
func (Findable) decl()  {}
func (Operation) decl() {}
func (Client) decl()    {}
