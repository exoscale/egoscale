// A toolbox to check/generate CloudStack API commands from the JSON description

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"text/template"

	"github.com/exoscale/egoscale"
)

// must be sorted
var ignoredFields = []string{
	"customid",
	"deploymentplanner",
	"displaynetwork",
	"displayvm",
	"displayvolume",
	"forvpc",
	"haenable",
	"hostid",
	"hypervisor",
	"iscallerchilddomain",
	"isdynamicallyscalable",
	"podid",
	"policyid",
}

var source = flag.String("apis", "listApis.json", "listApis response in JSON")
var cmd = flag.String("cmd", "", "command name (e.g. listZones)")
var rtype = flag.String("type", "", "command return type name (e.g. Zone)")
var interfaces = flag.String("interface", "", "interface(s) to be filled")

var apiTypes = map[string]string{
	"short":   "int16",
	"integer": "int",
	"long":    "int64",
	"map":     "map[string]string",
	"list":    "[]struct{}",
	"set":     "[]struct{}",
	"uuid":    "*UUID",
	"boolean": "*bool",
	"date":    "string",
}

// fieldInfo represents the inner details of a field
type fieldInfo struct {
	Var       *types.Var
	OmitEmpty bool
	Doc       string
}

// response represents a response struc for a command
type response struct {
	name     string
	s        *types.Struct
	position token.Pos
	fields   map[string]fieldInfo
	errors   map[string][]error
}

// command represents a struct within the source code
type command struct {
	name        string
	description string
	sync        string
	s           *types.Struct
	position    token.Pos
	fields      map[string]fieldInfo
	errors      map[string][]error
	response    *response
}

func newCommand(obj types.Object) *command {
	return &command{
		name:     obj.Name(),
		s:        obj.Type().Underlying().(*types.Struct),
		position: obj.Pos(),
		fields:   map[string]fieldInfo{},
		errors:   map[string][]error{},
	}
}

func (c *command) setResponse(r *command) {
	c.response = &response{
		name:     r.name,
		s:        r.s,
		position: r.position,
		fields:   map[string]fieldInfo{},
		errors:   map[string][]error{},
	}
}

func (c *command) Check(api egoscale.API) {
	c.description = strings.Trim(api.Description, " ")

	if api.IsAsync {
		c.sync = " (A)"
	}

	c.CheckFields(api)
	c.CheckParams(api.Params)

	if c.response != nil && len(api.Response) > 0 {
		c.CheckResponse(api.Response)
	}
}

func (c *command) CheckResponse(response []egoscale.APIField) {
	errs := c.response.errors
	for i := 0; i < c.response.s.NumFields(); i++ {
		f := c.response.s.Field(i)
		if !f.IsField() || !f.Exported() {
			continue
		}

		tag := (reflect.StructTag)(c.response.s.Tag(i))
		var name string
		if match, ok := tag.Lookup("json"); !ok {
			n := f.Name()
			errs[n] = append(errs[n], errors.New("field error: no json annotation found"))
			continue
		} else {
			names := strings.Split(match, ",")
			name = names[0]
		}

		doc := tag.Get("doc")

		c.response.fields[name] = fieldInfo{
			Var: f,
			Doc: doc,
		}
	}

	for _, p := range response {
		n := p.Name
		index := sort.SearchStrings(ignoredFields, p.Name)
		ignored := index < len(ignoredFields) && ignoredFields[index] == p.Name
		if ignored {
			continue
		}

		field, ok := c.response.fields[p.Name]
		description := strings.Trim(p.Description, " ")

		if !ok {
			doc := ""
			if description != "" {
				doc = fmt.Sprintf(" doc:%q", description)
			}

			apiType, ok := apiTypes[p.Type]
			if !ok {
				apiType = p.Type
			}

			errs[n] = append(errs[n], fmt.Errorf("missing field:\n\t%s %s `json:\"%s,omitempty\"%s`", strings.Title(p.Name), apiType, p.Name, doc))
			continue
		}
		delete(c.fields, p.Name)

		typename := field.Var.Type().String()

		if field.Doc != description {
			if field.Doc == "" {
				errs[n] = append(errs[n], fmt.Errorf("missing doc:\n\t\t`doc:%q`", description))
			} else {
				errs[n] = append(errs[n], fmt.Errorf("wrong doc want %q got %q", description, field.Doc))
			}
		}

		expected := ""
		switch p.Type {
		case "short":
			if typename != "int16" {
				expected = "int16"
			}
		case "int":
		case "integer":
			// uint are used by port and icmp types
			if typename != "int" && typename != "uint16" && typename != "uint8" {
				expected = "int"
			}
		case "long":
			if typename != "int64" && typename != "uint64" {
				expected = "int64"
			}
		case "boolean":
			if typename != "bool" && typename != "*bool" {
				expected = "bool"
			}
		case "string":
		case "date":
		case "tzdate":
		case "imageformat":
			if typename != "string" {
				expected = "string"
			}
		case "uuid":
			if typename != "*egoscale.UUID" {
				expected = "*UUID"
			}
		case "list":
			if !strings.HasPrefix(typename, "[]") {
				expected = "[]string"
			}
		case "map":
		case "set":
			if !strings.HasPrefix(typename, "[]") {
				expected = "array"
			}
		case "state":
			// skip
		default:
			errs[n] = append(errs[n], fmt.Errorf("unknown type %q <=> %q", p.Type, field.Var.Type().String()))
		}

		if expected != "" {
			errs[n] = append(errs[n], fmt.Errorf("expected to be a %s, got %q", expected, typename))
		}
	}

	for name := range c.fields {
		errs[name] = append(errs[name], errors.New("extra field found"))
	}
}

func (c *command) CheckFields(api egoscale.API) {
	hasMeta := false

	for i := 0; i < c.s.NumFields(); i++ {
		f := c.s.Field(i)

		if !f.IsField() || !f.Exported() {
			if f.Name() != "_" {
				continue
			}

			tag := (reflect.StructTag)(c.s.Tag(i))
			name, nameOK := tag.Lookup("name")
			description, descriptionOK := tag.Lookup("description")
			if !nameOK || !descriptionOK {
				c.errors["_"] = append(c.errors["_"], fmt.Errorf("meta field incomplete, wanted\n\t_ bool `name:%q description:%q`", api.Name, c.description))
			} else {
				if name != api.Name || description != c.description {
					c.errors["_"] = append(c.errors["_"], fmt.Errorf("meta field incorrect, got name:%q description:%q, wanted\n\t_ bool `name:%q description:%q`", name, description, api.Name, c.description))
				}
			}

			hasMeta = true
			continue
		}

		name := ""
		var omitempty bool
		tag := (reflect.StructTag)(c.s.Tag(i))
		if match, ok := tag.Lookup("json"); !ok {
			n := f.Name()
			c.errors[n] = append(c.errors[n], errors.New("field error: no json annotation found"))
			continue
		} else {
			parts := strings.Split(match, ",")
			name = parts[0]
			omitempty = len(parts) > 1 && parts[1] == "omitempty"
		}

		doc := tag.Get("doc")

		c.fields[name] = fieldInfo{
			Var:       f,
			OmitEmpty: omitempty,
			Doc:       doc,
		}
	}

	if !hasMeta {
		c.errors["_"] = append(c.errors["_"], fmt.Errorf("meta field missing, wanted\n\t_ bool `name:%q description:%q`", api.Name, api.Description))
	}
}

func (c *command) CheckParams(params []egoscale.APIParam) {
	for _, p := range params {
		n := p.Name
		index := sort.SearchStrings(ignoredFields, p.Name)
		ignored := index < len(ignoredFields) && ignoredFields[index] == p.Name
		if ignored {
			continue
		}
		field, ok := c.fields[p.Name]
		description := strings.Trim(p.Description, " ")

		omit := ""
		if !p.Required {
			omit = ",omitempty"
		}

		if !ok {
			doc := ""
			if description != "" {
				doc = fmt.Sprintf(" doc:%q", description)
			}

			apiType, ok := apiTypes[p.Type]
			if !ok {
				apiType = p.Type
			}

			c.errors[n] = append(c.errors[n], fmt.Errorf("missing field:\n\t%s %s `json:\"%s%s\"%s`", strings.Title(p.Name), apiType, p.Name, omit, doc))
			continue
		}
		delete(c.fields, p.Name)

		typename := field.Var.Type().String()

		if field.Doc != description {
			if field.Doc == "" {
				c.errors[n] = append(c.errors[n], fmt.Errorf("missing doc:\n\t\t`doc:%q`", description))
			} else {
				c.errors[n] = append(c.errors[n], fmt.Errorf("wrong doc want %q got %q", description, field.Doc))
			}
		}

		if p.Required == field.OmitEmpty {
			c.errors[n] = append(c.errors[n], fmt.Errorf("wrong omitempty, want `json:\"%s%s\"`", p.Name, omit))
			continue
		}

		expected := ""
		switch p.Type {
		case "short":
			if typename != "int16" {
				expected = "int16"
			}
		case "int":
		case "integer":
			// uint are used by port and icmp types
			if typename != "int" && typename != "uint16" && typename != "uint8" {
				expected = "int"
			}
		case "long":
			if typename != "int64" && typename != "uint64" {
				expected = "int64"
			}
		case "boolean":
			if typename != "bool" && typename != "*bool" {
				expected = "bool"
			}
		case "string":
		case "date":
		case "tzdate":
		case "imageformat":
			if typename != "string" {
				expected = "string"
			}
		case "uuid":
			if typename != "*egoscale.UUID" {
				expected = "*UUID"
			}
		case "list":
			if !strings.HasPrefix(typename, "[]") {
				expected = "[]string"
			}
		case "map":
		case "set":
			if !strings.HasPrefix(typename, "[]") {
				expected = "array"
			}
		default:
			c.errors[n] = append(c.errors[n], fmt.Errorf("unknown type %q <=> %q", p.Type, field.Var.Type().String()))
		}

		if expected != "" {
			c.errors[n] = append(c.errors[n], fmt.Errorf("expected to be a %s, got %q", expected, typename))
		}
	}

	for name := range c.fields {
		c.errors[name] = append(c.errors[name], errors.New("extra field found"))
	}
}

func loadGoSources() (*types.Info, *token.FileSet) {
	files, err := filepath.Glob("*.go")
	if err != nil {
		panic(err)
	}
	fset := token.NewFileSet()
	astFiles := make([]*ast.File, len(files))
	for i, file := range files {
		f, er := parser.ParseFile(fset, file, nil, 0)
		if er != nil {
			panic(er)
		}
		astFiles[i] = f
	}

	config := types.Config{
		Importer: importer.For("source", nil),
	}

	info := &types.Info{
		Defs: make(map[*ast.Ident]types.Object),
	}

	_, err = config.Check("egoscale", fset, astFiles, info)
	if err != nil {
		_, e := fmt.Fprintf(os.Stderr, err.Error())
		if e != nil {
			panic(e)
		}
		os.Exit(1)
	}

	return info, fset
}

func checkSource(source, cmd, rtype string) {
	info, fset := loadGoSources()

	commands := make(map[string]*command)

	for _, obj := range info.Defs {
		if obj == nil || !obj.Exported() {
			continue
		}

		typ := obj.Type().Underlying()

		switch typ.(type) {
		case *types.Struct:
			c := newCommand(obj)
			commands[strings.ToLower(c.name)] = c
		}
	}

	sourceFile, _ := os.Open(source)

	decoder := json.NewDecoder(sourceFile)
	apis := new(egoscale.ListAPIsResponse)
	if err := decoder.Decode(&apis); err != nil {
		panic(err)
	}

	for _, a := range apis.API {
		name := strings.ToLower(a.Name)

		if command, ok := commands[name]; ok {
			if cmd == "" || strings.ToLower(cmd) == name {
				if rtype != "" {
					if resp, ok := commands[strings.ToLower(rtype)]; ok {
						command.setResponse(resp)
					}
				}
				command.Check(a)
				// too much information
				//} else {
				//fmt.Fprintf(os.Stderr, "command %q is missing\n", name)
			}
		}
	}

	names := make([]string, 0, len(commands))
	for name := range commands {
		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		c := commands[name]
		pos := fset.Position(c.position)
		er := len(c.errors)

		if cmd == "" {
			if er != 0 {
				fmt.Printf("%5d %s: %s%s\n", er, pos, c.name, c.sync)
			}
		} else if strings.ToLower(cmd) == name {
			errs := make([]string, 0, len(c.errors))
			for k, es := range c.errors {
				var b strings.Builder
				for i, e := range es {
					if i > 0 {
						if _, err := fmt.Fprintln(&b, ""); err != nil {
							panic(e)
						}
					}
					if _, err := fmt.Fprintf(&b, "%s: %s", k, e.Error()); err != nil {
						panic(err)
					}
				}
				errs = append(errs, b.String())
			}
			sort.Strings(errs)
			for _, e := range errs {
				fmt.Println(e)
			}
			fmt.Printf("\n%s: %s%s has %d error(s)\n", pos, c.name, c.sync, er)

			if c.response != nil {
				fmt.Println("")
				errs = make([]string, 0, len(c.response.errors))
				for k, es := range c.response.errors {
					for _, e := range es {
						errs = append(errs, fmt.Sprintf("%s: %s", k, e.Error()))
					}
				}

				sort.Strings(errs)
				for _, e := range errs {
					fmt.Println(e)
				}
				fmt.Printf("\n%s: %s has %d error(s)\n", fset.Position(c.response.position), c.response.name, len(errs))
			}

			os.Exit(er)
		}
	}

	if cmd != "" {
		fmt.Printf("%s not found\n", cmd)
		os.Exit(1)
	}
}

func generateInterface(interfaces, typeName string) {
	if !strings.HasPrefix(typeName, "List") || !strings.HasSuffix(typeName, "s") {
		fmt.Printf("Error: typeName must be of form List<type>s, got %q\n", typeName)
	}

	end := len(typeName) - 1
	if strings.HasSuffix(typeName, "ses") {
		end--
	}
	keyName := typeName[4:end]

	if interfaces == "Listable" {
		t := template.Must(template.New("listable").Parse(listableTemplate))

		fileName := fmt.Sprintf("%s_response.go", strings.ToLower(typeName[4:]))
		file, _ := os.Create(fileName)

		err := t.Execute(file, struct {
			Package string
			Type    string
			Key     string
		}{
			Package: "egoscale",
			Type:    typeName,
			Key:     keyName,
		})

		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}

		return
	}

	fmt.Printf("unknown interface: %q", interfaces)
}

func main() {
	flag.Parse()

	if *interfaces != "" {
		generateInterface(*interfaces, flag.Arg(0))
		return
	}

	if *source != "" {
		checkSource(*source, *cmd, *rtype)
	}
}

const listableTemplate = `// Code generated; DO NOT EDIT.

package {{.Package}}

import "fmt"

func ({{.Type}}) response() interface{} {
	return new({{.Type}}Response)
}

// SetPage sets the current apge
func (ls *{{.Type}}) SetPage(page int) {
	ls.Page = page
}

// SetPageSize sets the page size
func (ls *{{.Type}}) SetPageSize(pageSize int) {
	ls.PageSize = pageSize
}

// each triggers the callback for each, valid answer or any non 404 issue
func ({{.Type}}) each(resp interface{}, callback IterateItemFunc) {
	items, ok := resp.(*{{.Type}}Response)
	if !ok {
		callback(nil, fmt.Errorf("wrong type, {{.Type}}Response was expected, got %T", resp))
		return
	}

	for i := range items.{{.Key}} {
		if !callback(&items.{{.Key}}[i], nil) {
			break
		}
	}
}
`
