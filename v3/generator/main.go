package main

import (
	"fmt"
	"os"
)

type Map map[string][]Spec

type Spec struct {
	RootName string
	Fns      []Fn
}

type Fn struct {
	// Mandatory
	Name     string
	OAPIName string

	// Optional
	ResDefOverride         string
	ResPassthroughOverride string

	// Computed, don't set values manually
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
	case "list-unused":
		//ListUnused()
	}
}
