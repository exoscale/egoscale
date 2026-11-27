package problems

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"

	"github.com/pb33f/libopenapi"

	"github.com/exoscale/egoscale/v3/generator/helpers"
)

// problemType is an entry of the x-problems types of the spec
type problemType struct {
	Parent      string `yaml:"parent"`
	Status      int    `yaml:"status"`
	Title       string `yaml:"title"`
	Description string `yaml:"description"`
}

// Generate writes an error value per error type the services declare in the x-problems of the spec, which
// errors.Is matches against the errors of the client. The shared types are left out: the status errors, eg
// ErrNotFound, already match them
func Generate(doc libopenapi.Document, path, packageName string) error {
	r, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	var problems struct {
		Types map[string]problemType `yaml:"types"`
	}
	if r.Model.Extensions != nil {
		if node, ok := r.Model.Extensions.Get("x-problems"); ok {
			if err := node.Decode(&problems); err != nil {
				return fmt.Errorf("x-problems: %w", err)
			}
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	content, err := render(packageName, problems.Types)
	if err != nil {
		return err
	}
	output.Write(content)

	formatted, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}
	return os.WriteFile(path, formatted, os.ModePerm)
}

func render(packageName string, types map[string]problemType) ([]byte, error) {
	slugs := make([]string, 0, len(types))
	for slug, t := range types {
		if t.Parent != "" {
			slugs = append(slugs, slug)
		}
	}
	sort.Strings(slugs)

	out := &bytes.Buffer{}
	fmt.Fprintf(out, "package %s\n\n", packageName)
	fmt.Fprintf(out, "// The error types of the services: errors.Is(err, ErrKmsKeyNotFound) reports whether the API answered err\n")
	fmt.Fprintf(out, "// with this type. The status errors, eg ErrNotFound, match the shared type each of them refines\n")
	fmt.Fprintf(out, "var (\n")
	names := map[string]string{}
	for _, slug := range slugs {
		t := types[slug]
		name := "Err" + helpers.ToCamel(slug)
		if other, ok := names[name]; ok {
			return nil, fmt.Errorf("error types %q and %q have the same Go name %s", other, slug, name)
		}
		names[name] = slug
		doc := strings.TrimSpace(t.Description)
		if doc == "" {
			doc = t.Title
		}
		fmt.Fprintf(out, "\t// %s (%s, %d): %s\n", name, slug, t.Status, strings.ReplaceAll(doc, "\n", " "))
		fmt.Fprintf(out, "\t%s = &ProblemType{ID: %q, Status: %d, Title: %q}\n", name, "#/problems/"+slug, t.Status, t.Title)
	}
	fmt.Fprintf(out, ")\n")
	return out.Bytes(), nil
}
