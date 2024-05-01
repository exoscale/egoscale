package client

import (
	"bytes"
	"fmt"
	"go/format"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
	"github.com/pb33f/libopenapi/orderedmap"
)

// Generate go client from OpenAPI spec servers into a go file.
func Generate(doc libopenapi.Document, path, packageName string) error {
	r, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	output.WriteString(fmt.Sprintf(`package %s
	import (
		"fmt"
		"net/http"
		"context"
		"runtime"
		"time"

		"github.com/exoscale/egoscale/v3/credentials"
		"github.com/exoscale/egoscale/version"
		"github.com/go-playground/validator/v10"
	)
	`, packageName))

	for _, s := range r.Model.Servers {
		if !strings.Contains(s.URL, "api") {
			// Skip generating code for pre-production "ppapi" server.
			continue
		}

		srv, err := renderClient(s)
		if err != nil {
			return err
		}
		output.Write(srv)
	}

	if os.Getenv("GENERATOR_DEBUG") == "client" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

// Template use client.tmpl file
type Template struct {
	Enum           string
	ServerEndpoint string
}

// renderClient using the client.tmpl template.
func renderClient(s *v3.Server) ([]byte, error) {
	var client Template

	if orderedmap.Len(s.Variables) == 0 {
		return nil, fmt.Errorf("no server variables defined")
	}

	for pair := s.Variables.First(); pair != nil; pair = pair.Next() {
		k, v := pair.Key(), pair.Value()

		if k != "zone" {
			// Supporting only zone variable for Exoscale
			continue
		}

		enum := ""
		for _, z := range v.Enum {
			url := strings.Replace(s.URL, "{zone}", z, 1)
			enum += fmt.Sprintf("%s Endpoint = %q\n", helpers.ToCamel(z), url)
		}

		client = Template{
			ServerEndpoint: helpers.ToCamel(v.Default),
			Enum:           enum,
		}
	}

	t, err := template.New("client.tmpl").ParseFiles("./client/client.tmpl")
	if err != nil {
		return nil, err
	}

	output := bytes.NewBuffer([]byte{})
	if err := t.Execute(output, client); err != nil {
		log.Fatal(err)
	}
	return output.Bytes(), nil
}
