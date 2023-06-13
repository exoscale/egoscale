package main

import (
	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/consumerapi/generator"
)

func main() {
	if err := generator.Generate(); err != nil {
		panic(err)
	}
}
