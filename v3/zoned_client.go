package v3

import (
	"errors"
	"fmt"
	"sync"

	"github.com/exoscale/egoscale/v3/oapi"
)

const (
	defaultHostPattern = "https://api-%s.exoscale.com/v2"
)

var (
	// Default zones list (available in oapi code).
	// When new zone is added or existing removed this slice needs to be updated.
	// First zone in the slice is used as default in DefaultZonedClient.
	defaultZones = []oapi.ZoneName{
		oapi.ChGva2,
		oapi.AtVie1,
		oapi.AtVie2,
		oapi.BgSof1,
		oapi.ChDk2,
		oapi.DeFra1,
		oapi.DeMuc1,
	}
)

// ZonedClient is a Exoscale API Client that can communicate with API servers in different zones.
// It has the same interface as Client and uses currently seleced zone to run API calls.
// Consumer is expected to select zone before invoking API calls.
type ZonedClient struct {
	zones       map[oapi.ZoneName]*oapi.ClientWithResponses
	currentZone oapi.ZoneName
	mx          sync.RWMutex

	Client
}

// NewZonedClient creates a new ZonedClient using URL pattern, list of zones and Client options.
// URL pattern must be a valid URL with exactly one substitution verb '%s', for example:
//
//	https://api-%s.exoscale.com/v2
//
// ClientOpt options will be passed down to Client as provided.
func NewZonedClient(urlPattern string, zones []oapi.ZoneName, opts ...ClientOpt) (*ZonedClient, error) {
	if len(zones) == 0 {
		return nil, errors.New("List of zones cannot be empty")
	}

	zonedClient := ZonedClient{
		zones: map[oapi.ZoneName]*oapi.ClientWithResponses{},
	}

	for _, zone := range zones {
		client, err := NewClient(fmt.Sprintf(urlPattern, zone), opts...)
		if err != nil {
			return nil, err
		}

		if zonedClient.creds == nil {
			zonedClient.creds = client.creds
		}

		zonedClient.zones[zone] = client.oapiClient
	}

	// Set default zone to first zone in the provided slice
	zonedClient.currentZone = zones[0]
	zonedClient.oapiClient = zonedClient.zones[zones[0]]

	return &zonedClient, nil
}

// DefaultZonedClient creates a ZonedClient with preset API URL pattern and zone and provided options.
// This is what should be used by default.
func DefaultZonedClient(opts ...ClientOpt) (*ZonedClient, error) {
	return NewZonedClient(defaultHostPattern, defaultZones, opts...)
}

// Zone returns the current zone identifier.
func (c *ZonedClient) Zone() oapi.ZoneName {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.currentZone
}

// SetZone selects the current zone.
func (c *ZonedClient) SetZone(z oapi.ZoneName) {
	c.mx.Lock()
	c.oapiClient = c.zones[z]
	c.mx.Unlock()
}

// InZone selects the current zone and returns instance of the ZonedClient so the methods may be chained:
//
//	zonedClient.InZone(oapi.ChGva2).OAPIClient()...
func (c *ZonedClient) InZone(z oapi.ZoneName) *ZonedClient {
	c.SetZone(z)
	return c
}

// OAPIClient returns configured instance of OpenAPI generated (low-level) API client in the selected zone.
func (c *ZonedClient) OAPIClient() *oapi.ClientWithResponses {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.Client.OAPIClient()
}
