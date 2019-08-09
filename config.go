package egoscale

import (
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/pelletier/go-toml"
)

const (
	// APIKeyEnvvar is the environment variable indicating the Exoscale API key.
	APIKeyEnvvar = "EXOSCALE_API_KEY"

	// APISecretEnvvar is the environment variable indicating the Exoscale API secret.
	APISecretEnvvar = "EXOSCALE_API_SECRET"

	// ComputeAPIEndpointEnvvar is the environment variable indicating an alternative Exoscale Compute API endpoint.
	ComputeAPIEndpointEnvvar = "EXOSCALE_COMPUTE_API_ENDPOINT"

	// DNSAPIEndpointEnvvar is the environment variable indicating an alternative Exoscale DNS API endpoint.
	DNSAPIEndpointEnvvar = "EXOSCALE_DNS_API_ENDPOINT"

	// RunstatusAPIEndpointEnvvar is the environment variable indicating an alternative Exoscale Runstatus API endpoint.
	RunstatusAPIEndpointEnvvar = "EXOSCALE_RUNSTATUS_API_ENDPOINT"

	// StorageAPIEndpointEnvvar is the environment variable indicating an alternative Exoscale Storage API endpoint.
	StorageAPIEndpointEnvvar = "EXOSCALE_STORAGE_API_ENDPOINT"

	// StorageZoneEnvvar is the environment variable indicating an Exoscale Storage zone.
	StorageZoneEnvvar = "EXOSCALE_STORAGE_ZONE"

	// ConfigFileEnvvar is the environment variable indicating an alternative configuration file location (default is
	// "$HOME/.exoscale/config.toml").
	ConfigFileEnvvar = "EXOSCALE_CONFIG_FILE"

	traceModeEnvvar = "EXOSCALE_TRACE"

	configVersion = 1
)

var (
	defaultConfigFile = path.Join(os.Getenv("HOME"), ".exoscale", "config.toml")
)

// ConfigProfile represents an Exoscale API client configuration profile.
type ConfigProfile struct {
	// Name represents the name of the profile.
	Name string `toml:"name"`

	// APIKey represents the profile Exoscale client API key.
	APIKey string `toml:"api_key"`

	// APISecret represents the profile Exoscale client API secret.
	APISecret string `toml:"api_secret"`

	// ComputeAPIEndpoint represents an alternative Exoscale Compute API endpoint.
	ComputeAPIEndpoint string `toml:"compute_api_endpoint"`

	// DNSAPIEndpoint represents an alternative Exoscale DNS API endpoint.
	DNSAPIEndpoint string `toml:"dns_api_endpoint"`

	// RunstatusAPIEndpoint represents an alternative Exoscale Runstatus API endpoint.
	RunstatusAPIEndpoint string `toml:"runstatus_api_endpoint"`

	// StorageAPIEndpoint represents an alternative Exoscale Storage API endpoint.
	StorageAPIEndpoint string `toml:"storage_api_endpoint"`

	// StorageZone represents an Exoscale Object Storage zone.
	StorageZone string `toml:"storage_zone"`
}

type config struct {
	Version        int             `toml:"version"`
	DefaultProfile string          `toml:"default_profile"`
	Profiles       []ConfigProfile `toml:"profiles"`
}

func loadConfig(path string) (*config, error) {
	var config config

	tc, err := toml.LoadFile(path)
	if err != nil {
		return nil, err
	}

	if err := tc.Unmarshal(&config); err != nil {
		return nil, err
	}

	if config.Version < configVersion {
		// TODO: legacy configuration format migration
	}

	return &config, nil
}

func (c *config) getProfile(name string) (*ConfigProfile, error) {
	if len(c.Profiles) == 0 {
		return nil, errors.New("no profiles configured")
	}

	if name != "" || c.DefaultProfile != "" {
		if name == "" && c.DefaultProfile != "" {
			name = c.DefaultProfile
		}

		for _, profile := range c.Profiles {
			profile := profile
			if profile.Name == name {
				return &profile, nil
			}
		}

		return nil, fmt.Errorf("profile %q not found", name)
	}

	return &c.Profiles[0], nil
}
