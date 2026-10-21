// Package build turns an OpenAPI v3 model into ir declarations.
// All the OpenAPI spec interpretation of the generator lives here.
package build

import (
	"errors"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/ir"
	"github.com/exoscale/egoscale/v3/generator/naming"
)

// Builder builds the generated files declarations from an OpenAPI model.
// The same Builder must build all the files of a Go package,
// to detect type names declared twice or never declared.
type Builder struct {
	model *v3.Document
	cfg   config.Config
	names *naming.Namer

	// declared maps declared type names to where they come from in the spec.
	declared map[string]string
	// used maps the type names used by declarations to where the first use comes from.
	used map[string]string
	// appliedRefOverrides are the config ref overrides used while building.
	appliedRefOverrides map[refOverride]bool
	// decls are the declarations of the file being built.
	decls []ir.Decl
}

// Load parses an OpenAPI spec into a v3 model. The spec is OpenAPI 3.0,
// where properties next to a $ref are ignored: they are neither turned into
// allOf schemas (OpenAPI 3.1 behavior) nor merged into the referenced schema.
func Load(spec []byte) (*v3.Document, error) {
	c := datamodel.NewDocumentConfiguration()
	c.TransformSiblingRefs = false
	c.MergeReferencedProperties = false

	doc, err := libopenapi.NewDocumentWithConfiguration(spec, c)
	if err != nil {
		return nil, fmt.Errorf("load spec: %w", err)
	}
	model, err := doc.BuildV3Model()
	if err != nil {
		return nil, fmt.Errorf("build model: %w", err)
	}

	return &model.Model, nil
}

// New returns a Builder of the model.
func New(model *v3.Document, cfg config.Config) *Builder {
	return &Builder{
		model:    model,
		cfg:      cfg,
		names:    naming.New(cfg.Acronyms),
		declared: map[string]string{},
		used:     map[string]string{},

		appliedRefOverrides: map[refOverride]bool{},
	}
}

// refOverride is a config ref override: a reference within a Go type.
type refOverride struct {
	typeName, ref string
}

// predeclared are the Go predeclared types used by the generated code.
var predeclared = []string{"any", "bool", "byte", "float32", "float64", "int", "int32", "int64", "string"}

// Check returns an error for every type used but never declared, which would
// not compile, and for every config ref override that didn't apply, to keep
// the config in sync with the spec. Call it once every file is built.
func (b *Builder) Check() error {
	var errs []error
	for _, name := range slices.Sorted(maps.Keys(b.used)) {
		_, declared := b.declared[name]
		if declared || strings.Contains(name, ".") ||
			slices.Contains(predeclared, name) || slices.Contains(b.cfg.ExternalTypes, name) {
			continue
		}
		errs = append(errs, fmt.Errorf("type %s used by %s is not declared", name, b.used[name]))
	}
	for _, typeName := range slices.Sorted(maps.Keys(b.cfg.RefOverrides)) {
		for _, ref := range slices.Sorted(maps.Keys(b.cfg.RefOverrides[typeName])) {
			if !b.appliedRefOverrides[refOverride{typeName, ref}] {
				errs = append(errs, fmt.Errorf("ref override %s %s: no such reference in the spec", typeName, ref))
			}
		}
	}

	return errors.Join(errs...)
}

// emit appends a declaration to the file being built.
// It fails if a type it declares is already declared.
func (b *Builder) emit(d ir.Decl, origin string) error {
	for _, name := range typeNames(d) {
		if prev, ok := b.declared[name]; ok {
			return fmt.Errorf("type %s from %s is already declared from %s", name, origin, prev)
		}
		b.declared[name] = origin
	}
	for _, name := range usedTypes(d) {
		if _, ok := b.used[name]; !ok {
			b.used[name] = origin
		}
	}
	b.decls = append(b.decls, d)

	return nil
}

// typeNames returns the names of the types a declaration declares.
func typeNames(d ir.Decl) []string {
	switch d := d.(type) {
	case ir.TypeDef:
		return []string{d.Name}
	case ir.Enum:
		return []string{d.Name}
	case ir.Struct:
		return []string{d.Name}
	case ir.QueryOpts:
		return []string{d.Type}
	case ir.Aliases:
		names := make([]string, 0, len(d.Aliases))
		for _, a := range d.Aliases {
			names = append(names, a.Name)
		}
		return names
	case ir.Constants:
		// Declared by the constants template.
		return []string{"Endpoint"}
	}

	return nil
}

// usedTypes returns the names of the types a declaration uses.
func usedTypes(d ir.Decl) []string {
	var types []ir.Type
	switch d := d.(type) {
	case ir.TypeDef:
		types = append(types, d.Type)
	case ir.Enum:
		types = append(types, d.Base)
	case ir.Struct:
		for _, f := range d.Fields {
			types = append(types, f.Type)
		}
	case ir.Aliases:
		for _, a := range d.Aliases {
			types = append(types, ir.Named(a.Target))
		}
	case ir.QueryOpts:
		for _, f := range d.Funcs {
			types = append(types, f.ArgType)
		}
	case ir.Findable:
		types = append(types, ir.Named(d.ListType), ir.Named(d.ItemType))
	case ir.Operation:
		for _, p := range d.Params {
			types = append(types, p.Type)
		}
		types = append(types, d.Body, d.Return)
		for _, r := range d.ErrorResponses {
			types = append(types, r.Type)
		}
	case ir.Constants:
		// Used by the constants template.
		types = append(types, ir.Named("ZoneName"))
	}

	var names []string
	for _, t := range types {
		names = append(names, namedTypes(t)...)
	}

	return names
}

// namedTypes returns the named types of a type expression.
func namedTypes(t ir.Type) []string {
	switch t := t.(type) {
	case ir.Named:
		return []string{string(t)}
	case ir.Pointer:
		return namedTypes(t.Elem)
	case ir.Slice:
		return namedTypes(t.Elem)
	case ir.Map:
		return namedTypes(t.Elem)
	}

	return nil
}

// file returns the file of the declarations emitted since the last call.
func (b *Builder) file(packageName string, imports ...string) ir.File {
	f := ir.File{Package: packageName, Imports: imports, Decls: b.decls}
	b.decls = nil

	return f
}

func (b *Builder) skipAuth(operationID string) bool {
	return slices.Contains(b.cfg.SkipAuthOperations, operationID)
}
