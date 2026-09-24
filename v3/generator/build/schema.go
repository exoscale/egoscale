package build

import (
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi/datamodel/high/base"
	"github.com/pb33f/libopenapi/orderedmap"
	"gopkg.in/yaml.v3"

	"github.com/exoscale/egoscale/v3/generator/ir"
	"github.com/exoscale/egoscale/v3/generator/naming"
)

// Schema kinds, from the schema type or inferred when it is missing (see kindOf).
const (
	kindBoolean = "boolean"
	kindInteger = "integer"
	kindNumber  = "number"
	kindString  = "string"
	kindArray   = "array"
	kindObject  = "object"
	// kindMap is an object with additionalProperties, always rendered as map[string]T.
	kindMap = "map"
)

var alphanumericRe = regexp.MustCompile("^[a-zA-Z0-9]+$")

// Schemas builds the file of the spec component schemas.
func (b *Builder) Schemas(packageName string) (ir.File, error) {
	if b.model.Components == nil || orderedmap.Len(b.model.Components.Schemas) == 0 {
		return ir.File{}, errors.New("no schema found in the spec")
	}

	for pair := orderedmap.SortAlpha(b.model.Components.Schemas).First(); pair != nil; pair = pair.Next() {
		schemaName, proxy := pair.Key(), pair.Value()
		if slices.Contains(b.cfg.IgnoredSchemas, schemaName) {
			continue
		}

		s, err := proxy.BuildSchema()
		if err != nil {
			return ir.File{}, fmt.Errorf("schema %s: %w", schemaName, err)
		}
		if err := b.schema(b.names.Camel(schemaName), s, "#/components/schemas/"+schemaName); err != nil {
			return ir.File{}, fmt.Errorf("schema %s: %w", schemaName, err)
		}
	}

	if len(b.cfg.Aliases) > 0 {
		aliases := ir.Aliases{}
		for _, a := range b.cfg.Aliases {
			aliases.Aliases = append(aliases.Aliases, ir.Alias{Name: a.Name, Target: a.Target})
		}
		if err := b.emit(aliases, "config"); err != nil {
			return ir.File{}, err
		}
	}

	return b.file(packageName, "net", "time"), nil
}

// schema emits the declaration of the type name from a schema,
// preceded by the declarations of its inline sub schemas.
// origin locates the schema in the spec for error messages.
func (b *Builder) schema(name string, s *base.Schema, origin string) error {
	doc := schemaDoc(s)

	switch kindOf(s) {
	case kindBoolean, kindInteger, kindNumber, kindString:
		if len(s.Enum) > 0 {
			enum, ok, err := b.enum(name, s)
			if err != nil {
				return err
			}
			if !ok {
				return fmt.Errorf("%s: enum values can't be rendered as Go constants", name)
			}
			enum.Doc = doc
			return b.emit(enum, origin)
		}

		typ, err := simpleType(s)
		if err != nil {
			return fmt.Errorf("%s: %w", name, err)
		}
		return b.emit(ir.TypeDef{Doc: doc, Name: name, Type: typ}, origin)
	case kindArray:
		typ, err := b.array(name, s, true, name, origin)
		if err != nil {
			return err
		}
		return b.emit(ir.TypeDef{Doc: doc, Name: name, Type: typ}, origin)
	case kindObject:
		st, err := b.object(name, s, origin)
		if err != nil {
			return err
		}
		st.Doc = doc
		return b.emit(st, origin)
	case kindMap:
		typ, err := b.simpleMap(name, s, name, origin)
		if err != nil {
			return err
		}
		return b.emit(ir.TypeDef{Doc: doc, Name: name, Type: typ}, origin)
	default:
		return fmt.Errorf("%s: schema type %v not implemented", name, s.Type)
	}
}

// enum returns the enum type name of a simple schema.
// ok is false if its values can't be rendered as Go constant names.
func (b *Builder) enum(name string, s *base.Schema) (_ ir.Enum, ok bool, _ error) {
	for _, e := range s.Enum {
		if strings.Contains(e.Value, ",") {
			return ir.Enum{}, false, nil
		}
		if len(e.Value) >= 1 && !alphanumericRe.MatchString(e.Value[:1]) {
			return ir.Enum{}, false, nil
		}
	}

	// x-go-enum-varnames provides explicit suffixes for each enum value
	// It must be used in the same order as the enum array.
	// Useful when: '1g.24gb-me' and '1g.24gb+me' returns the same suffix '1g24gbMe'.
	var varnames []string
	if s.Extensions != nil {
		if node, ok := s.Extensions.Get("x-go-enum-varnames"); ok {
			for _, n := range node.Content {
				varnames = append(varnames, n.Value)
			}
		}
	}
	useVarnames := len(varnames) == len(s.Enum)

	base, err := simpleType(s)
	if err != nil {
		return ir.Enum{}, false, fmt.Errorf("%s: %w", name, err)
	}

	enum := ir.Enum{Name: name, Base: base}
	for i, e := range s.Enum {
		literal := e.Value
		if base == ir.Named("string") {
			literal = fmt.Sprintf("%q", e.Value)
		}
		suffix := b.names.Camel(e.Value)
		if useVarnames {
			suffix = varnames[i]
		}
		enum.Values = append(enum.Values, ir.EnumValue{Name: name + suffix, Literal: literal})
	}

	return enum, true, nil
}

// array returns the slice type of an array schema. An inline object items schema
// is declared as typeName, suffixed with "Items" for a root schema.
// scope is the type name used to look up reference overrides.
func (b *Builder) array(typeName string, s *base.Schema, root bool, scope, origin string) (ir.Type, error) {
	if s.Items == nil {
		return nil, fmt.Errorf("array %s: items is nil", typeName)
	}
	if !s.Items.IsA() {
		return nil, fmt.Errorf("array %s: invalid spec version", typeName)
	}

	item, err := s.Items.A.BuildSchema()
	if err != nil {
		return nil, fmt.Errorf("array %s: build schema: %w", typeName, err)
	}
	if s.Items.A.IsReference() {
		return ir.Slice{Elem: b.ref(s.Items.A.GetReference(), scope)}, nil
	}

	if item.AdditionalProperties != nil {
		m, err := b.simpleMap(typeName, item, scope, origin)
		if err != nil {
			return nil, err
		}
		return ir.Slice{Elem: m}, nil
	}

	if isSimple(item) {
		typ, err := simpleType(item)
		if err != nil {
			return nil, fmt.Errorf("array %s: %w", typeName, err)
		}
		return ir.Slice{Elem: typ}, nil
	}

	if root {
		typeName += "Items"
	}
	if err := b.schema(typeName, item, origin); err != nil {
		return nil, err
	}

	return ir.Slice{Elem: ir.Named(typeName)}, nil
}

// object returns the struct of an object schema.
func (b *Builder) object(typeName string, s *base.Schema, origin string) (ir.Struct, error) {
	st := ir.Struct{Name: typeName}

	for pair := orderedmap.SortAlpha(s.Properties).First(); pair != nil; pair = pair.Next() {
		propName, proxy := pair.Key(), pair.Value()
		prop := proxy.Schema()
		if prop == nil {
			continue
		}

		propName = b.cfg.PropName(typeName, propName)

		nullable := prop.Nullable != nil && *prop.Nullable
		if slices.Contains(prop.Type, "null") {
			if len(nonNullTypes(prop.Type)) != 1 {
				return ir.Struct{}, fmt.Errorf("property %q: nullable type must contain exactly one non-null type", propName)
			}
			nullable = true
		}
		// https://github.com/pb33f/libopenapi/issues/283
		// Wait for a fix to remove this check func.
		if isNullableReference(proxy.GetReferenceNode()) {
			nullable = true
		}

		kind := kindOf(prop)
		required := slices.Contains(s.Required, propName)
		field := ir.Field{
			Doc:  schemaDoc(prop),
			Name: b.names.Camel(propName),
			JSON: propName,
			// Do not remove omitempty on simple nullable types.
			OmitEmpty: !required && (!nullable || isSimple(prop)),
			Validate:  validation(prop, required),
		}
		nestedName := typeName + field.Name

		switch {
		case proxy.IsReference():
			field.Type = b.ref(proxy.GetReference(), typeName)
			if kind != kindMap && (nullable || !isSimple(prop) || kind == kindBoolean) {
				field.Type = ir.Pointer{Elem: field.Type}
			}
		case kind == kindArray:
			typ, err := b.array(nestedName, prop, false, typeName, origin)
			if err != nil {
				return ir.Struct{}, err
			}
			field.Type = typ
		case isSimple(prop):
			if len(prop.Enum) > 0 {
				enum, ok, err := b.enum(nestedName, prop)
				if err != nil {
					return ir.Struct{}, err
				}
				if ok {
					if err := b.emit(enum, origin); err != nil {
						return ir.Struct{}, err
					}
					field.Type = ir.Named(nestedName)
					break
				}
			}

			typ, err := simpleType(prop)
			if err != nil {
				return ir.Struct{}, fmt.Errorf("property %q: %w", propName, err)
			}
			field.Type = typ
			// For simple types, we generate fields without pointers and rely on the JSON "omitempty" tag
			// to represent unset values in most cases. The exception is for boolean types, where we must
			// use a pointer to distinguish between "false" and "unset". If an unset value is required for
			// any other simple type, the OpenAPI spec must explicitly set "nullable: true" for that field.
			// Enum values are never pointers.
			if len(prop.Enum) == 0 && (nullable || kind == kindBoolean) {
				field.Type = ir.Pointer{Elem: typ}
			}
		case orderedmap.Len(prop.Properties) == 0 && prop.AdditionalProperties == nil:
			// This is an OpenAPI free form object (deprecated).
			// https://docs.42crunch.com/latest/content/oasv3/datavalidation/schema/v3-schema-object-without-properties.htm
			// We recommend to use AdditionalProperties instead.
			field.Type = ir.Map{Elem: ir.Any}
		case kind == kindMap:
			typ, err := b.simpleMap(nestedName, prop, typeName, origin)
			if err != nil {
				return ir.Struct{}, err
			}
			field.Type = typ
		default:
			if err := b.schema(nestedName, prop, origin); err != nil {
				return ir.Struct{}, err
			}
			field.Type = ir.Pointer{Elem: ir.Named(nestedName)}
		}

		st.Fields = append(st.Fields, field)
	}

	return st, nil
}

// simpleMap returns the map type of a schema with additionalProperties.
// An inline object values schema is declared as typeName.
func (b *Builder) simpleMap(typeName string, s *base.Schema, scope, origin string) (ir.Type, error) {
	// https://swagger.io/docs/specification/data-models/dictionaries/#free-form
	// There is two case for a free form object:
	//  - additionalProperties: true
	//  - additionalProperties: {}
	// Here is the libopenapi representation of it:

	//  - additionalProperties: true
	if s.AdditionalProperties.IsB() {
		return ir.Map{Elem: ir.Any}, nil
	}

	//  - additionalProperties: object
	if !s.AdditionalProperties.IsA() {
		return nil, fmt.Errorf("additional properties in: %s not supported", typeName)
	}

	proxy := s.AdditionalProperties.A
	if proxy.IsReference() {
		return ir.Map{Elem: b.ref(proxy.GetReference(), scope)}, nil
	}

	values := proxy.Schema()
	if values == nil {
		return nil, fmt.Errorf("additional properties in: %s: %w", typeName, proxy.GetBuildError())
	}
	//  - additionalProperties: {} empty object
	if len(values.Type) == 0 && orderedmap.Len(values.Extensions) == 0 {
		return ir.Map{Elem: ir.Any}, nil
	}
	if isSimple(values) {
		typ, err := simpleType(values)
		if err != nil {
			return nil, fmt.Errorf("additional properties in: %s: %w", typeName, err)
		}
		return ir.Map{Elem: typ}, nil
	}

	if err := b.schema(typeName, values, origin); err != nil {
		return nil, err
	}

	return ir.Map{Elem: ir.Named(typeName)}, nil
}

// ref returns the type of a schema reference, overridden in the config for the scope type name.
func (b *Builder) ref(reference, scope string) ir.Type {
	if typ, ok := b.cfg.RefType(scope, reference); ok {
		return ir.Named(typ)
	}

	return ir.Named(b.names.Camel(filepath.Base(reference)))
}

// validation returns the go-playground/validator rules of a field.
func validation(s *base.Schema, required bool) []string {
	var rules []string
	exclusiveMin := s.ExclusiveMinimum != nil && s.ExclusiveMinimum.IsA() && s.ExclusiveMinimum.A
	exclusiveMax := s.ExclusiveMaximum != nil && s.ExclusiveMaximum.IsA() && s.ExclusiveMaximum.A
	minOp, maxOp := "gte", "lte"
	if exclusiveMin {
		minOp = "gt"
	}
	if exclusiveMax {
		maxOp = "lt"
	}

	if s.MinLength != nil {
		rules = append(rules, fmt.Sprintf("%s=%v", minOp, *s.MinLength))
	}
	if s.MaxLength != nil {
		rules = append(rules, fmt.Sprintf("%s=%v", maxOp, *s.MaxLength))
	}
	if s.Minimum != nil {
		rules = append(rules, fmt.Sprintf("%s=%v", minOp, *s.Minimum))
	}
	if s.Maximum != nil {
		rules = append(rules, fmt.Sprintf("%s=%v", maxOp, *s.Maximum))
	}

	if required {
		return append([]string{"required"}, rules...)
	}
	// The omitempty rule is only needed in front of other rules.
	if len(rules) == 0 {
		return nil
	}

	return append([]string{"omitempty"}, rules...)
}

// simpleType returns the Go type of a simple schema (see isSimple).
func simpleType(s *base.Schema) (ir.Type, error) {
	if s.Extensions != nil {
		if typ, ok := s.Extensions.Get("x-go-type"); ok {
			return ir.Named(typ.Value), nil
		}
	}

	switch s.Format {
	case "":
	case "date-time":
		return ir.Named("time.Time"), nil
	case "uuid":
		return ir.Named("UUID"), nil
	case "ipv4", "ip":
		return ir.Named("net.IP"), nil
	case "uri-reference":
		return ir.Named("string"), nil
	case "byte":
		return ir.Named("[]byte"), nil
	case "double":
		return ir.Named("float64"), nil
	case "float":
		return ir.Named("float32"), nil
	case "int64", "int32":
		return ir.Named(s.Format), nil
	default:
		return nil, fmt.Errorf("schema format %q not implemented", s.Format)
	}

	types := nonNullTypes(s.Type)
	if len(types) == 0 {
		// Happens on oneOf/anyOf schemas, not supported yet.
		slog.Error("invalid spec: no type definition, Please fix the OpenApi Spec! Returning type any",
			slog.String("description", s.Description))
		return ir.Any, nil
	}

	switch types[0] {
	case kindBoolean:
		return ir.Named("bool"), nil
	case kindInteger:
		return ir.Named("int"), nil
	case kindNumber:
		return ir.Named("float64"), nil
	}

	return ir.Named(types[0]), nil
}

// kindOf returns the kind of a schema: its first non-null type,
// inferred from its content if missing, or kindMap with additionalProperties.
func kindOf(s *base.Schema) string {
	if s.AdditionalProperties != nil &&
		(s.AdditionalProperties.IsA() || (s.AdditionalProperties.IsB() && s.AdditionalProperties.B)) {
		return kindMap
	}

	if len(s.Type) == 0 {
		if orderedmap.Len(s.Properties) > 0 {
			return kindObject
		}
		if s.Items != nil {
			return kindArray
		}
		return ""
	}

	if types := nonNullTypes(s.Type); len(types) > 0 {
		return types[0]
	}

	return ""
}

// isSimple returns true if the schema is a scalar type, false otherwise.
func isSimple(s *base.Schema) bool {
	switch kindOf(s) {
	case kindObject, kindMap, kindArray:
		return false
	}

	return true
}

func nonNullTypes(types []string) []string {
	return slices.DeleteFunc(slices.Clone(types), func(t string) bool { return t == "null" })
}

func schemaDoc(s *base.Schema) string {
	doc := s.Description
	if doc == "" {
		doc = s.Title
	}

	return naming.Doc(doc)
}

// https://github.com/pb33f/libopenapi/issues/283
func isNullableReference(node *yaml.Node) bool {
	if node == nil || node.Content == nil {
		return false
	}

	for i, c := range node.Content {
		if c.Value == "nullable" {
			if i+1 < len(node.Content) && node.Content[i+1].Value == "true" {
				return true
			}
		}
	}

	return false
}
