// Package build turns an OpenAPI v3 model into ir declarations.
// All the OpenAPI spec interpretation of the generator lives here.
package build

import (
	"fmt"
	"slices"

	"github.com/pb33f/libopenapi"
	"github.com/pb33f/libopenapi/datamodel"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/ir"
	"github.com/exoscale/egoscale/v3/generator/naming"
)

// Builder builds the generated files declarations from an OpenAPI model.
// The same Builder must build all the files of a Go package,
// to detect type names declared twice.
type Builder struct {
	model *v3.Document
	cfg   config.Config
	names *naming.Namer

	// declared maps declared type names to where they come from in the spec.
	declared map[string]string
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
	}
}

// emit appends a declaration to the file being built.
// It fails if the declaration name is already declared.
func (b *Builder) emit(d ir.Decl, origin string) error {
	if name := declName(d); name != "" {
		if prev, ok := b.declared[name]; ok {
			return fmt.Errorf("type %s from %s is already declared from %s", name, origin, prev)
		}
		b.declared[name] = origin
	}
	b.decls = append(b.decls, d)

	return nil
}

func declName(d ir.Decl) string {
	switch d := d.(type) {
	case ir.TypeDef:
		return d.Name
	case ir.Enum:
		return d.Name
	case ir.Struct:
		return d.Name
	case ir.QueryOpts:
		return d.Type
	}

	return ""
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
