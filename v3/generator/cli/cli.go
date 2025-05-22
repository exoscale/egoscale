package cli

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/exoscale/egoscale/v3/generator/schemas"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
)

// cliCommand holds data for a single CLI command function.
type cliCommand struct {
	FuncName       string
	OperationID    string
	Params         []string
	HelpAndParse   string
	RequestBuilder string
	Args           string
}

// switchCase holds data for a single switch case in main().
type switchCase struct {
	OperationID string
	FuncName    string
	Summary     string
	Description string
}

// cliTemplateData is the root data passed to the CLI template.
type cliTemplateData struct {
	Commands    []cliCommand
	SwitchCases []switchCase
}

const mainTemplate = `
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	v3 "github.com/exoscale/egoscale/v3"
	"github.com/exoscale/egoscale/v3/credentials"
)

func printUsage() {
	commands := []struct {
		Name, Doc string
	}{
		{{- range .SwitchCases }}
		{"{{ .OperationID }}", {{ if .Summary }}{{ printf "%q" .Summary }}{{ else }}""{{ end }}},
		{{- end }}
	}
	maxLen := 0
	for _, c := range commands {
		if l := len(c.Name); l > maxLen {
			maxLen = l
		}
	}
	fmt.Println("Usage: " + os.Args[0] + " <command>")
	fmt.Println("Available commands:")
	for _, c := range commands {
		fmt.Printf("  %-*s %s\n", maxLen, c.Name, c.Doc)
	}
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Make credentials configurable via flags/env
	client, err := v3.NewClient(credentials.NewEnvCredentials())
	if err != nil {
		fmt.Println("failed to create client:", err)
		os.Exit(1)
	}

	switch os.Args[1] {
	{{- range .SwitchCases }}
	case "{{ .OperationID }}":
		{{ .FuncName }}Cmd(client)
	{{- end }}
	default:
		fmt.Println("unknown command:", os.Args[1])
		printUsage()
		os.Exit(1)
	}
}

{{ range .Commands }}
func {{ .FuncName }}Cmd(client *v3.Client) {
{{- if .Params }}
	flagset := flag.NewFlagSet("{{ .OperationID }}", flag.ExitOnError)
{{- range .Params }}
	{{ . }}
{{- end }}
{{ .HelpAndParse }}
{{ end -}}
{{ .RequestBuilder }}
	resp, err := client.{{ .FuncName }}({{ .Args }})
	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
	fmt.Printf("response: %+v\n", resp)
}
{{ end }}
`

func Generate(doc libopenapi.Document, path, packageName string) error {
	model, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	var (
		data        cliTemplateData
		switchCases = make([]switchCase, 0)
		commands    = make([]cliCommand, 0)
	)

	// Prepare switch cases and commands
	for pair := orderedmap.SortAlpha(model.Model.Paths.PathItems).First(); pair != nil; pair = pair.Next() {
		pathStr, pathItems := pair.Key(), pair.Value()
		for opPair := orderedmap.SortAlpha(pathItems.GetOperations()).First(); opPair != nil; opPair = opPair.Next() {
			_, operation := opPair.Key(), opPair.Value()

			funcName := helpers.ToCamel(operation.OperationId)
			if funcName == "" {
				funcName = helpers.ToCamel(pathStr)
			}

			// Add operation doc (summary/description) for usage
			summary := strings.TrimSpace(operation.Summary)
			description := strings.TrimSpace(operation.Description)
			switchCases = append(switchCases, switchCase{
				OperationID: operation.OperationId,
				FuncName:    funcName,
				Summary:     summary,
				Description: description,
			})

			params := make([]string, 0)
			args := []string{"context.Background()"}

			// Parameters (path/query)
			for _, param := range operation.Parameters {
				paramVar := helpers.ToLowerCamel(param.Name) + "Flag"
				paramType := schemas.RenderSimpleType(param.Schema.Schema())
				if paramType == "UUID" {
					paramType = "string"
				}
				params = append(params, fmt.Sprintf("var %s %s\n\tflagset.%sVar(&%s, \"%s\", %s, \"\")", paramVar, paramType, helpers.ToCamel(paramType), paramVar, helpers.ToCamel(param.Name), zeroValue(paramType)))
				if schemas.RenderSimpleType(param.Schema.Schema()) == "UUID" {
					paramVar = "v3.UUID(" + paramVar + ")"
				}
				args = append(args, paramVar)
			}

			// Request body: generate flags for each field in the request struct
			reqBodyType := funcName + "Request"
			hasRequestBody := false
			var reqBodyFields []reqBodyField
			if operation.RequestBody != nil {
				if media, ok := operation.RequestBody.Content.Get("application/json"); ok && media.Schema != nil {
					hasRequestBody = true
					schemaName := reqBodyType
					reqBodyFields = getRequestBodyFields(media.Schema, schemaName, "")
					for _, f := range reqBodyFields {
						flagVar := "req" + helpers.ToCamel(f.StructPath) + "Flag"
						flagType := f.Type
						if f.Type == "UUID" {
							flagType = "string"
						}
						flagDefault := zeroValue(flagType)
						params = append(params, fmt.Sprintf("var %s %s", flagVar, flagType))
						params = append(params, fmt.Sprintf("flagset.%sVar(&%s, \"%s\", %s, %q)", strings.Title(flagType), flagVar, f.Flag, flagDefault, f.Doc))
					}
				}
			}

			var requestBuilder strings.Builder
			if hasRequestBody {
				requestBuilder.WriteString("\t// Build request body struct from flags\n")
				requestBuilder.WriteString(fmt.Sprintf("\tvar req v3.%s\n", reqBodyType))
				requestBuilder.WriteString(buildNestedAssignments(reqBodyFields, "req"))
				args = append(args, "req")
			}

			helpAndParse := ""
			if len(params) > 0 {
				helpAndParse = `
	// Print help if no args or --help is present
	if len(os.Args) <= 2 || (len(os.Args) > 2 && (os.Args[2] == "-h" || os.Args[2] == "--help")) {
		flagset.Usage()
		os.Exit(0)
	}
	flagset.Parse(os.Args[2:])
`
			}

			commands = append(commands, cliCommand{
				FuncName:       funcName,
				OperationID:    operation.OperationId,
				Params:         params,
				HelpAndParse:   helpAndParse,
				RequestBuilder: requestBuilder.String(),
				Args:           strings.Join(args, ", "),
			})
		}
	}

	data.SwitchCases = switchCases
	data.Commands = commands

	tmpl := template.Must(template.New("cli").Parse(mainTemplate))

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return err
	}

	if os.Getenv("GENERATOR_DEBUG") == "cli" {
		fmt.Println(buf.String())
	}

	// Format and write the generated code
	content, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0644)
}

// zeroValue returns the zero value for a Go type as a string.
func zeroValue(goType string) string {
	switch goType {
	case "int", "int64", "float64":
		return "0"
	case "bool":
		return "false"
	default:
		return "\"\""
	}
}

// getRequestBodyFields extracts the fields for a request body struct from the schema.
// Returns a slice of {Name, Type, Flag, Doc}.
// Now supports nested objects with flag prefixing.
type reqBodyField struct {
	Type       string
	StructPath string
	Flag       string
	Doc        string
	ParentType string
}

func getRequestBodyFields(schemaProxy *base.SchemaProxy, schemaName string, parentObj string) []reqBodyField {
	fields := []reqBodyField{}
	s, err := schemaProxy.BuildSchema()
	if err != nil || s == nil {
		return fields
	}
	// Only handle object type
	if len(s.Type) == 0 || s.Type[0] != "object" {
		return fields
	}
	for pair := orderedmap.SortAlpha(s.Properties).First(); pair != nil; pair = pair.Next() {
		propName, propSchema := pair.Key(), pair.Value()
		prop := propSchema.Schema()
		if prop == nil {
			continue
		}

		flagName := propName
		if parentObj != "" {
			flagName = schemaName + "." + propName
		}

		structPath := helpers.ToCamel(propName)
		if parentObj != "" {
			structPath = helpers.ToCamel(schemaName) + "." + helpers.ToCamel(propName)
		}

		if !schemas.IsSimpleSchema(prop) {

			parentObj := helpers.ToCamel(schemaName) + helpers.ToCamel(propName)
			if propSchema.IsReference() {
				parentObj = helpers.RenderReference(propSchema.GetReference())
			}
			// Nested object: recurse, prefix flag and name
			nestedFields := getRequestBodyFields(propSchema, flagName, parentObj)
			// Do NOT append a flag for the object itself, only for its children
			fields = append(fields, nestedFields...)
			continue
		}
		goType := schemas.RenderSimpleType(prop)
		doc := ""
		if prop.Description != "" {
			doc = prop.Description
		} else if prop.Title != "" {
			doc = prop.Title
		}

		doc = strings.ReplaceAll(doc, "\n", " ")
		doc = strings.ReplaceAll(doc, "\t", "")
		// REMOVE trailing space between words

		fields = append(fields, reqBodyField{
			Type:       goType,
			StructPath: structPath,
			Flag:       flagName,
			Doc:        doc,
			ParentType: parentObj,
		})
	}
	return fields
}

// buildNestedAssignments generates Go code to assign flat flags to nested struct fields.
// It instantiates nested structs only if any of their fields are set.
func buildNestedAssignments(fields []reqBodyField, rootVar string) string {
	var assignments strings.Builder

	var lastParent string
	for i := len(fields) - 1; i >= 0; i-- {
		f := fields[i]

		if f.ParentType != "" {
			paths := strings.Split(f.StructPath, ".")
			path := strings.TrimRight(f.StructPath, "."+paths[len(paths)-1])

			varRef := "req" + helpers.ToCamel(f.StructPath) + "Flag"
			if f.Type == "UUID" {
				varRef = "v3.UUID(" + varRef + ")"
			}

			createStruct := ""
			if f.ParentType != lastParent {
				createStruct = fmt.Sprintf("req.%s = &v3.%s{}", path, f.ParentType)
			}

			assignments.WriteString(fmt.Sprintf(`if %s != %s {
				%s
				req.%s = %s
			}
			`, varRef, zeroValue(f.Type), createStruct, f.StructPath, varRef))

			lastParent = f.ParentType
			continue
		}

		assignments.WriteString(fmt.Sprintf("req.%s = req%sFlag\n", f.StructPath, f.StructPath))
	}

	return assignments.String()
}
