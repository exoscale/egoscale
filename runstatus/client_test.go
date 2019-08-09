package runstatus

import (
	"context"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type clientTestSuite struct {
	suite.Suite
}

func (t *clientTestSuite) TestNewClient() {
	client, err := NewClient(context.Background(), "", "", nil)
	assert.EqualError(t.T(), err, egoerr.ErrMissingAPICredentials.Error())
	assert.Empty(t.T(), client)

	client, err = NewClient(
		context.Background(),
		"apiKey",
		"apiSecret",
		&ClientOpts{
			APIEndpoint: "apiEndpoint",
			Tracing:     true,
		})
	if err != nil {
		t.FailNow("client instantiation failed", err)
	}
	assert.NotEmpty(t.T(), client)
	assert.Equal(t.T(), "apiEndpoint", client.apiEndpoint)
	assert.True(t.T(), client.tracing)
}

func TestAccRunstatusClientTestSuite(t *testing.T) {
	suite.Run(t, new(clientTestSuite))
}
