package build

import (
	"errors"
	"log/slog"
	"strings"

	"github.com/pb33f/libopenapi/orderedmap"

	"github.com/exoscale/egoscale/v3/generator/ir"
)

// Constants builds the file of the zone endpoints, from the spec server.
func (b *Builder) Constants(packageName string) (ir.File, error) {
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

	constants := ir.Constants{}
	// Supporting only zone variable for Exoscale.
	if zone, ok := srv.Variables.Get("zone"); ok {
		for _, z := range zone.Enum {
			constants.Endpoints = append(constants.Endpoints, ir.Endpoint{
				Name: b.names.Camel(z),
				Zone: z,
				URL:  strings.Replace(srv.URL, "{zone}", z, 1),
			})
		}
	}
	if err := b.emit(constants, "#/servers/0"); err != nil {
		return ir.File{}, err
	}

	return b.file(packageName), nil
}
