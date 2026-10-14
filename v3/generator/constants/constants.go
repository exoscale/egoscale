package constants

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"log/slog"
	"os"
	"strings"
	"text/template"

	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/exoscale/egoscale/v3/generator/helpers"
)

// Generate constants from OpenAPI spec servers into a go file.
func Generate(doc libopenapi.Document, path, packageName string) error {
	r, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	fmt.Fprintf(output, "package %s\n", packageName)

	if len(r.Model.Servers) == 0 {
		slog.Error("skip constants generation", slog.Int("servers_len", len(r.Model.Servers)))
		return nil
	}

	// The spec is returning only production server.
	if len(r.Model.Servers) != 1 {
		slog.Error("more than one server found", slog.Int("servers_len", len(r.Model.Servers)))
	}

	srv, err := renderConstants(r.Model.Servers[0])
	if err != nil {
		return err
	}
	output.Write(srv)

	if os.Getenv("GENERATOR_DEBUG") == "constants" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

// Template use constants.tmpl file
type Template struct {
	Enum          string
	ZoneEndpoints string
}

// renderConstants using the constants.tmpl template.
func renderConstants(s *v3.Server) ([]byte, error) {
	var tmpl Template

	if orderedmap.Len(s.Variables) == 0 {
		return nil, fmt.Errorf("no server variables defined")
	}

	for pair := s.Variables.First(); pair != nil; pair = pair.Next() {
		k, v := pair.Key(), pair.Value()

		if k != "zone" {
			// Supporting only zone variable for Exoscale
			continue
		}

		enum := make([]string, 0, len(v.Enum))
		endpoints := make([]string, 0, len(v.Enum))
		for _, z := range v.Enum {
			url := strings.Replace(s.URL, "{zone}", z, 1)
			enum = append(enum, fmt.Sprintf("%s Endpoint = %q", helpers.ToCamel(z), url))
			endpoints = append(endpoints, fmt.Sprintf("ZoneName(%q): %s,", z, helpers.ToCamel(z)))
		}

		tmpl = Template{
			Enum:          strings.Join(enum, "\n"),
			ZoneEndpoints: strings.Join(endpoints, "\n"),
		}
	}

	t, err := template.New("constants.tmpl").ParseFiles("./constants/constants.tmpl")
	if err != nil {
		return nil, err
	}

	output := bytes.NewBuffer([]byte{})
	if err := t.Execute(output, tmpl); err != nil {
		log.Fatal(err)
	}
	return output.Bytes(), nil
}
