// +build testacc

package compute

import (
	"errors"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type instanceTypeFixture struct {
	c   *Client
	req *egoapi.ListServiceOfferings
	res *egoapi.ServiceOffering
}

func newInstanceTypeFixture(c *Client, opts ...instanceTypeFixtureOpt) *instanceTypeFixture {
	var fixture = &instanceTypeFixture{
		c:   c,
		req: &egoapi.ListServiceOfferings{},
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *instanceTypeFixture) setup() (*instanceTypeFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}

	nt := res.(*egoapi.ListServiceOfferingsResponse).Count
	switch {
	case nt == 0:
		return nil, errors.New("instance type not found")

	case nt > 1:
		return nil, errors.New("multiple results returned, expected only one")

	default:
		for _, instanceType := range res.(*egoapi.ListServiceOfferingsResponse).ServiceOffering {
			instanceType := instanceType
			f.res = &instanceType
		}
	}

	return f, nil
}

func (f *instanceTypeFixture) teardown() error { // nolint:unused,deadcode
	return nil
}

type instanceTypeFixtureOpt func(*instanceTypeFixture)

func instanceTypeFixtureOptID(id string) instanceTypeFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTypeFixture) { f.req.ID = egoapi.MustParseUUID(id) }
}

func instanceTypeFixtureOptName(name string) instanceTypeFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTypeFixture) { f.req.Name = name }
}

func (t *accTestSuite) withInstanceTypeFixture(f func(*instanceTypeFixture), opts ...instanceTypeFixtureOpt) {
	if len(opts) == 0 {
		opts = append(opts, instanceTypeFixtureOptID(testInstanceServiceOfferingID))
	}

	instanceTypeFixture, err := newInstanceTypeFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("instance type fixture setup failed", err)
	}

	f(instanceTypeFixture)
}

type instanceTypeTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *instanceTypeTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *instanceTypeTestSuite) TestListInstanceTypes() {
	var expectedInstanceTypes = []string{
		"Micro",
		"Tiny",
		"Small",
		"Medium",
		"Large",
		"Extra-large",
		"Huge",
		"Mega",
		"Titan",
		"Jumbo",
	}

	instanceTypes, err := t.client.ListInstanceTypes()
	if err != nil {
		t.FailNow("instance types listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(instanceTypes), len(expectedInstanceTypes))
}

func (t *instanceTypeTestSuite) TestGetInstanceTypeByID() {
	instanceType, err := t.client.GetInstanceTypeByID(testInstanceServiceOfferingID)
	if err != nil {
		t.FailNow("instance type retrieval by ID failed", err)
	}
	assert.Equal(t.T(), testInstanceServiceOfferingID, instanceType.ID)
	assert.Equal(t.T(), testInstanceServiceOfferingName, instanceType.Name)

	instanceType, err = t.client.GetInstanceTypeByID("00000000-0000-0000-0000-000000000000")
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), instanceType)
}

func (t *instanceTypeTestSuite) TestGetInstanceTypeByName() {
	instanceType, err := t.client.GetInstanceTypeByName(testInstanceServiceOfferingName)
	if err != nil {
		t.FailNow("instance type retrieval by name failed", err)
	}
	assert.Equal(t.T(), testInstanceServiceOfferingID, instanceType.ID)
	assert.Equal(t.T(), testInstanceServiceOfferingName, instanceType.Name)

	instanceType, err = t.client.GetInstanceTypeByName("lolnope")
	assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
	assert.Empty(t.T(), instanceType)
}

func TestAccComputeInstanceTypeTestSuite(t *testing.T) {
	suite.Run(t, new(instanceTypeTestSuite))
}
