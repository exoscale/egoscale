package generator

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	wrapperTemplate = `package consumerapi

import (
	"context"
	"encoding/json"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/oapi"
)

// THIS CODE IS GENERATED

type {{.Platform}} struct {
	*oapi.ClientWithResponses
}
type {{.Asset}} struct {
	*oapi.ClientWithResponses
}

func ({{ .Asset | ToLower }} *{{.Asset}}) {{.Operation}}{{.Params}} *{{.JSON200Type}} {
	resp, err2 := {{ .Asset | ToLower }}.ClientWithResponses.{{.Operation}}{{.Asset}}WithResponse({{.ParamsList}})
	if err2 != nil {
		panic(err2)
	}

	return resp.JSON200
}

func ({{ .Platform | ToLower }} *{{.Platform}}) {{.Asset}}() *{{.Asset}} {
	return &{{.Asset}}{
		{{ .Platform | ToLower }}.ClientWithResponses,
	}
}

func (client *Client) {{.Platform}}() *{{.Platform}} {
	return &{{.Platform}}{
		client.ClientWithResponses,
	}
}
`
)

var platformNumber = 0

func generateWrapper(asset, operation, json200Type, params string, paramsList []string) {
	funcMap := template.FuncMap{
		"ToLower": strings.ToLower,
	}

	t, err := template.New(asset + operation + "Tmpl").Funcs(funcMap).Parse(wrapperTemplate)
	if err != nil {
		panic(err)
	}

	f, err3 := os.Create("../../gen/consumerapi/" + asset + operation + ".go")
	if err3 != nil {
		panic(err3)
	}

	paramsListStr := ""
	if len(paramsList) > 0 {
		paramsListStr = paramsList[0]
		for _, param := range paramsList[1:] {
			newParam := param
			if newParam == "reqEditors" {
				newParam = newParam + "..."
			}
			paramsListStr += ", " + newParam
		}
	}

	platform := "ExoPlatform" + strconv.Itoa(platformNumber)
	platformNumber++
	if asset == "Instance" {
		platform = "ComputeAPI"
	}
	err2 := t.Execute(f, map[string]string{
		"Platform":    platform,
		"Asset":       asset,
		"Operation":   operation,
		"JSON200Type": json200Type,
		"Params":      params,
		"ParamsList":  paramsListStr,
	})
	if err2 != nil {
		panic(err2)
	}
}

func wrapGetFuncs(fset *token.FileSet, pkg *ast.Package) {
	for _, pkgFile := range pkg.Files {
		for _, decl := range pkgFile.Decls {
			switch castedDecl := decl.(type) {
			case *ast.FuncDecl:
				if castedDecl.Type.Results == nil {
					continue
				}

				returnVals := castedDecl.Type.Results.List
				if len(returnVals) != 2 {
					continue
				}

				getPrefix := "Get"
				fname := castedDecl.Name.Name
				if len(fname) <= len(getPrefix) {
					continue
				}

				if !strings.HasPrefix(fname, getPrefix) {
					continue
				}

				withResponseSuffix := "WithResponse"
				if !strings.HasSuffix(fname, withResponseSuffix) {
					continue
				}

				asset := strings.TrimSuffix(fname, withResponseSuffix)
				asset = asset[len(getPrefix):]
				if strings.HasPrefix(asset, "Dbaas") {
					continue
				}

				if castedDecl.Type == nil || castedDecl.Type.Results == nil || len(castedDecl.Type.Results.List) < 1 {
					continue
				}
				respType := castedDecl.Type.Results.List[0].Type

				start := fset.Position(respType.Pos())
				end := fset.Position(respType.End())
				filename := start.Filename
				respTypeStr, err := readPart(filepath.Join(oapiPath, filename), int64(start.Offset), int64(end.Offset))
				if err != nil {
					panic(err)
				}

				respTypeStr = respTypeStr[1:]

				json200Type := findJSON200Type(fset, pkg, respTypeStr)
				if json200Type == "" {
					continue
				}

				pType := prepareType(json200Type)
				if pType == "" {
					continue
				}

				paramsList, params := getParams(fset, castedDecl)

				generateWrapper(asset, "Get", pType, params, paramsList)
			}
		}
	}
}

func getParams(fset *token.FileSet, fd *ast.FuncDecl) ([]string, string) {
	var paramList []string
	for _, param := range fd.Type.Params.List {
		for _, name := range param.Names {
			paramList = append(paramList, name.Name)
		}
	}

	content := getNode(fset, fd.Type.Params)
	return paramList, strings.ReplaceAll(content, "RequestEditorFn", "oapi.RequestEditorFn")
}

func getNode(fset *token.FileSet, node ast.Node) string {
	start := fset.Position(node.Pos())
	end := fset.Position(node.End())
	filename := start.Filename
	content, err := readPart(filepath.Join(oapiPath, filename), int64(start.Offset), int64(end.Offset))
	if err != nil {
		panic(err)
	}

	return content
}

func prepareType(typ string) string {
	if strings.HasPrefix(typ, "*struct") {
		return ""
		//return typ[1:]
	}

	if strings.HasPrefix(typ, "*[") {
		return ""
	}

	return "oapi." + typ[1:]
}

func findJSON200Type(fset *token.FileSet, pkg *ast.Package, typName string) string {
	for _, pkgFile := range pkg.Files {
		for _, decl := range pkgFile.Decls {
			switch castedDecl := decl.(type) {
			case *ast.GenDecl:
				for _, s := range castedDecl.Specs {
					switch sCasted := s.(type) {
					case *ast.TypeSpec:
						if sCasted.Name.Name == typName {
							switch castedStruct := sCasted.Type.(type) {
							case *ast.StructType:
								for _, f := range castedStruct.Fields.List {
									isJSON200 := false
									for _, i2 := range f.Names {
										if i2.Name == "JSON200" {
											isJSON200 = true
										}
									}

									if !isJSON200 {
										continue
									}

									tt := f.Type
									start := fset.Position(tt.Pos())
									end := fset.Position(tt.End())
									filename := start.Filename
									ff, err := readPart(filepath.Join(oapiPath, filename), int64(start.Offset), int64(end.Offset))
									if err != nil {
										panic(err)
									}

									return ff
								}
							}
						}
					}
				}
			}
		}
	}

	return ""
}

const oapiPath = "../../gen/oapi"

func Generate() error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, oapiPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	pkgName := "oapi"
	pkg, ok := pkgs[pkgName]
	if !ok {
		return fmt.Errorf("can't find package %q", pkgName)
	}

	wrapGetFuncs(fset, pkg)

	return nil
}

func readPart(filename string, start, end int64) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	size := end - start
	buffer := make([]byte, size)
	offset := int64(start)

	_, err = file.ReadAt(buffer, offset)
	if err != nil && err != io.EOF {
		return "", err
	}

	return string(buffer), nil
}
