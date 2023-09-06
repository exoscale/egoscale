package v3

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
)

const (
	DefaultHostPattern = "https://api-%s.exoscale.com/v2"

	EnvKeyAPIEndpointPattern = "EXOSCALE_API_ENDPOINT_PATTERN"
	EnvKeyAPIEndpointZones   = "EXOSCALE_API_ENDPOINT_ZONES"
)

var (
	// DefaultZones list.
	// When new zone is added or existing removed this slice needs to be updated.
	// First zone in the slice is used as default in DefaultZonedClient.
	DefaultZones = []ZoneName{
		ChGva2,
		AtVie1,
		AtVie2,
		BgSof1,
		ChDk2,
		DeFra1,
		DeMuc1,
	}
)

// ZonedClient is an Exoscale API Client that can communicate with API servers in different zones.
// It has the same interface as Client and uses currently selected zone to run API calls.
// Consumer is expected to select zone before invoking API calls.
type ZonedClient struct {
	zones       map[ZoneName]string
	currentZone ZoneName
	mx          sync.RWMutex

	Client
}

// NewZonedClient creates a new ZonedClient using URL pattern, list of zones and Client options.
// URL pattern must be a valid URL with exactly one substitution verb '%s', for example:
//
//	https://api-%s.exoscale.com/v2
//
// ClientOpt options will be passed down to Client as provided.
// If EXOSCALE_API_ENDPOINT_PATTERN environment variable is set, it replaces urlPattern.
// If EXOSCALE_API_ENDPOINT_ZONES environment variable is set (CSV format), it replaces zones.
func NewZonedClient(urlPattern string, zones []ZoneName, opts ...ClientOpt) (*ZonedClient, error) {
	if len(zones) == 0 {
		return nil, errors.New("list of zones cannot be empty")
	}

	// Env overrides
	if h := os.Getenv(EnvKeyAPIEndpointPattern); h != "" {
		urlPattern = h
	}
	if z := os.Getenv(EnvKeyAPIEndpointZones); z != "" {
		zones = []ZoneName{}
		parts := strings.Split(z, ",")
		for _, part := range parts {
			zones = append(zones, ZoneName(part))
		}
	}

	zonedClient := ZonedClient{
		zones: map[ZoneName]string{},
	}

	for i, zone := range zones {
		url := fmt.Sprintf(urlPattern, zone)

		// Ensure the server URL always has a trailing slash
		if !strings.HasSuffix(url, "/") {
			url += "/"
		}

		zonedClient.zones[zone] = url

		// We use first zone in the provided sice as default.
		if i == 0 {
			zonedClient.currentZone = zone

			client, err := NewClient(url, opts...)
			if err != nil {
				return nil, err
			}

			zonedClient.Client = *client
		}
	}

	return &zonedClient, nil
}

// DefaultClient creates a ZonedClient with preset API URL pattern and zone and provided options.
// This is what should be used by default.
func DefaultClient(opts ...ClientOpt) (*ZonedClient, error) {
	return NewZonedClient(DefaultHostPattern, DefaultZones, opts...)
}

// Zone returns the current zone identifier.
func (c *ZonedClient) Zone() ZoneName {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.currentZone
}

// SetZone selects the current zone.
func (c *ZonedClient) SetZone(z ZoneName) {
	c.mx.Lock()
	c.server = c.zones[z]
	c.currentZone = z
	c.mx.Unlock()
}

// InZone selects the instance of the Client in selected zone so the methods may be chained:
//
//	zonedClient.InZone(ChGva2).IAM()...
func (c *ZonedClient) InZone(z ZoneName) *Client {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return &Client{
		server:     c.zones[z],
		httpClient: c.httpClient,
		reqEditors: c.reqEditors,
		logger:     c.logger,
		creds:      c.creds,
	}
}

// ForEachZone runs function f in each configured zone.
// Argument of function f is configured Client for the zone.
func (c *ZonedClient) ForEachZone(f func(c *Client, zone ZoneName)) {
	for z, url := range c.zones {
		f(&Client{
			server:     url,
			httpClient: c.httpClient,
			reqEditors: c.reqEditors,
			logger:     c.logger,
			creds:      c.creds,
		}, z)
	}
}
