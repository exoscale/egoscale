package v2

import (
	"context"
	"net/http"

	"github.com/exoscale/egoscale/v2/oapi"
	"github.com/stretchr/testify/mock"
)

var (
	testAPIKey       = new(testSuite).randomString(10)
	testAPIKeyName   = new(testSuite).randomString(10)
	testAPIKeyRoleID = new(testSuite).randomString(10)
	testAPIKeySecret = new(testSuite).randomString(10)
)

func (ts *testSuite) TestClient_GetAPIKey() {
	ts.mock().
		On("GetApiKeyWithResponse",
			mock.Anything,                 // ctx
			testAPIKey,                    // key
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetApiKeyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.IamApiKey{
				Key:    &testAPIKey,
				Name:   &testAPIKeyName,
				RoleId: &testAPIKeyRoleID,
			},
		}, nil)

	expected := &APIKey{
		Key:    &testAPIKey,
		Name:   &testAPIKeyName,
		RoleID: &testAPIKeyRoleID,
	}

	actual, err := ts.client.GetAPIKey(context.Background(), testZone, testAPIKey)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListAPIKeys() {
	ts.mock().
		On("ListApiKeysWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListApiKeysResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				ApiKeys *[]oapi.IamApiKey `json:"api-keys,omitempty"`
			}{
				ApiKeys: &[]oapi.IamApiKey{
					{
						Key:    &testAPIKey,
						Name:   &testAPIKeyName,
						RoleId: &testAPIKeyRoleID,
					},
				},
			},
		}, nil)

	expected := []*APIKey{
		&APIKey{
			Key:    &testAPIKey,
			Name:   &testAPIKeyName,
			RoleID: &testAPIKeyRoleID,
		},
	}

	actual, err := ts.client.ListAPIKeys(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateAPIKey() {
	ts.mock().
		On(
			"CreateApiKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateApiKeyJSONRequestBody{
					Name:   testAPIKeyName,
					RoleId: testAPIKeyRoleID,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateApiKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.IamApiKeyCreated{
					Key:    &testAPIKey,
					Name:   &testAPIKeyName,
					RoleId: &testAPIKeyRoleID,
					Secret: &testAPIKeySecret,
				},
			},
			nil,
		)

	ts.mock().
		On("GetApiKeyWithResponse",
			mock.Anything,                 // ctx
			testAPIKey,                    // key
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetApiKeyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.IamApiKey{
				Key:    &testAPIKey,
				Name:   &testAPIKeyName,
				RoleId: &testAPIKeyRoleID,
			},
		}, nil)

	expected := &APIKey{
		Key:    &testAPIKey,
		Name:   &testAPIKeyName,
		RoleID: &testAPIKeyRoleID,
	}

	actual, secret, err := ts.client.CreateAPIKey(context.Background(), testZone, &APIKey{
		Name:   &testAPIKeyName,
		RoleID: &testAPIKeyRoleID,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
	ts.Require().Equal(secret, testAPIKeySecret)
}

func (ts *testSuite) TestClient_DeleteAPIKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteApiKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // key
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testAPIKey, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteApiKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testAPIKey, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testAPIKey, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteAPIKey(
		context.Background(),
		testZone,
		&APIKey{Key: &testAPIKey},
	))
	ts.Require().True(deleted)
}
