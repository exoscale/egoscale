package egoscale

import (
	"context"
	"fmt"
	"os"

	"github.com/exoscale/egoscale/compute"
	"github.com/exoscale/egoscale/dns"
	egoerr "github.com/exoscale/egoscale/error"
	"github.com/exoscale/egoscale/runstatus"
	"github.com/exoscale/egoscale/storage"
	"github.com/pkg/errors" // TODO: replace with Go2-style error wrapping (golang.org/x/xerrors)
)

// Client represents an Exoscale API client.
type Client struct {
	Compute   *compute.Client
	DNS       *dns.Client
	Storage   *storage.Client
	Runstatus *runstatus.Client

	tracing bool
}

// ConfigFunc is a function returning a client configuration profile, or an error if the profile failed to be provided.
type ConfigFunc func() (*ConfigProfile, error)

// ConfigFromProfile returns a ConfigFunc closure generating a client configuration profile from a litterial
// ConfigProfile struct.
func ConfigFromProfile(profile ConfigProfile) ConfigFunc {
	return func() (*ConfigProfile, error) {
		if os.Getenv(traceModeEnvvar) != "" {
			fmt.Fprintf(os.Stderr, "*** ConfigFromProfile: reading configuration from profile: %#v\n", profile)
		}

		return &profile, nil
	}
}

// ConfigFromFile returns a ConfigFunc closure generating a client configuration profile from a configuration file
// located at path.
func ConfigFromFile(path string) ConfigFunc {
	return func() (*ConfigProfile, error) {
		if os.Getenv(traceModeEnvvar) != "" {
			fmt.Fprintf(os.Stderr, "*** ConfigFromFile: reading configuration from file %q\n", path)
		}

		config, err := loadConfig(path)
		if err != nil {
			if os.Getenv(traceModeEnvvar) != "" {
				fmt.Fprintf(os.Stderr, "*** ConfigFromFile: error: %s\n", err)
			}
			return nil, err
		}

		return config.getProfile("")
	}
}

// ConfigFromEnv returns a ConfigFunc closure generating a client configuration profile from user environment (see
// *Envvar constants documentation for a list of supported environment variables).
func ConfigFromEnv() ConfigFunc {
	return func() (*ConfigProfile, error) {
		if os.Getenv(traceModeEnvvar) != "" {
			fmt.Fprintf(os.Stderr, "*** ConfigFromEnv: reading configuration from environment\n")
		}

		return &ConfigProfile{
			APIKey:               os.Getenv(APIKeyEnvvar),
			APISecret:            os.Getenv(APISecretEnvvar),
			ComputeAPIEndpoint:   os.Getenv(ComputeAPIEndpointEnvvar),
			DNSAPIEndpoint:       os.Getenv(DNSAPIEndpointEnvvar),
			RunstatusAPIEndpoint: os.Getenv(RunstatusAPIEndpointEnvvar),
			StorageAPIEndpoint:   os.Getenv(StorageAPIEndpointEnvvar),
			StorageZone:          os.Getenv(StorageZoneEnvvar),
		}, nil
	}
}

// NewClient returns an Exoscale API client with argument cfs being a list of ConfigFunc functions returning a client
// configuration profile, or returns a non-nil error if client configuration failed. If no ConfigFunc is provided, the
// client initialization falls back to the ConfigFromFile($HOME/.exoscale/config.toml) then ConfigFromEnv(). In case
// multiple ConfigFunc are provided, they are processed in the provided order and the first successful execution halts
// the configuration lookup and the returned configuration profile is used to initialize the client.
func NewClient(cfs ...ConfigFunc) (*Client, error) {
	var (
		profile *ConfigProfile
		err     error
		tracing bool
	)

	if os.Getenv(traceModeEnvvar) != "" {
		tracing = true
	}

	// In case no ConfigFunc are provided, fall back to default config file then environment variables.
	// If the ConfigFileEnvvar is non-empty, use it instead of the default config file location.
	if len(cfs) == 0 {
		configFile := defaultConfigFile
		if os.Getenv(ConfigFileEnvvar) != "" {
			configFile = os.Getenv(ConfigFileEnvvar)
		}
		cfs = []ConfigFunc{ConfigFromFile(configFile), ConfigFromEnv()}
	}

	// We stop at the first successful configuration retrieval
	for _, cf := range cfs {
		if profile, err = cf(); err == nil {
			break
		}
	}
	if profile == nil {
		return nil, errors.New("no configuration")
	}

	if profile.APIKey == "" || profile.APISecret == "" {
		return nil, egoerr.ErrMissingAPICredentials
	}

	computeClient, err := compute.NewClient(
		context.Background(),
		profile.APIKey,
		profile.APISecret,
		&compute.ClientOpts{
			APIEndpoint: profile.ComputeAPIEndpoint,
			Tracing:     tracing,
		})
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize Compute API client")
	}

	dnsClient, err := dns.NewClient(
		context.Background(),
		profile.APIKey,
		profile.APISecret,
		&dns.ClientOpts{
			APIEndpoint: profile.DNSAPIEndpoint,
			Tracing:     tracing,
		})
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize DNS API client")
	}

	storageClient, err := storage.NewClient(
		context.Background(),
		profile.APIKey,
		profile.APISecret,
		&storage.ClientOpts{
			APIEndpoint: profile.StorageAPIEndpoint,
			Zone:        profile.StorageZone,
			Tracing:     tracing,
		})
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize Storage API client")
	}

	runstatusClient, err := runstatus.NewClient(
		context.Background(),
		profile.APIKey,
		profile.APISecret,
		&runstatus.ClientOpts{
			APIEndpoint: profile.RunstatusAPIEndpoint,
			Tracing:     tracing,
		})
	if err != nil {
		return nil, errors.Wrap(err, "unable to initialize Runstatus API client")
	}

	return &Client{
		Compute:   computeClient,
		DNS:       dnsClient,
		Storage:   storageClient,
		Runstatus: runstatusClient,
		tracing:   tracing,
	}, nil
}
