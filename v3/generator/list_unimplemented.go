package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

// ListUnimplemented goes though oapi client functions and prints thouse that are not yet mapped to Consumer API.
func ListUnimplemented() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "../oapi/oapi.gen.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	implemented := map[string]struct{}{}
	for _, entities := range APIMap {
		for _, entity := range entities {
			for _, fn := range entity.Fns {
				implemented[fn.OAPIName] = struct{}{}
			}
		}
	}

	ast.Inspect(node, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok {
			return true
		}
		if strings.HasSuffix(fn.Name.Name, "WithResponse") && !strings.HasSuffix(fn.Name.Name, "WithBodyWithResponse") {
			n := strings.TrimSuffix(fn.Name.Name, "WithResponse")
			if _, ok := implemented[n]; !ok {
				fmt.Println(n)
			}
		}
		return true
	})
}
