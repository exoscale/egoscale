package main

import (
	"context"
	"fmt"

	"github.com/sauterp/egoscale-rewrite-experiments/oapi-codegen/gen/consumerapi"
)

func main() {
	oapiClient, err := newClient()
	if err != nil {
		panic(err)
	}

	c := consumerapi.Client{
		ClientWithResponses: oapiClient,
	}

	ctx := context.Background()

	i := c.ComputeAPI().Instance().Get(ctx, "5bd8d8c8-c9a8-45dc-ace0-351b30bbb700")
	fmt.Printf("i.Name: %v\n", *i.Name)
}
