// +build testacc

package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type bucketTestSuite struct {
	suite.Suite
	client *Client
}

func (t *bucketTestSuite) SetupTest() {
	var err error

	if t.client, err = testClientFromEnv(); err != nil {
		t.FailNow("unable to initialize API client", err)
	}
}

// func (t *bucketTestSuite) TestCreateBucket() {
// }

func (t *bucketTestSuite) TestListBuckets() {
	_, teardown, err := bucketFixture("")
	if err != nil {
		t.FailNow("bucket fixture setup failed", err)
	}
	defer teardown() // nolint:errcheck

	// We cannot guarantee that there will be only our resources,
	// so we ensure we get at least our fixture bucket
	buckets, err := t.client.ListBuckets()
	if err != nil {
		t.FailNow("buckets listing failed", err)
	}
	assert.GreaterOrEqual(t.T(), len(buckets), 1)
}

// func (t *bucketTestSuite) TestGetBucket() {
// }

// func (t *bucketTestSuite) TestDeleteBucket() {
// }

func TestAccStorageBucketTestSuite(t *testing.T) {
	suite.Run(t, new(bucketTestSuite))
}
