// +build testacc

package compute

import (
	"errors"
	"fmt"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type zoneFixture struct {
	c   *Client
	req *egoapi.ListZones
	res *egoapi.Zone
}

func newZoneFixture(c *Client, opts ...zoneFixtureOpt) *zoneFixture {
	var fixture = &zoneFixture{
		c:   c,
		req: &egoapi.ListZones{},
	}

	// Fixture default options
	for _, opt := range []zoneFixtureOpt{
		zoneFixtureOptName(testZoneName),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *zoneFixture) setup() (*zoneFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}

	nz := res.(*egoapi.ListZonesResponse).Count
	switch {
	case nz == 0:
		return nil, fmt.Errorf("zone %q not found", f.req.Name)

	case nz > 1:
		return nil, errors.New("multiple results returned, expected only one")

	default:
		for _, zone := range res.(*egoapi.ListZonesResponse).Zone {
			zone := zone
			f.res = &zone
		}
	}

	return f, nil
}

func (f *zoneFixture) teardown() error { // nolint:unused,deadcode
	return nil
}

type zoneFixtureOpt func(*zoneFixture)

func zoneFixtureOptName(name string) zoneFixtureOpt { // nolint:unused,deadcode
	return func(f *zoneFixture) { f.req.Name = name }
}

func (t *accTestSuite) withZoneFixture(f func(*zoneFixture), opts ...zoneFixtureOpt) {
	zoneFixture, err := newZoneFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("zone fixture setup failed", err)
	}

	f(zoneFixture)
}

type zoneTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *zoneTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *zoneTestSuite) TestListZones() {
	var expectedZones = []string{
		"at-vie-1",
		"bg-sof-1",
		"ch-dk-2",
		"ch-gva-2",
		"de-fra-1",
		"de-muc-1",
	}

	zones, err := t.client.ListZones()
	if err != nil {
		t.FailNow("zones listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(zones), len(expectedZones))
}

func (t *zoneTestSuite) TestGetZoneByID() {
	zone, err := t.client.GetZoneByID(testZoneID)
	if err != nil {
		t.FailNow("zone retrieval by ID failed", err)
	}
	assert.Equal(t.T(), testZoneID, zone.ID)
	assert.Equal(t.T(), testZoneName, zone.Name)

	zone, err = t.client.GetZoneByID("00000000-0000-0000-0000-000000000000")
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), zone)
}

func (t *zoneTestSuite) TestGetZoneByName() {
	zone, err := t.client.GetZoneByName(testZoneName)
	if err != nil {
		t.FailNow("zone retrieval by name failed", err)
	}
	assert.Equal(t.T(), testZoneID, zone.ID)
	assert.Equal(t.T(), testZoneName, zone.Name)

	zone, err = t.client.GetZoneByName("lolnope")
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), zone)
}

func TestAccComputeZoneTestSuite(t *testing.T) {
	suite.Run(t, new(zoneTestSuite))
}
