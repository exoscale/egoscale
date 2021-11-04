package oapi

import (
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/suite"
)

type testSuite struct {
	suite.Suite

	client oapiClient
}

func (ts *testSuite) SetupTest() {
	ts.client = new(oapiClientMock)
}

func (ts *testSuite) TearDownTest() {
	ts.client = nil
}

func (ts *testSuite) randomID() string {
	id, err := uuid.NewV4()
	if err != nil {
		ts.T().Fatalf("unable to generate a new UUID: %s", err)
	}
	return id.String()
}

func TestSuiteOAPITestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}
