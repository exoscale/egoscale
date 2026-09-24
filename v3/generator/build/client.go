package build

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

// Client builds the file of the API client, from the spec server zone endpoints.
func (b *Builder) Client(packageName string) (ir.File, error) {
	// The spec is returning only production server.
	if len(b.model.Servers) == 0 {
		return ir.File{}, errors.New("no server found in the spec")
	}
	if len(b.model.Servers) != 1 {
		slog.Warn("more than one server found, using the first one", slog.Int("servers_len", len(b.model.Servers)))
	}

	srv := b.model.Servers[0]
	if orderedmap.Len(srv.Variables) == 0 {
		return ir.File{}, errors.New("no server variables defined")
	}

	client := ir.Client{}
	// Supporting only zone variable for Exoscale.
	if zone, ok := srv.Variables.Get("zone"); ok {
		for _, z := range zone.Enum {
			client.Endpoints = append(client.Endpoints, ir.Endpoint{
				Name: b.names.Camel(z),
				URL:  strings.Replace(srv.URL, "{zone}", z, 1),
			})
		}
	}
	if err := b.emit(client, "#/servers/0"); err != nil {
		return ir.File{}, err
	}

	return b.file(packageName,
		"context", "fmt", "io", "log", "net/http", "runtime", "time", "",
		"github.com/exoscale/egoscale/v3/credentials",
		"github.com/go-playground/validator/v10",
		"github.com/hashicorp/go-retryablehttp",
	), nil
}
