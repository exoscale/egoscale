package main

import (
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"

	"github.com/exoscale/egoscale/v3/generator/client"
	"github.com/exoscale/egoscale/v3/generator/config"
	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/exoscale/egoscale/v3/generator/operations"
	"github.com/exoscale/egoscale/v3/generator/schemas"
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
	helpers.Configure(cfg)

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

	if err := schemas.Generate(doc, filepath.Join(genPathDir, "/schemas.go"), packageName); err != nil {
		return fmt.Errorf("schemas: %w", err)
	}
	if err := client.Generate(doc, filepath.Join(genPathDir, "/client.go"), packageName); err != nil {
		return fmt.Errorf("client: %w", err)
	}
	if err := operations.Generate(doc, filepath.Join(genPathDir, "/operations.go"), packageName); err != nil {
		return fmt.Errorf("operations: %w", err)
	}

	return nil
}
