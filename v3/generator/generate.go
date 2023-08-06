package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

func Generate() {
	tpl := template.Must(template.ParseFiles("templates/resource.tmpl"))

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "../oapi/oapi.gen.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for group, specs := range APIMap {
		for _, spec := range specs {
			for i, fn := range spec.Fns {
				args := FuncOptArgs(node, fn.OAPIName)
				if len(args) > 0 {
					fn.OptArgsDef = strings.Join(args, ", ")
					argKeys := make([]string, 0, len(args))
					for _, arg := range args {
						argKeys = append(argKeys, strings.Split(arg, " ")[0])
					}
					fn.OptArgsPassthrough = strings.Join(argKeys, ", ")
				}

				if fn.ResDefOverride != "" && fn.ResPassthroughOverride != "" {
					fn.ResDef = fn.ResDefOverride
					fn.ResPassthrough = fn.ResPassthroughOverride
				} else {
					resType, subpath := FuncResp(node, fn.OAPIName)
					fn.ResDef = "*" + resType
					fn.ResPassthrough = "resp.JSON200"
					if subpath != "" {
						fn.ResPassthrough = fmt.Sprintf("%s.%s", fn.ResPassthrough, subpath)
					}
					if strings.HasPrefix(fn.ResDef, "*[]") {
						fn.ResPassthrough = "*" + fn.ResPassthrough
						fn.ResDef = strings.TrimPrefix(fn.ResDef, "*")
					}
				}

				spec.Fns[i] = fn
			}

			f, err := os.Create(
				fmt.Sprintf(
					"../api/%s/%s.gen.go",
					group,
					strings.ToLower(spec.RootName),
				),
			)
			if err != nil {
				panic(err)
			}

			err = tpl.Execute(f, &spec)
			if err != nil {
				panic(err)
			}
		}
	}
}

// FuncOptArgs looks for function signature in oapi/oapi.gen.go and returns non-default arguments.
// Default arguments are context.Context and oapi.RequestEditorsFn.
func FuncOptArgs(node *ast.File, name string) []string {
	params := []string{}
	found := false

	// Search function arguments
	ast.Inspect(node, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			if fn.Name.Name == name+"WithResponse" {
				found = true
				// Skip first (context.Context) and last (RequestEditorsFn...) parameter as they are always the same.
				for i := 1; i < len(fn.Type.Params.List)-1; i++ {
					param := fn.Type.Params.List[i]

					switch t := param.Type.(type) {
					case *ast.SelectorExpr: //external type
						params = append(params, fmt.Sprintf(
							"%s %s.%s",
							param.Names[0].Name,
							t.X.(*ast.Ident).Name,
							t.Sel.Name,
						))
					case *ast.Ident: //internal or builtin type
						var format string
						if strings.HasSuffix(t.Name, "RequestBody") {
							// internal type, append package name
							format = "%s oapi.%s"
						} else {
							// builtin
							format = "%s %s"
						}
						params = append(params, fmt.Sprintf(
							format,
							param.Names[0].Name,
							t.Name,
						))
					default:
						// We are triggering specific panic which will print useful error message.
						_ = t.(*ast.SelectorExpr)
					}
				}

				return false
			}
		}

		return true
	})

	if !found {
		fmt.Printf("not found in oapi/oapi.gen.go: %s\n", name+"WithResponse")
		os.Exit(1)
	}

	return params
}

// FuncResp looks for a response body structure in oapi/oapi.gen.go.
// Returns type of JSON200 or nested struct attribute type and it's name.
// Pointers are implicit.
//
// For nested struct support is limited to single attribute structs with attribute type one of:
// - defined oapi struct type;
// - slice of defined oapi struct types;
//
// For example if we have the following types in oapi/oapi.gen.go:
//
//	type GetAccessKeyResponse struct {
//		Body         []byte
//		HTTPResponse *http.Response
//		JSON200      *AccessKey
//	}
//	type ListAccessKeysResponse struct {
//		Body         []byte
//		HTTPResponse *http.Response
//		JSON200      *struct {
//			AccessKeys *[]AccessKey `json:"access-keys,omitempty"`
//		}
//	}
//
// Calling FuncResp(node,"GetAccessKey") will return "AccessKey" and "",
// whle calling FuncResp(node,"ListAccessKeys") will return "[]AccessKey" and "AccessKeys".
func FuncResp(node *ast.File, name string) (resType, subpath string) {
	// Search response struct
	ast.Inspect(node, func(n ast.Node) bool {
		if spec, ok := n.(*ast.TypeSpec); ok {
			if str, ok := spec.Type.(*ast.StructType); ok {
				if spec.Name.Name != name+"Response" {
					return true
				}

				// Go to JSON200 attribute
				for _, f := range str.Fields.List {
					if f.Names[0].Name != "JSON200" {
						continue
					}

					// Supported JSON200 types:
					switch x := f.Type.(*ast.StarExpr).X.(type) {
					case *ast.Ident: //defined oapi type
						resType = fmt.Sprintf("oapi.%s", x.Name)
					case *ast.StructType: //nested struct, limited support for single attribute structs
						if len(x.Fields.List) != 1 {
							fmt.Printf("found %s response, but has unsupported nested struct")
							os.Exit(1)
						}

						y := x.Fields.List[0]
						subpath = y.Names[0].Name //path after JSON200

						switch z := y.Type.(*ast.StarExpr).X.(type) {
						case *ast.Ident: //defined oapi type
							resType = fmt.Sprintf("oapi.%s", z.Name)
						case *ast.ArrayType:
							resType = fmt.Sprintf("[]oapi.%s", z.Elt.(*ast.Ident).Name)
						default:
							// Panic which prints useful error message.
							_ = z.(*ast.Ident)
						}
					default:
						// Panic which prints useful error message.
						_ = x.(*ast.Ident)
					}
				}
			}
		}

		return true
	})

	return
}
