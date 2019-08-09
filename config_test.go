package egoscale

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type configTestSuite struct {
	suite.Suite
	dir string
}

func (t *configTestSuite) SetupTest() {
	t.dir = os.TempDir()
}

func (t *configTestSuite) TearDownSuite() {
	os.RemoveAll(t.dir)
}

func (t *configTestSuite) TestLoadConfig() {
	var file = path.Join(t.dir, "config.toml")

	assert.Empty(t.T(), configFileFixture(file, fmt.Sprintf(`
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

	config, err := loadConfig(file)
	assert.Empty(t.T(), err)
	assert.Len(t.T(), config.Profiles, 2)
	assert.Equal(t.T(), []ConfigProfile{
		{Name: "alice", APIKey: testAliceAPIKey, APISecret: testAliceAPISecret},
		{Name: "bob", APIKey: testBobAPIKey, APISecret: testBobAPISecret},
	}, config.Profiles)
}

func (t *configTestSuite) TestConfigGetProfile() {
	var file = path.Join(t.dir, "config.toml")

	assert.Empty(t.T(), configFileFixture(file, fmt.Sprintf(`
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

	config, err := loadConfig(file)
	assert.Empty(t.T(), err)

	profile, err := config.getProfile("")
	assert.Empty(t.T(), err)
	assert.Equal(t.T(), &ConfigProfile{
		Name:      "alice",
		APIKey:    testAliceAPIKey,
		APISecret: testAliceAPISecret,
	}, profile)

	profile, err = config.getProfile("bob")
	assert.Empty(t.T(), err)
	assert.Equal(t.T(), &ConfigProfile{
		Name:      "bob",
		APIKey:    testBobAPIKey,
		APISecret: testBobAPISecret,
	}, profile)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(configTestSuite))
}
