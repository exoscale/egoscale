package egoscale

import "io/ioutil"

const (
	testAliceAPIKey    = "alice_api_key"
	testAliceAPISecret = "alice_api_secret"
	testBobAPIKey      = "bob_api_key"
	testBobAPISecret   = "bob_api_secret"
)

func configFileFixture(path string, content string) error {
	return ioutil.WriteFile(path, []byte(content), 0600)
}
