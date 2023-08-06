package main

import (
	"fmt"
	"os"
)

// Map hold all mapping information that allows us to generate Consumer API.
// Map is split into groups (subfolder) and each group contains one or many entities.
type Map map[string][]Entity

// Entity represents a resource in Exoscale platform.
type Entity struct {
	// Mandatory
	RootName string //struct name where API functions will attach (filename is strcase.ToSnake(RootName))
	Fns      []Fn   //list of API calls

	// Computed (don't set values manually)
	OAPITypesImport bool
	Package         string
}

type Fn struct {
	// Mandatory
	Name     string //name of the API function
	OAPIName string //base name of the function in oapi

	// Optional
	ResDefOverride         string //override for response definition string
	ResPassthroughOverride string //override for response passthrough string

	// Computed (don't set values manually)
	OptArgsDef         string
	OptArgsPassthrough string
	ResDef             string
	ResPassthrough     string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("available commands: generate, list-unused")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		Generate()
	case "list-unimplemented":
		ListUnimplemented()
	}
}
