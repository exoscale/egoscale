package schemas

import (
	"bytes"
	"go/format"
	"os"
	"strings"

	"fmt"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/low"
)

// TODO fix the OpenApi spec (duplicated resources)
var ignoredList = map[string]struct{}{
	"snapshot-export": {},
}

func Generate(doc libopenapi.Document, path, packageName string) error {
	result, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	output.WriteString(fmt.Sprintf(`package %s

import (
	"net"
	"time"
)
`, packageName))

	err := helpers.ForEachMapSorted(result.Model.Components.Schemas, func(k string, v any) error {
		if _, ok := ignoredList[k]; ok {
			return nil
		}

		r, err := RenderSchema(k, v.(*base.SchemaProxy))
		if err != nil {
			return err
		}
		output.Write(r)
		output.WriteString("\n")
		return nil
	})
	if err != nil {
		return err
	}

	if os.Getenv("GENERATOR_DEBUG") == "schemas" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

// RenderSchema returns generated GO code from an OpenAPI Schema.
func RenderSchema(name string, s *base.SchemaProxy) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	sc, err := s.BuildSchema()
	if err != nil {
		return nil, err
	}

	name = helpers.ToCamel(name)
	if err := renderSchemaInternal(name, sc, map[[32]byte]bool{}, output); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// RenderSimpleType returns the type of a simple go type,
// not an object, map, array...etc.
// Add more simple type here.
func RenderSimpleType(s *base.Schema) string {
	if s.Format != "" {
		if s.Format == "date-time" {
			return "time.Time"
		}
		if s.Format == "uuid" {
			return "UUID"
		}
		if s.Format == "ipv4" {
			return "net.IP"
		}

		return s.Format
	}

	if len(s.Type) == 0 {
		// This should never happen.
		panic("invalid spec, no type definition")
	}
	if s.Type[0] == "boolean" {
		return "bool"
	}
	if s.Type[0] == "integer" {
		return "int"
	}

	return s.Type[0]
}

func renderSchemaInternal(n string, s *base.Schema, known map[[32]byte]bool, output *bytes.Buffer) error {
	doc := renderDoc(s) + "\n"
	inferType(s)

	// TODO: list type alternatives somehow?
	// In OpenAPI versions 2 and 3.0, this Type is a single value,
	// so array will only ever have one value in version 3.1,
	// Type can be multiple values
	typ := ""
	for _, t := range s.Type {
		// Find the first non-null type and use that for now.
		if t != "null" {
			typ = t
			break
		}
	}

	switch typ {
	case "boolean", "integer", "number", "string":
		output.WriteString(doc)
		var definition string

		if len(s.Enum) > 0 {
			definition = renderSimpleTypeEnum(n, s)
		} else {
			definition = "type " + n + " " + RenderSimpleType(s) + "\n"
		}

		output.WriteString(definition)
		return nil
	case "array":
		output.WriteString(doc)
		array, err := renderArray(n, s, known, output)
		if err != nil {
			return err
		}
		output.WriteString("type " + n + " " + array + "\n")
		return nil
	case "object":
		object, err := renderObject(n, s, known, output)
		if err != nil {
			return err
		}
		output.WriteString(doc)
		output.WriteString(object)
		return nil
	// map represents an AdditionalProperties, it will always be map[string]Type
	case "map":
		output.WriteString(doc)
		Map, err := renderSimpleMap(n, s, known, output)
		if err != nil {
			return err
		}
		output.WriteString("type " + n + " " + Map + "\n")
		return nil
	default:
		panic("type: " + typ + " not implemented")
	}
}

func renderSimpleTypeEnum(n string, s *base.Schema) string {
	typ := RenderSimpleType(s)
	definition := "type " + n + " " + typ + "\n"
	definition += "const (\n"

	for _, e := range s.Enum {
		value := helpers.ToCamel(fmt.Sprint(e))
		if typ == "string" {
			value = fmt.Sprintf(`"%v"`, e)
		}
		definition += n + helpers.ToCamel(fmt.Sprintf("%v", e)) + " " + n + " = " + value + "\n"
	}
	definition += ")\n"

	return definition
}

func renderArray(n string, s *base.Schema, known map[[32]byte]bool, output *bytes.Buffer) (string, error) {
	definition := "[]"

	if s.Items == nil {
		return "", fmt.Errorf("array: items is nil")
	}
	if !s.Items.IsA() {
		return "", fmt.Errorf("array: invalid spec version")
	}

	item, err := s.Items.A.BuildSchema()
	if err != nil {
		return "", fmt.Errorf("array: build schema: %w", err)
	}
	isReference := s.Items.A.IsReference()
	if isReference {
		return definition + helpers.RenderReference(s.Items.A.GetReference()), nil
	}

	if item.AdditionalProperties != nil {
		Map, err := renderSimpleMap(n, item, known, output)
		if err != nil {
			return "", err
		}
		return definition + Map, nil
	}

	simple := IsSimpleSchema(item)
	if simple {
		return definition + RenderSimpleType(item), nil
	}

	hash := item.GoLow().Hash()
	if !known[hash] {
		known[hash] = true
		if err := renderSchemaInternal(n, item, known, output); err != nil {
			return "", err
		}
		known[hash] = false
		return definition + n, nil
	}

	panic("array: not reachable")
}

func renderValidation(s *base.Schema, required bool) string {
	ops := []string{}

	if required {
		ops = append(ops, "required")
	} else {
		ops = append(ops, "omitempty")
	}

	if s.MinLength != nil {
		op := "gte"
		if s.ExclusiveMinimum != nil && s.ExclusiveMinimum.IsA() && s.ExclusiveMinimum.A {
			op = "gt"
		}
		ops = append(ops, fmt.Sprintf("%s=%v", op, *s.MinLength))
	}
	if s.MaxLength != nil {
		op := "lte"
		if s.ExclusiveMaximum != nil && s.ExclusiveMaximum.IsA() && s.ExclusiveMaximum.A {
			op = "lt"
		}
		ops = append(ops, fmt.Sprintf("%s=%v", op, *s.MaxLength))
	}
	if s.Minimum != nil {
		op := "gte"
		if s.ExclusiveMinimum != nil && s.ExclusiveMinimum.IsA() && s.ExclusiveMinimum.A {
			op = "gt"
		}
		ops = append(ops, fmt.Sprintf("%s=%v", op, *s.Minimum))
	}
	if s.Maximum != nil {
		op := "lte"
		if s.ExclusiveMaximum != nil && s.ExclusiveMaximum.IsA() && s.ExclusiveMaximum.A {
			op = "lt"
		}
		ops = append(ops, fmt.Sprintf("%s=%v", op, *s.Maximum))
	}

	if len(ops) == 0 {
		return ""
	}
	if len(ops) == 1 && ops[0] == "omitempty" {
		return ""
	}

	validation := `validate:"`
	for i, op := range ops {
		if i == len(ops)-1 {
			validation += op
			break
		}

		validation += op + ","

	}

	return validation + `"`
}

func renderObject(n string, s *base.Schema, known map[[32]byte]bool, output *bytes.Buffer) (string, error) {
	definition := "type " + n + " struct {\n"
	err := helpers.ForEachMapSorted(s.Properties, func(name string, v any) error {
		properties := v.(*base.SchemaProxy)

		prop := properties.Schema()
		if prop == nil {
			return nil
		}

		doc := renderDoc(prop)
		if doc != "" {
			definition += doc + "\n"
		}

		tag := fmt.Sprintf(" `json:\"%s,omitempty\"", name)

		required := "*"
		req := isRequiredField(name, s)
		if req {
			tag = fmt.Sprintf(" `json:%q", name)
		}
		validation := renderValidation(prop, req)
		if validation != "" {
			tag += " " + validation
		}
		tag += "`"

		camelName := helpers.ToCamel(name)

		if properties.IsReference() {
			referenceName := helpers.RenderReference(properties.GetReference())
			if prop.AdditionalProperties != nil {
				required = ""
			}
			definition += camelName + " " + required + referenceName + tag + "\n"
			return nil
		}

		if prop.Type[0] == "array" {
			array, err := renderArray(n+camelName, prop, known, output)
			if err != nil {
				return err
			}
			definition += camelName + " " + array + tag + "\n"
			return nil
		}

		if IsSimpleSchema(prop) {
			if len(prop.Enum) > 0 {
				output.WriteString(renderSimpleTypeEnum(n+camelName, prop))
				definition += camelName + " " + n + camelName + tag + "\n"
				return nil
			}

			if prop.Type[0] == "string" || prop.Type[0] == "integer" {
				definition += camelName + " " + RenderSimpleType(prop) + tag + "\n"
				return nil
			}

			definition += camelName + " " + required + RenderSimpleType(prop) + tag + "\n"

			return nil
		}

		if len(prop.Properties) == 0 && prop.AdditionalProperties == nil {
			definition += camelName + " map[string]any" + tag + "\n"
			return nil
		}

		if prop.AdditionalProperties != nil {
			Map, err := renderSimpleMap(n+camelName, prop, known, output)
			if err != nil {
				return err
			}
			definition += camelName + " " + Map + tag + "\n"
			return nil
		}

		hash := prop.GoLow().Hash()
		if !known[hash] {
			known[hash] = true
			if err := renderSchemaInternal(n+camelName, prop, known, output); err != nil {
				return err
			}
			known[hash] = false
			definition += camelName + " " + required + n + camelName + tag + "\n"
		}

		return nil
	})
	if err != nil {
		return "", err
	}

	return definition + "}\n\n", nil
}

func isRequiredField(name string, s *base.Schema) bool {
	for _, req := range s.Required {
		if req == name {
			return true
		}
	}

	return false
}

// renderSimpleMap represents AdditionalProperties, it's always a map[string]Type
func renderSimpleMap(n string, s *base.Schema, known map[[32]byte]bool, output *bytes.Buffer) (string, error) {
	definition := "map[string]"

	switch ap := s.AdditionalProperties.(type) {
	case *base.SchemaProxy:
		break
	case bool, map[low.KeyReference[string]]low.ValueReference[interface{}]:
		return definition + "any", nil
	default:
		return "", fmt.Errorf("additional properties in: %s not supported: %#v", n, ap)
	}

	sp := s.AdditionalProperties.(*base.SchemaProxy)
	if sp.IsReference() {
		return definition + helpers.RenderReference(sp.GetReference()), nil
	}

	addl := sp.Schema()
	simple := IsSimpleSchema(addl)
	if simple {
		return definition + RenderSimpleType(addl), nil
	}

	hash := addl.GoLow().Hash()
	if !known[hash] {
		known[hash] = true
		if err := renderSchemaInternal(n, addl, known, output); err != nil {
			return "", err
		}
		known[hash] = false
		return definition + n, nil
	}

	panic("not reachable")
}

// inferType fixes missing type if it is missing & can be inferred
func inferType(s *base.Schema) {
	if len(s.Type) == 0 {
		if s.Items != nil {
			s.Type = []string{"array"}
		}

		if len(s.Properties) > 0 {
			s.Type = []string{"object"}
		}
	}

	// AdditionalProperties will always be map[string]Type
	if s.AdditionalProperties != nil {
		s.Type = []string{"map"}
	}
}

// IsSimpleSchema returns whether this schema is a scalar or array as these
// can't be circular references. Objects result in `false` and that triggers
// circular ref checks.
func IsSimpleSchema(s *base.Schema) bool {
	if len(s.Type) == 0 {
		return true
	}

	return s.Type[0] != "object" && s.Type[0] != "map"
}

func renderDoc(s *base.Schema) string {
	doc := s.Description
	if doc == "" {
		doc = s.Title
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
