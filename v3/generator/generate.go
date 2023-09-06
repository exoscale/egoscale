package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"

	"github.com/iancoleman/strcase"
)

// Generate is a command the generates Consumer API according to the mapping in the mapping.go file.
func Generate() {
	tpl := template.Must(template.ParseFiles("templates/resource.tmpl"))

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "../client.gen.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	for group, entities := range APIMap {
		for _, entity := range entities {
			entity.Package = group

			for i, fn := range entity.Fns {
				args := FuncOptArgs(node, fn.OAPIName)
				if len(args) > 0 {
					fn.OptArgsDef = strings.Join(args, ", ")
					argKeys := make([]string, 0, len(args))
					for _, arg := range args {
						argKeys = append(argKeys, strings.Split(arg, " ")[0])
					}
					fn.OptArgsPassthrough = strings.Join(argKeys, ", ")
					if strings.Contains(fn.OptArgsDef, "openapi_types") {
						entity.OAPITypesImport = true
					}
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

				entity.Fns[i] = fn
			}

			f, err := os.Create(
				fmt.Sprintf(
					"../%s_%s.gen.go",
					group,
					strcase.ToSnake(entity.RootName),
				),
			)
			if err != nil {
				panic(err)
			}

			err = tpl.Execute(f, &entity)
			if err != nil {
				panic(err)
			}
		}
	}
}

// FuncOptArgs looks for function signature in client.gen.go and returns non-default arguments (default is server url).
func FuncOptArgs(node *ast.File, name string) []string {
	params := []string{}
	found := false //to print useful error if function signature is not found.

	// Search function arguments
	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if fn.Name.Name != "new"+name+"Request" {
			return true
		}
		found = true
		// Skip first (server string) parameter.
		for i := 1; i < len(fn.Type.Params.List); i++ {
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
				params = append(params, fmt.Sprintf(
					"%s %s",
					param.Names[0].Name,
					t.Name,
				))
			default:
				// Panic to print useful error message.
				_ = t.(*ast.Ident)
			}
		}

		return true
	})

	if !found {
		fmt.Printf("not found in client.gen.go: %s\n", "new"+name+"Request")
		os.Exit(1)
	}

	return params
}

// FuncResp looks for a response body structure in client.gen.go.
// Returns type of JSON200 or nested struct attribute type and it's name.
// Pointers are implicit.
//
// For nested struct support is limited to single attribute structs with attribute type one of:
// - defined struct type;
// - slice of defined struct types;
//
// For example if we have the following types in client.gen.go:
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
		spec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		str, ok := spec.Type.(*ast.StructType)
		if !ok {
			return true
		}
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
			case *ast.Ident: //defined type
				resType = fmt.Sprintf("%s", x.Name)
			case *ast.ArrayType: //slice (usually List functions)
				resType = fmt.Sprintf("[]%s", x.Elt.(*ast.Ident).Name)
			case *ast.StructType: //nested struct, limited support for single attribute structs
				if len(x.Fields.List) != 1 {
					fmt.Printf("found %s response, but has unsupported nested struct\n", name)
					os.Exit(1)
				}

				y := x.Fields.List[0]
				subpath = y.Names[0].Name //path after JSON200

				switch z := y.Type.(*ast.StarExpr).X.(type) {
				case *ast.Ident: //defined type
					resType = fmt.Sprintf("%s", z.Name)
				case *ast.ArrayType: //slice (usually List functions)
					resType = fmt.Sprintf("[]%s", z.Elt.(*ast.Ident).Name)
				default:
					// Panic to print useful error message.
					_ = z.(*ast.Ident)
				}
			default:
				// Panic to print useful error message.
				_ = x.(*ast.Ident)
			}
		}

		return true
	})

	return
}
