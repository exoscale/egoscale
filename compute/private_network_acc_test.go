// +build testacc

package compute

import (
	"encoding/json"
	"net"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type privateNetworkFixture struct {
	c   *Client
	req *egoapi.CreateNetwork
	res *egoapi.Network
}

func newPrivateNetworkFixture(c *Client, opts ...privateNetworkFixtureOpt) *privateNetworkFixture {
	var fixture = &privateNetworkFixture{
		c:   c,
		req: &egoapi.CreateNetwork{},
	}

	// Fixture default options
	for _, opt := range []privateNetworkFixtureOpt{
		privateNetworkFixtureOptZone(testZoneID),
		privateNetworkFixtureOptName(testPrefix + "-" + testRandomString()),
		privateNetworkFixtureOptDescription(testDescription),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *privateNetworkFixture) setup() (*privateNetworkFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoapi.Network)

	return f, nil
}

func (f *privateNetworkFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DeleteNetwork{ID: f.res.ID})
	return f.c.csError(err)
}

type privateNetworkFixtureOpt func(*privateNetworkFixture)

func privateNetworkFixtureOptZone(id string) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.ZoneID = egoapi.MustParseUUID(id) }
}

func privateNetworkFixtureOptName(name string) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.Name = name }
}

func privateNetworkFixtureOptDescription(description string) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.DisplayText = description }
}

func privateNetworkFixtureOptStartIP(ip net.IP) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.StartIP = ip }
}

func privateNetworkFixtureOptEndIP(ip net.IP) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.EndIP = ip }
}

func privateNetworkFixtureOptNetmask(ip net.IP) privateNetworkFixtureOpt { // nolint:unused,deadcode
	return func(f *privateNetworkFixture) { f.req.Netmask = ip }
}

func (t *accTestSuite) withPrivateNetworkFixture(f func(*privateNetworkFixture), opts ...privateNetworkFixtureOpt) {
	privateNetworkFixture, err := newPrivateNetworkFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("Private Network fixture setup failed", err)
	}

	f(privateNetworkFixture)
}

type privateNetworkTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *privateNetworkTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *privateNetworkTestSuite) TestCreatePrivateNetwork() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {
		var (
			zone                  = t.client.zoneFromAPI(zoneFixture.res)
			privateNetworkName    = testPrefix + "-" + testRandomString()
			privateNetworkStartIP = net.ParseIP("10.0.0.1")
			privateNetworkEndIP   = net.ParseIP("10.0.0.5")
			privateNetworkNetmask = net.ParseIP("255.255.255.0")
		)

		privateNetwork, err := t.client.CreatePrivateNetwork(
			zone,
			&PrivateNetworkCreateOpts{
				Name:        privateNetworkName,
				Description: testDescription,
				StartIP:     privateNetworkStartIP,
				EndIP:       privateNetworkEndIP,
				Netmask:     privateNetworkNetmask,
			})
		if err != nil {
			t.FailNow("Private Network creation failed", err)
		}
		assert.NotEmpty(t.T(), privateNetwork.ID)

		actualPrivateNetwork := egoapi.Network{}
		if err := json.Unmarshal(privateNetwork.Raw(), &actualPrivateNetwork); err != nil {
			t.FailNow("unable to unmarshal raw resource", err)
		}

		assert.Equal(t.T(), zone.ID, privateNetwork.Zone.ID)
		assert.Equal(t.T(), privateNetworkName, actualPrivateNetwork.Name)
		assert.Equal(t.T(), privateNetworkName, privateNetwork.Name)
		assert.Equal(t.T(), testDescription, actualPrivateNetwork.DisplayText)
		assert.Equal(t.T(), testDescription, privateNetwork.Description)
		assert.True(t.T(), privateNetwork.StartIP.Equal(privateNetworkStartIP))
		assert.True(t.T(), privateNetwork.EndIP.Equal(privateNetworkEndIP))
		assert.True(t.T(), privateNetwork.Netmask.Equal(privateNetworkNetmask))

		if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DeleteNetwork{
			ID: egoapi.MustParseUUID(privateNetwork.ID),
		}); err != nil {
			t.FailNow("Private Network deletion failed", err)
		}
	})
}

func (t *privateNetworkTestSuite) TestListPrivateNetworks() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withPrivateNetworkFixture(func(privateNetworkFixture *privateNetworkFixture) {
			defer privateNetworkFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)
			privateNetworks, err := t.client.ListPrivateNetworks(zone)
			if err != nil {
				t.FailNow("Private Networks listing failed", err)
			}

			// We cannot guarantee that there will be only our resources in the
			// testing environment, so we ensure we get at least our fixture Private Network
			assert.GreaterOrEqual(t.T(), len(privateNetworks), 1)
		})
	})
}

func (t *privateNetworkTestSuite) TestGetPrivateNetwork() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withPrivateNetworkFixture(func(privateNetworkFixture *privateNetworkFixture) {
			defer privateNetworkFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)
			privateNetwork, err := t.client.GetPrivateNetwork(zone, privateNetworkFixture.res.ID.String())
			if err != nil {
				t.FailNow("Private Network retrieval failed", err)
			}
			assert.Equal(t.T(), privateNetworkFixture.res.ID.String(), privateNetwork.ID)

			privateNetwork, err = t.client.GetPrivateNetwork(zone, "00000000-0000-0000-0000-000000000000")
			assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
			assert.Empty(t.T(), privateNetwork)
		})
	})
}

// TODO: func (t *privateNetworkTestSuite) TestPrivateNetworkInstances() {}

// TODO: func (t *privateNetworkTestSuite) TestPrivateNetworkAttachInstance() {}

// TODO: func (t *privateNetworkTestSuite) TestPrivateNetworkDetachInstance() {}

func (t *privateNetworkTestSuite) TestPrivateNetworkUpdate() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withPrivateNetworkFixture(func(privateNetworkFixture *privateNetworkFixture) {
			defer privateNetworkFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)
			privateNetwork, err := t.client.privateNetworkFromAPI(privateNetworkFixture.res)
			if err != nil {
				t.FailNow("Private Network fixture setup failed", err)
			}

			privateNetworkNameEdited := privateNetwork.Name + " (edited)"
			privateNetworkDescriptionEdited := privateNetwork.Description + " (edited)"
			privateNetworkStartIPEdited := net.ParseIP("10.0.0.1")
			privateNetworkEndIPEdited := net.ParseIP("10.0.0.100")
			privateNetworkNetmaskEdited := net.ParseIP("255.0.0.0")

			if err := privateNetwork.Update(
				&PrivateNetworkUpdateOpts{
					Name:        privateNetworkNameEdited,
					Description: privateNetworkDescriptionEdited,
					StartIP:     privateNetworkStartIPEdited,
					EndIP:       privateNetworkEndIPEdited,
					Netmask:     privateNetworkNetmaskEdited,
				}); err != nil {
				t.FailNow("Private Network update failed", err)
			}

			res, err := t.client.c.ListWithContext(t.client.ctx, &egoapi.Network{
				ID:     egoapi.MustParseUUID(privateNetwork.ID),
				ZoneID: egoapi.MustParseUUID(zone.ID),
			})
			if err != nil || len(res) == 0 {
				t.FailNow("Private Network retrieval failed", err)
			}
			actualPrivateNetwork := res[0].(*egoapi.Network)

			assert.Equal(t.T(), privateNetworkNameEdited, actualPrivateNetwork.Name)
			assert.Equal(t.T(), privateNetworkNameEdited, privateNetwork.Name)
			assert.Equal(t.T(), privateNetworkDescriptionEdited, actualPrivateNetwork.DisplayText)
			assert.Equal(t.T(), privateNetworkDescriptionEdited, privateNetwork.Description)
			assert.Equal(t.T(), privateNetworkStartIPEdited, actualPrivateNetwork.StartIP)
			assert.Equal(t.T(), privateNetworkStartIPEdited, privateNetwork.StartIP)
			assert.Equal(t.T(), privateNetworkEndIPEdited, actualPrivateNetwork.EndIP)
			assert.Equal(t.T(), privateNetworkEndIPEdited, privateNetwork.EndIP)
			assert.Equal(t.T(), privateNetworkNetmaskEdited, actualPrivateNetwork.Netmask)
			assert.Equal(t.T(), privateNetworkNetmaskEdited, privateNetwork.Netmask)
		},
			privateNetworkFixtureOptStartIP(net.ParseIP("10.0.0.10")),
			privateNetworkFixtureOptEndIP(net.ParseIP("10.0.0.50")),
			privateNetworkFixtureOptNetmask(net.ParseIP("255.255.255.0")),
		)
	})
}

func (t *privateNetworkTestSuite) TestPrivateNetworkDelete() {
	t.withPrivateNetworkFixture(func(privateNetworkFixture *privateNetworkFixture) {
		privateNetwork, err := t.client.privateNetworkFromAPI(privateNetworkFixture.res)
		if err != nil {
			t.FailNow("Private Network fixture setup failed", err)
		}

		if err = privateNetwork.Delete(); err != nil {
			t.FailNow("Private Network deletion failed", err)
		}
		assert.Empty(t.T(), privateNetwork.ID)
		assert.Empty(t.T(), privateNetwork.Name)
		assert.Empty(t.T(), privateNetwork.Description)
		assert.Empty(t.T(), privateNetwork.StartIP)
		assert.Empty(t.T(), privateNetwork.EndIP)
		assert.Empty(t.T(), privateNetwork.Netmask)

		r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.Network{ID: privateNetworkFixture.res.ID})
		assert.Len(t.T(), r, 0)
	})
}

func TestAccComputePrivateNetworkTestSuite(t *testing.T) {
	suite.Run(t, new(privateNetworkTestSuite))
}
