package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pb33f/libopenapi"

	"github.com/exoscale/egoscale/v3/generator/cli"
	"github.com/exoscale/egoscale/v3/generator/client"
	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/exoscale/egoscale/v3/generator/operations"
	"github.com/exoscale/egoscale/v3/generator/schemas"
)

//go:generate go run main.go ./source.yaml ../ v3

func main() {
	if len(os.Args) <= 3 {
		fmt.Printf("%s <openAPI-spec.json|yaml> <path generation> <package name>\n", os.Args[0])
		return
	}
	openAPISpec := os.Args[1]
	genPathDir := os.Args[2]
	packageName := os.Args[3]

	buf, err := os.ReadFile(openAPISpec)
	if err != nil {
		log.Fatal(err)
	}

	doc, err := libopenapi.NewDocument(buf)
	if err != nil {
		log.Fatal(err)
	}

	if err := os.MkdirAll(genPathDir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := schemas.Generate(doc, filepath.Join(genPathDir, "schemas.go"), packageName); err != nil {
		log.Fatal("schemas: ", err)
	}
	if err := client.Generate(doc, filepath.Join(genPathDir, "client.go"), packageName); err != nil {
		log.Fatal("client: ", err)
	}
	if err := operations.Generate(doc, filepath.Join(genPathDir, "operations.go"), packageName); err != nil {
		log.Fatal("operations: ", err)
	}

	if err := os.MkdirAll(genPathDir+"/cli", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	if err := cli.Generate(doc, filepath.Join(genPathDir, "cli", "cli.go"), packageName); err != nil {
		log.Fatal("schemas: ", err)
	}

}

func init() {
	// Additional that are not found here
	// https://github.com/BluntSporks/abbreviation/blob/master/acronyms.go
	// helpers package handle stardard Acronyms.
	helpers.ConfigureAcronym("ssh", "SSH")
	// Exoscale Specifics
	helpers.ConfigureAcronym("iam", "IAM")
	helpers.ConfigureAcronym("sks", "SKS")
	helpers.ConfigureAcronym("sos", "SOS")
	helpers.ConfigureAcronym("dbaas", "DBAAS")
	helpers.ConfigureAcronym("ppapi", "PPAPI")
}
