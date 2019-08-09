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

type instanceTemplateFixture struct {
	c   *Client
	req *egoapi.ListTemplates
	res *egoapi.Template
}

func newInstanceTemplateFixture(c *Client, opts ...instanceTemplateFixtureOpt) *instanceTemplateFixture {
	var fixture = &instanceTemplateFixture{
		c: c,
		req: &egoapi.ListTemplates{
			TemplateFilter: "featured", // TODO: make a map of public/internal filter labels
		},
	}

	// Fixture default options
	for _, opt := range []instanceTemplateFixtureOpt{
		instanceTemplateFixtureOptZone(testZoneID),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *instanceTemplateFixture) setup() (*instanceTemplateFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}

	nt := res.(*egoapi.ListTemplatesResponse).Count
	switch {
	case nt == 0:
		return nil, errors.New("instance template not found")

	case nt > 1:
		return nil, errors.New("multiple results returned, expected only one")

	default:
		for _, instanceTemplate := range res.(*egoapi.ListTemplatesResponse).Template {
			instanceTemplate := instanceTemplate
			f.res = &instanceTemplate
		}
	}

	return f, nil
}

func (f *instanceTemplateFixture) teardown() error { // nolint:unused,deadcode
	return nil
}

type instanceTemplateFixtureOpt func(*instanceTemplateFixture)

func instanceTemplateFixtureOptZone(id string) instanceTemplateFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTemplateFixture) { f.req.ZoneID = egoapi.MustParseUUID(id) }
}

func instanceTemplateFixtureOptID(id string) instanceTemplateFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTemplateFixture) { f.req.ID = egoapi.MustParseUUID(id) }
}

func instanceTemplateFixtureOptName(name string) instanceTemplateFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTemplateFixture) { f.req.Name = name }
}

func instanceTemplateFixtureOptFilter(filter string) instanceTemplateFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceTemplateFixture) { f.req.TemplateFilter = filter }
}

func (t *accTestSuite) withInstanceTemplateFixture(f func(instance *instanceTemplateFixture),
	opts ...instanceTemplateFixtureOpt) {
	if len(opts) == 0 {
		opts = append(opts, instanceTemplateFixtureOptID(testInstanceTemplateID))
	}

	instanceTemplateFixture, err := newInstanceTemplateFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("instance template fixture setup failed", err)
	}

	f(instanceTemplateFixture)
}

type instanceTemplateTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *instanceTemplateTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *instanceTemplateTestSuite) TestListInstanceTemplates() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {
		zone := t.client.zoneFromAPI(zoneFixture.res)

		instanceTemplates, err := t.client.ListInstanceTemplates(zone, "", "exoscale")
		if err != nil {
			t.FailNow("instance types listing failed", err)
		}
		assert.GreaterOrEqual(t.T(), len(instanceTemplates), 10)
	})
}

func (t *instanceTemplateTestSuite) TestGetInstanceTemplate() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {
		zone := t.client.zoneFromAPI(zoneFixture.res)

		instanceTemplate, err := t.client.GetInstanceTemplate(zone, testInstanceTemplateID, "exoscale")
		if err != nil {
			t.FailNow("instance template retrieval failed", err)
		}
		assert.Equal(t.T(), testInstanceTemplateID, instanceTemplate.ID)
		assert.Equal(t.T(), testInstanceTemplateName, instanceTemplate.Name)

		instanceTemplate, err = t.client.GetInstanceTemplate(zone, "00000000-0000-0000-0000-000000000000", "exoscale")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), instanceTemplate)
	})
}

func TestAccComputeInstanceTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(instanceTemplateTestSuite))
}
