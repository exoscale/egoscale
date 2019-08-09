// +build testacc

package runstatus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type pageTestSuite struct {
	suite.Suite
	client *Client
}

func (t *pageTestSuite) SetupTest() {
	var err error

	if t.client, err = testClientFromEnv(); err != nil {
		t.FailNow("unable to initialize API client", err)
	}
}

func (t *pageTestSuite) TestListPages() {
	_, teardown, err := pageFixture("")
	if err != nil {
		t.FailNow("page fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	// We cannot guarantee that there will be only our resources,
	// so we ensure we get at least our fixture page
	pages, err := t.client.ListPages()
	if err != nil {
		t.FailNow("pages listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(pages), 1)
}

func TestAccRunstatusPageTestSuite(t *testing.T) {
	suite.Run(t, new(pageTestSuite))
}
