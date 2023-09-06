package main

import (
	"fmt"
	"os"
)

// Map holds all mapping information that allows us to generate API.
// Map is split into platform sections and each section contains one or many resources.
type Map map[string][]Resource

// Resource represents a resource in Exoscale platform.
type Resource struct {
	// Mandatory
	RootName string //struct name where API functions will attach
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
		fmt.Println("Available commands: generate, list-unimplemented")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "generate":
		Generate()
	case "list-unimplemented":
		ListUnimplemented()
	default:
		fmt.Println("Invalid command.")
		fmt.Println("Available commands: generate, list-unimplemented")
		os.Exit(1)
	}
}
