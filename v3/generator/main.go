package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/pb33f/libopenapi"

	"github.com/exoscale/egoscale/v3/generator/build"
	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/ir"
	"github.com/exoscale/egoscale/v3/generator/render"
)

//go:generate go run . ./source.yaml ../ v3

//go:embed config.yaml
var configYAML []byte

func main() {
	if len(os.Args) <= 3 {
		fmt.Printf("%s <openAPI-spec.json|yaml> <path generation> <package name>\n", os.Args[0])
		return
	}
	openAPISpec := os.Args[1]
	genPathDir := os.Args[2]
	packageName := os.Args[3]

	if err := run(openAPISpec, genPathDir, packageName); err != nil {
		log.Fatal(err)
	}
}

func run(openAPISpec, genPathDir, packageName string) error {
	cfg, err := config.Parse(configYAML)
	if err != nil {
		return err
	}
	buf, err := os.ReadFile(openAPISpec)
	if err != nil {
		return err
	}

	doc, err := libopenapi.NewDocument(buf)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(genPathDir, os.ModePerm); err != nil {
		return err
	}

	model, errs := doc.BuildV3Model()
	if err := errors.Join(errs...); err != nil {
		return fmt.Errorf("build model: %w", err)
	}

	// Build every file before writing any, not to leave a partially generated package.
	b := build.New(&model.Model, cfg)
	files := []struct {
		name    string
		build   func(string) (ir.File, error)
		content []byte
	}{
		{name: "schemas.go", build: b.Schemas},
		{name: "client.go", build: b.Client},
		{name: "operations.go", build: b.Operations},
	}
	for i, f := range files {
		file, err := f.build(packageName)
		if err != nil {
			return fmt.Errorf("%s: %w", f.name, err)
		}
		if files[i].content, err = source(f.name, file); err != nil {
			return fmt.Errorf("%s: %w", f.name, err)
		}
	}

	for _, f := range files {
		if err := os.WriteFile(filepath.Join(genPathDir, f.name), f.content, 0o644); err != nil {
			return err
		}
	}

	return nil
}

// source renders a file, printing it unformatted when GENERATOR_DEBUG
// is its name without extension (e.g. GENERATOR_DEBUG=schemas).
func source(name string, f ir.File) ([]byte, error) {
	if os.Getenv("GENERATOR_DEBUG") == strings.TrimSuffix(name, ".go") {
		src, err := render.Source(f)
		if err != nil {
			return nil, err
		}
		fmt.Println(string(src))
	}

	return render.File(f)
}
