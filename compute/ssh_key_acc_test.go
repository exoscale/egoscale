// +build testacc

package compute

import (
	"encoding/json"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type sshKeyFixture struct {
	c   *Client
	req *egoapi.CreateSSHKeyPair
	res *egoapi.SSHKeyPair
}

func newSSHKeyFixture(c *Client, opts ...sshKeyFixtureOpt) *sshKeyFixture {
	var fixture = &sshKeyFixture{
		c:   c,
		req: &egoapi.CreateSSHKeyPair{},
	}

	// Fixture default options
	for _, opt := range []sshKeyFixtureOpt{
		sshKeyFixtureOptName(testPrefix + "-" + testRandomString()),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *sshKeyFixture) setup() (*sshKeyFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoapi.SSHKeyPair)

	return f, nil
}

func (f *sshKeyFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DeleteSSHKeyPair{Name: f.res.Name})
	return f.c.csError(err)
}

type sshKeyFixtureOpt func(*sshKeyFixture)

func sshKeyFixtureOptName(name string) sshKeyFixtureOpt { // nolint:unused,deadcode
	return func(f *sshKeyFixture) { f.req.Name = name }
}

func (t *accTestSuite) withSSHKeyFixture(f func(*sshKeyFixture), opts ...sshKeyFixtureOpt) {
	sshKeyFixture, err := newSSHKeyFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("SSH key fixture setup failed", err)
	}

	f(sshKeyFixture)
}

type sshKeyTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *sshKeyTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *sshKeyTestSuite) TestCreateSSHKey() {
	var sshKeyName = testPrefix + "-" + testRandomString()

	sshKey, err := t.client.CreateSSHKey(sshKeyName)
	if err != nil {
		t.FailNow("SSH key creation failed", err)
	}
	assert.NotEmpty(t.T(), sshKey.Name)

	actualSSHKey := egoapi.SSHKeyPair{}
	if err := json.Unmarshal(sshKey.Raw(), &actualSSHKey); err != nil {
		t.FailNow("unable to unmarshal raw resource", err)
	}

	assert.Equal(t.T(), sshKeyName, actualSSHKey.Name)
	assert.Equal(t.T(), sshKeyName, sshKey.Name)
	assert.NotEmpty(t.T(), sshKey.Fingerprint)
	assert.NotEmpty(t.T(), sshKey.PrivateKey)

	if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DeleteSSHKeyPair{Name: sshKey.Name}); err != nil {
		t.FailNow("SSH key deletion failed", err)
	}
}

func (t *sshKeyTestSuite) TestRegisterSSHKey() {
	var (
		sshKeyName   = testPrefix + "-" + testRandomString()
		sshPublicKey = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDGRYWaNYBG/Ld3ZnXGsK9pZ" +
			"l9kT3B6GXvsslgy/LCjkJvDIP+nL+opAArKZD1P1+SGylCLt8ISdJNNGLtxKp9CL12EGAYqd" +
			"Dvm5PurkpqIkEsfhsIG4dne9hNu7ZW8aHGHDWM62/4uiWOKtbGdv/P33L/FepzypwpivFsaX" +
			"wPYVunAgoBQLUAmj/xcwtx7cvKS4zdj0+Iu21CIGU9wsH3ZLS34QiXtCGJyMOp158qld9Oeu" +
			"s3Y/7DQ4w5XvfGn9sddxHOSMwUlNiFVty673X3exgMIc8psZOsHvWZPS0zWx9gEDE95cUU10" +
			"K6u4vzTr2O6fgDOQBynEUw3CDiHvwRD alice@example.net"
		sshKeyFingerprint = "a0:25:fa:32:c0:18:7a:f8:e8:b2:3b:30:d8:ca:9a:2e"
	)

	sshKey, err := t.client.RegisterSSHKey(sshKeyName, sshPublicKey)
	if err != nil {
		t.FailNow("SSH key registration failed", err)
	}
	assert.Equal(t.T(), sshKeyName, sshKey.Name)
	assert.Equal(t.T(), sshKeyFingerprint, sshKey.Fingerprint)

	if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DeleteSSHKeyPair{Name: sshKey.Name}); err != nil {
		t.FailNow("SSH key deletion failed", err)
	}
}

func (t *sshKeyTestSuite) TestListSSHKeys() {
	t.withSSHKeyFixture(func(sshKeyFixture *sshKeyFixture) {
		defer sshKeyFixture.teardown() // nolint:errcheck

		sshKeys, err := t.client.ListSSHKeys()
		if err != nil {
			t.FailNow("SSH keys listing failed", err)
		}

		// We cannot guarantee that there will be only our resources in the
		// testing environment, so we ensure we get at least our fixture SSH key
		assert.GreaterOrEqual(t.T(), len(sshKeys), 1)
	})
}

func (t *sshKeyTestSuite) TestGetSSHKey() {
	var sshKeyName = testPrefix + "-" + testRandomString()

	t.withSSHKeyFixture(func(sshKeyFixture *sshKeyFixture) {
		defer sshKeyFixture.teardown() // nolint:errcheck

		sshKey, err := t.client.GetSSHKey(sshKeyFixture.res.Name)
		if err != nil {
			t.FailNow("SSH key retrieval failed", err)
		}
		assert.Equal(t.T(), sshKeyFixture.res.Name, sshKey.Name)

		sshKey, err = t.client.GetSSHKey("lolnope")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), sshKey)
	}, sshKeyFixtureOptName(sshKeyName))
}

func (t *sshKeyTestSuite) TesteSSHKeyDelete() {
	t.withSSHKeyFixture(func(sshKeyFixture *sshKeyFixture) {
		defer sshKeyFixture.teardown() // nolint:errcheck

		sshKey := t.client.sshKeyFromAPI(sshKeyFixture.res)
		sshKeyName := sshKey.Name

		if err := sshKey.Delete(); err != nil {
			t.FailNow("SSH key deletion failed", err)
		}
		assert.Empty(t.T(), sshKey.Name)
		assert.Empty(t.T(), sshKey.Fingerprint)
		assert.Empty(t.T(), sshKey.PrivateKey)

		r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.SSHKeyPair{Name: sshKeyName})
		assert.Len(t.T(), r, 0)
	})
}

func TestAccComputeSSHKeyTestSuite(t *testing.T) {
	suite.Run(t, new(sshKeyTestSuite))
}
