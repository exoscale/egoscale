package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

func (ts *testSuite) TestClient_StopPgDatabaseMigration() {
	var (
		testDatabaseName   = "testdb"
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		stopped            = false
	)

	ts.mock().
		On(
			"StopDbaasPgMigrationWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // name
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(oapi.DbaasServiceName(testDatabaseName), args.Get(1))
			stopped = true
		}).
		Return(
			&oapi.StopDbaasPgMigrationResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testDatabaseName, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testDatabaseName, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.StopPgDatabaseMigration(
		context.Background(),
		testZone,
		testDatabaseName,
	))
	ts.Require().True(stopped)
}
