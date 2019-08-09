package egoscale

import (
	"fmt"
	"os"
	"path"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type clientTestSuite struct {
	suite.Suite
	dir string
}

func (t *clientTestSuite) SetupTest() {
	t.dir = os.TempDir()
}

func (t *clientTestSuite) TearDownSuite() {
	os.RemoveAll(t.dir)
}

func (t *clientTestSuite) TestConfigFromProfile() {
	cf := ConfigFromProfile(ConfigProfile{Name: "alice"})
	profile, err := cf()
	assert.Empty(t.T(), err)
	assert.Equal(t.T(), &ConfigProfile{Name: "alice"}, profile)
}

func (t *clientTestSuite) TestConfigFromFile() {
	var file = path.Join(t.dir, "config.toml")

	assert.Empty(t.T(), configFileFixture(file, fmt.Sprintf(`
default_profile = "bob"

[[profiles]]
name = "alice"
api_key = "%s"
api_secret = "%s"

[[profiles]]
name = "bob"
api_key = "%s"
api_secret = "%s"
`,
		testAliceAPIKey,
		testAliceAPISecret,
		testBobAPIKey,
		testBobAPISecret)))

	cf := ConfigFromFile(file)
	profile, err := cf()
	assert.Empty(t.T(), err)
	assert.Equal(t.T(), &ConfigProfile{
		Name:      "bob",
		APIKey:    testBobAPIKey,
		APISecret: testBobAPISecret,
	}, profile)
}

func (t *clientTestSuite) TestConfigFromEnv() {
	var (
		apiKey               = "apiKey"
		apiSecret            = "apiSecret"
		computeAPIEndpoint   = "computeAPIEndpoint"
		dnsAPIEndpoint       = "dnsAPIEndpoint"
		runstatusAPIEndpoint = "runstatusAPIEndpoint"
		storageAPIEndpoint   = "storageAPIEndpoint"
		storageZone          = "storageZone"
	)

	os.Setenv(APIKeyEnvvar, apiKey)
	os.Setenv(APISecretEnvvar, apiSecret)
	os.Setenv(ComputeAPIEndpointEnvvar, computeAPIEndpoint)
	os.Setenv(DNSAPIEndpointEnvvar, dnsAPIEndpoint)
	os.Setenv(RunstatusAPIEndpointEnvvar, runstatusAPIEndpoint)
	os.Setenv(StorageAPIEndpointEnvvar, storageAPIEndpoint)
	os.Setenv(StorageZoneEnvvar, storageZone)
	defer func() {
		os.Unsetenv(APIKeyEnvvar)
		os.Unsetenv(APISecretEnvvar)
		os.Unsetenv(ComputeAPIEndpointEnvvar)
		os.Unsetenv(DNSAPIEndpointEnvvar)
		os.Unsetenv(RunstatusAPIEndpointEnvvar)
		os.Unsetenv(StorageAPIEndpointEnvvar)
		os.Unsetenv(StorageZoneEnvvar)
	}()

	cf := ConfigFromEnv()
	profile, err := cf()
	assert.Empty(t.T(), err)
	assert.Equal(t.T(), &ConfigProfile{
		APIKey:               apiKey,
		APISecret:            apiSecret,
		ComputeAPIEndpoint:   computeAPIEndpoint,
		DNSAPIEndpoint:       dnsAPIEndpoint,
		RunstatusAPIEndpoint: runstatusAPIEndpoint,
		StorageAPIEndpoint:   storageAPIEndpoint,
		StorageZone:          storageZone,
	}, profile)
}

func (t *clientTestSuite) TestNewClientNoConfig() {
	client, err := NewClient()
	assert.Empty(t.T(), client)
	assert.EqualError(t.T(), err, egoerr.ErrMissingAPICredentials.Error())
}

func (t *clientTestSuite) TestNewClientNoConfigWithEnv() {
	var (
		apiKey               = "apiKey"
		apiSecret            = "apiSecret"
		computeAPIEndpoint   = "computeAPIEndpoint"
		dnsAPIEndpoint       = "dnsAPIEndpoint"
		runstatusAPIEndpoint = "runstatusAPIEndpoint"
		storageAPIEndpoint   = "storageAPIEndpoint"
		storageZone          = "storageZone"
	)

	os.Setenv(APIKeyEnvvar, apiKey)
	os.Setenv(APISecretEnvvar, apiSecret)
	os.Setenv(ComputeAPIEndpointEnvvar, computeAPIEndpoint)
	os.Setenv(DNSAPIEndpointEnvvar, dnsAPIEndpoint)
	os.Setenv(RunstatusAPIEndpointEnvvar, runstatusAPIEndpoint)
	os.Setenv(StorageAPIEndpointEnvvar, storageAPIEndpoint)
	os.Setenv(StorageZoneEnvvar, storageZone)
	defer func() {
		os.Unsetenv(APIKeyEnvvar)
		os.Unsetenv(APISecretEnvvar)
		os.Unsetenv(ComputeAPIEndpointEnvvar)
		os.Unsetenv(DNSAPIEndpointEnvvar)
		os.Unsetenv(RunstatusAPIEndpointEnvvar)
		os.Unsetenv(StorageAPIEndpointEnvvar)
		os.Unsetenv(StorageZoneEnvvar)
	}()

	client, err := NewClient()
	assert.Empty(t.T(), err)
	assert.NotEmpty(t.T(), client)
}

func (t *clientTestSuite) TestNewClientNoConfigWithConfigFileEnv() {
	var file = path.Join(t.dir, "config.toml")

	assert.Empty(t.T(), configFileFixture(file, fmt.Sprintf(`
[[profiles]]
name = "alice"
api_key = "%s"
api_secret = "%s"
`,
		testAliceAPIKey,
		testAliceAPISecret)))

	os.Setenv(ConfigFileEnvvar, file)
	defer os.Unsetenv(ConfigFileEnvvar)

	client, err := NewClient()
	assert.Empty(t.T(), err)
	assert.NotEmpty(t.T(), client)
}

func (t *clientTestSuite) TestNewClientFromProfile() {
	client, err := NewClient(ConfigFromProfile(ConfigProfile{
		APIKey:    "apiKey",
		APISecret: "apiSecret",
	}))
	assert.Empty(t.T(), err)
	assert.NotEmpty(t.T(), client)
}

func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
