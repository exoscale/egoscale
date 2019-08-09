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

type antiAffinityGroupFixture struct {
	c   *Client
	req *egoapi.CreateAffinityGroup
	res *egoapi.AffinityGroup
}

func newAntiAffinityGroupFixture(c *Client, opts ...antiAffinityGroupFixtureOpt) *antiAffinityGroupFixture {
	var fixture = &antiAffinityGroupFixture{
		c:   c,
		req: &egoapi.CreateAffinityGroup{Type: "host anti-affinity"},
	}

	// Fixture default options
	for _, opt := range []antiAffinityGroupFixtureOpt{
		antiAffinityGroupFixtureOptName(testPrefix + "-" + testRandomString()),
		antiAffinityGroupFixtureOptDescription(testDescription),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *antiAffinityGroupFixture) setup() (*antiAffinityGroupFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoapi.AffinityGroup)

	return f, nil
}

func (f *antiAffinityGroupFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DeleteAffinityGroup{ID: f.res.ID})
	return f.c.csError(err)
}

type antiAffinityGroupFixtureOpt func(*antiAffinityGroupFixture)

func antiAffinityGroupFixtureOptName(name string) antiAffinityGroupFixtureOpt { // nolint:unused,deadcode
	return func(f *antiAffinityGroupFixture) { f.req.Name = name }
}

func antiAffinityGroupFixtureOptDescription(description string) antiAffinityGroupFixtureOpt { // nolint:unused,deadcode
	return func(f *antiAffinityGroupFixture) { f.req.Description = description }
}

func (t *accTestSuite) withAntiAffinityGroupFixture(f func(*antiAffinityGroupFixture),
	opts ...antiAffinityGroupFixtureOpt) {
	antiAffinityGroupFixture, err := newAntiAffinityGroupFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("Anti-Affinity Group fixture setup failed", err)
	}

	f(antiAffinityGroupFixture)
}

type antiAffinityGroupTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *antiAffinityGroupTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *antiAffinityGroupTestSuite) TestCreateAntiAffinityGroup() {
	var antiAffinityGroupName = testPrefix + "-" + testRandomString()

	antiAffinityGroup, err := t.client.CreateAntiAffinityGroup(
		antiAffinityGroupName,
		&AntiAffinityGroupCreateOpts{Description: testDescription},
	)
	if err != nil {
		t.FailNow("Anti-Affinity Group creation failed", err)
	}
	assert.NotEmpty(t.T(), antiAffinityGroup.ID)

	actualAntiAffinityGroup := egoapi.AffinityGroup{}
	if err := json.Unmarshal(antiAffinityGroup.Raw(), &actualAntiAffinityGroup); err != nil {
		t.FailNow("unable to unmarshal raw resource", err)
	}

	assert.Equal(t.T(), antiAffinityGroupName, actualAntiAffinityGroup.Name)
	assert.Equal(t.T(), antiAffinityGroupName, antiAffinityGroup.Name)
	assert.Equal(t.T(), testDescription, actualAntiAffinityGroup.Description)
	assert.Equal(t.T(), testDescription, antiAffinityGroup.Description)

	if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DeleteAffinityGroup{
		ID: egoapi.MustParseUUID(antiAffinityGroup.ID),
	}); err != nil {
		t.FailNow("Anti-Affinity Group deletion failed", err)
	}
}

func (t *antiAffinityGroupTestSuite) TestListAntiAffinityGroups() {
	t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
		defer antiAffinityGroupFixture.teardown() // nolint:errcheck

		antiAffinityGroups, err := t.client.ListAntiAffinityGroups()
		if err != nil {
			t.FailNow("Anti-Affinity Groups listing failed", err)
		}

		// We cannot guarantee that there will be only our resources in the
		// testing environment, so we ensure we get at least our fixture AAG
		assert.GreaterOrEqual(t.T(), len(antiAffinityGroups), 1)
	})
}

func (t *antiAffinityGroupTestSuite) TestGetAntiAffinityGroupByID() {
	t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
		defer antiAffinityGroupFixture.teardown() // nolint:errcheck

		antiAffinityGroup, err := t.client.GetAntiAffinityGroupByID(antiAffinityGroupFixture.res.ID.String())
		if err != nil {
			t.FailNow("Anti-Affinity Group retrieval by ID failed", err)
		}
		assert.Equal(t.T(), antiAffinityGroupFixture.res.ID.String(), antiAffinityGroup.ID)
		assert.Equal(t.T(), antiAffinityGroupFixture.res.Name, antiAffinityGroup.Name)
		assert.Equal(t.T(), antiAffinityGroupFixture.res.Description, antiAffinityGroup.Description)

		antiAffinityGroup, err = t.client.GetAntiAffinityGroupByID("00000000-0000-0000-0000-000000000000")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), antiAffinityGroup)
	})
}

func (t *antiAffinityGroupTestSuite) TestGetAntiAffinityGroupByName() {
	var antiAffinityGroupName = testPrefix + "-" + testRandomString()

	t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
		defer antiAffinityGroupFixture.teardown() // nolint:errcheck

		antiAffinityGroup, err := t.client.GetAntiAffinityGroupByName(antiAffinityGroupFixture.res.Name)
		if err != nil {
			t.FailNow("Anti-Affinity Group retrieval by name failed", err)
		}
		assert.Equal(t.T(), antiAffinityGroupFixture.res.ID.String(), antiAffinityGroup.ID)
		assert.Equal(t.T(), antiAffinityGroupFixture.res.Name, antiAffinityGroup.Name)
		assert.Equal(t.T(), antiAffinityGroupFixture.res.Description, antiAffinityGroup.Description)

		antiAffinityGroup, err = t.client.GetAntiAffinityGroupByName("lolnope")
		assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
		assert.Empty(t.T(), antiAffinityGroup)
	}, antiAffinityGroupFixtureOptName(antiAffinityGroupName))
}

func (t *antiAffinityGroupTestSuite) TestAntiAffinityGroupDelete() {
	t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
		antiAffinityGroup := t.client.antiAffinityGroupFromAPI(antiAffinityGroupFixture.res)

		if err := antiAffinityGroup.Delete(); err != nil {
			t.FailNow("SSH key deletion failed", err)
		}
		assert.Empty(t.T(), antiAffinityGroup.Name)
		assert.Empty(t.T(), antiAffinityGroup.Description)

		r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.AffinityGroup{ID: antiAffinityGroupFixture.res.ID})
		assert.Len(t.T(), r, 0)
	})
}

func TestAccComputeAntiAffinityGroupTestSuite(t *testing.T) {
	suite.Run(t, new(antiAffinityGroupTestSuite))
}
