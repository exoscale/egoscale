package schemas

import (
	"bytes"
	"go/format"
	"log/slog"
	"os"
	"strings"

	"fmt"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/datamodel/low"

	"github.com/exoscale/egoscale/v3/generator/helpers"
)

// TODO fix the OpenApi spec (duplicated resources)
var ignoredList = map[string]struct{}{
	"snapshot-export": {},
}

// Generate go models from OpenAPI spec schemas into a go file.
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

	err := helpers.ForEachMapSorted(result.Model.Components.Schemas, func(schemaName string, v any) error {
		if _, ok := ignoredList[schemaName]; ok {
			return nil
		}

		r, err := RenderSchema(schemaName, v.(*base.SchemaProxy))
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

// RenderSchema returns generated go code from an OpenAPI Schema proxy object.
func RenderSchema(schemaName string, s *base.SchemaProxy) ([]byte, error) {
	output := bytes.NewBuffer([]byte{})

	sc, err := s.BuildSchema()
	if err != nil {
		return nil, err
	}

	schemaName = helpers.ToCamel(schemaName)
	if err := renderSchemaInternal(schemaName, sc, output); err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

// RenderSimpleType returns the type of a simple go type,
// not an object, map, array...etc.
// This function is called if you are sure IsSimpleSchema(s *base.Schema) return true.
// Add more simple type here.
func RenderSimpleType(s *base.Schema) string {
	typ, ok := s.Extensions["x-go-type"]
	if ok {
		t, ok := typ.(string)
		if ok {
			return t
		}

		slog.Error(
			"invalid x-go-type extension type, fallback on original type",
			slog.Any("type", typ),
		)
	}

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
		slog.Error(
			"invalid spec: no type definition, Please fix the OpenApi Spec! Returning type any",
			slog.Any("schema", s),
		)
		return "any"
	}

	if s.Type[0] == "boolean" {
		return "bool"
	}
	if s.Type[0] == "integer" {
		return "int"
	}

	return s.Type[0]
}

// renderSchemaInternal render a given libopenapi Schema into a buffer.
// This function is mostly used recursively to render sub schemas object into this buffer.
//
// /!\ for every recursive call, make sure to check schema reference before,
// to prevent end up in infinite loop.
// That prevent embed hash cookies in the whole codebase to compare schemas:
// e.g: https://github.com/danielgtaylor/restish/blob/main/openapi/schema.go#L59
func renderSchemaInternal(schemaName string, s *base.Schema, output *bytes.Buffer) error {
	doc := renderDoc(s) + "\n"
	InferType(s)

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
			definition = renderSimpleTypeEnum(schemaName, s)
		} else {
			definition = "type " + schemaName + " " + RenderSimpleType(s) + "\n"
		}

		output.WriteString(definition)
		return nil
	case "array":
		output.WriteString(doc)
		array, err := renderArray(schemaName, s, output)
		if err != nil {
			return err
		}
		output.WriteString("type " + schemaName + " " + array + "\n")
		return nil
	case "object":
		object, err := renderObject(schemaName, s, output)
		if err != nil {
			return err
		}
		output.WriteString(doc)
		output.WriteString(object)
		return nil
	// map represents an OpenAPI AdditionalProperties, it will always be map[string]T
	case "map":
		output.WriteString(doc)
		Map, err := renderSimpleMap(schemaName, s, output)
		if err != nil {
			return err
		}
		typeString := "type " + schemaName + " "
		if s.Nullable != nil && *s.Nullable {
			typeString += "*"
		}
		typeString += Map + "\n"
		output.WriteString(typeString)
		return nil
	default:
		slog.Error("type not implemented", slog.String("type", typ))
		return nil
	}
}

func renderSimpleTypeEnum(typeName string, s *base.Schema) string {
	typ := RenderSimpleType(s)
	definition := "type " + typeName + " " + typ + "\n"
	definition += "const (\n"

	for _, e := range s.Enum {
		value := helpers.ToCamel(fmt.Sprint(e))
		if typ == "string" {
			value = fmt.Sprintf(`"%v"`, e)
		}
		definition += typeName + helpers.ToCamel(fmt.Sprintf("%v", e)) + " " + typeName + " = " + value + "\n"
	}
	definition += ")\n"

	return definition
}

func renderArray(typeName string, s *base.Schema, output *bytes.Buffer) (string, error) {
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
		Map, err := renderSimpleMap(typeName, item, output)
		if err != nil {
			return "", err
		}
		return definition + Map, nil
	}

	simple := IsSimpleSchema(item)
	if simple {
		return definition + RenderSimpleType(item), nil
	}

	// Render new object from array schema into the buffer.
	if err := renderSchemaInternal(typeName, item, output); err != nil {
		return "", err
	}

	return definition + typeName, nil
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
	// Remove the useless omitempty validation if it's the only one.
	// JSON omit empty will be already there.
	if len(ops) == 1 && ops[0] == "omitempty" {
		return ""
	}

	return fmt.Sprintf(`validate:"%s"`, strings.Join(ops, ","))
}

func renderObject(typeName string, s *base.Schema, output *bytes.Buffer) (string, error) {
	definition := "type " + typeName + " struct {\n"
	err := helpers.ForEachMapSorted(s.Properties, func(propName string, v any) error {
		properties := v.(*base.SchemaProxy)

		prop := properties.Schema()
		if prop == nil {
			return nil
		}

		var nullable = false
		if prop.Nullable != nil {
			nullable = *prop.Nullable
		}

		doc := renderDoc(prop)
		if doc != "" {
			definition += doc + "\n"
		}

		tag := fmt.Sprintf(" `json:\"%s,omitempty\"", propName)

		pointer := "*"
		req := isRequiredField(propName, s)
		if req {
			tag = fmt.Sprintf(" `json:%q", propName)
		}
		validation := renderValidation(prop, req)
		if validation != "" {
			tag += " " + validation
		}
		tag += "`"

		camelName := helpers.ToCamel(propName)

		if properties.IsReference() {
			referenceName := helpers.RenderReference(properties.GetReference())
			if prop.AdditionalProperties != nil {
				pointer = ""
			}
			if !nullable && IsSimpleSchema(prop) {
				pointer = ""
			}
			definition += camelName + " " + pointer + referenceName + tag + "\n"
			return nil
		}

		if prop.Type[0] == "array" {
			array, err := renderArray(typeName+camelName, prop, output)
			if err != nil {
				return err
			}
			definition += camelName + " " + array + tag + "\n"
			return nil
		}

		if IsSimpleSchema(prop) {
			// Render property type enum.
			if len(prop.Enum) > 0 {
				output.WriteString(renderSimpleTypeEnum(typeName+camelName, prop))
				definition += camelName + " " + typeName + camelName + tag + "\n"
				return nil
			}

			// To be discuss here, for the moment we bypass pointer on those types,
			// and use JSON omitempty. This will cover most of all case.
			// For the specific types left like in PUT requests schema,
			// we need to update the spec to put those type as nullable, take the instance-pool as good example,
			// (only use custom schema, not schema reference for PUT request).
			if !nullable && (prop.Type[0] == "string" || prop.Type[0] == "integer") {
				definition += camelName + " " + RenderSimpleType(prop) + tag + "\n"
				return nil
			}
			definition += camelName + " " + pointer + RenderSimpleType(prop) + tag + "\n"

			return nil
		}

		// This is an OpenAPI free form object (deprecated).
		// https://docs.42crunch.com/latest/content/oasv3/datavalidation/schema/v3-schema-object-without-properties.htm
		// We recommend to use AdditionalProperties instead.
		if len(prop.Properties) == 0 && prop.AdditionalProperties == nil {
			definition += camelName + " map[string]any" + tag + "\n"
			return nil
		}

		// Render additional properties (map).
		if prop.AdditionalProperties != nil {
			Map, err := renderSimpleMap(typeName+camelName, prop, output)
			if err != nil {
				return err
			}
			definition += camelName + " " + Map + tag + "\n"
			return nil
		}

		// Render new object from object property into the buffer.
		if err := renderSchemaInternal(typeName+camelName, prop, output); err != nil {
			return err
		}
		definition += camelName + " " + pointer + typeName + camelName + tag + "\n"

		return nil
	})
	if err != nil {
		return "", err
	}

	return definition + "}\n\n", nil
}

func isRequiredField(schemaName string, s *base.Schema) bool {
	for _, req := range s.Required {
		if req == schemaName {
			return true
		}
	}

	return false
}

// renderSimpleMap represents AdditionalProperties, it's always a map[string]Type
func renderSimpleMap(typeName string, s *base.Schema, output *bytes.Buffer) (string, error) {
	definition := "map[string]"

	switch ap := s.AdditionalProperties.(type) {
	case *base.SchemaProxy:
		break
	// https://swagger.io/docs/specification/data-models/dictionaries/#free-form
	// There is two case for a free form object:
	//  - additionalProperties: true
	//  - additionalProperties: {}
	// Here is the libopenapi representation of it.
	case bool, map[low.KeyReference[string]]low.ValueReference[interface{}]:
		return definition + "any", nil
	default:
		return "", fmt.Errorf("additional properties in: %s not supported: %#v", typeName, ap)
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

	// Render new object from AdditionalProperties schema into the buffer.
	if err := renderSchemaInternal(typeName, addl, output); err != nil {
		return "", err
	}

	return definition + typeName, nil
}

// InferType fixes missing type if it is missing & can be inferred
func InferType(s *base.Schema) {
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

	return helpers.RenderDoc(doc)
}
