package build

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

// problemType is an entry of the x-problems types of the spec.
type problemType struct {
	Parent      string `yaml:"parent"`
	Status      int    `yaml:"status"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// Problems builds the declarations of problems.go: an error value per error type the services declare in the
// x-problems of the spec. The shared types are left out: the status errors, e.g. ErrNotFound, already match them.
func (b *Builder) Problems(packageName string) (ir.File, error) {
	var problems struct {
		Types map[string]problemType `yaml:"types"`
	}
	if b.model.Extensions != nil {
		if node, ok := b.model.Extensions.Get("x-problems"); ok {
			if err := node.Decode(&problems); err != nil {
				return ir.File{}, fmt.Errorf("x-problems: %w", err)
			}
		}
	}

	decl := ir.ProblemTypes{}
	for _, slug := range slices.Sorted(maps.Keys(problems.Types)) {
		t := problems.Types[slug]
		if t.Parent == "" {
			continue
		}
		name := "Err" + b.names.Camel(slug)
		doc := strings.TrimSpace(t.Description)
		if doc == "" {
			doc = t.Title
		}
		decl.Types = append(decl.Types, ir.ProblemType{
			Doc:    fmt.Sprintf("// %s (%s, %d): %s", name, slug, t.Status, strings.ReplaceAll(doc, "\n", " ")),
			Name:   name,
			ID:     "#/problems/" + slug,
			Status: t.Status,
			Title:  t.Title,
		})
	}

	var decls []ir.Decl
	if err := b.emit(&decls, decl, "#/x-problems"); err != nil {
		return ir.File{}, err
	}

	return ir.File{Package: packageName, Decls: decls}, nil
}
