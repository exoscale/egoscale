package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testIAMAccessKeyKey                  = new(testSuite).randomString(10)
	testIAMAccessKeyName                 = new(testSuite).randomString(10)
	testIAMAccessKeyOperationName        = new(testSuite).randomString(10)
	testIAMAccessKeyResourceDomain       = oapi.AccessKeyResourceDomainSos
	testIAMAccessKeyResourceName         = new(testSuite).randomString(10)
	testIAMAccessKeyResourceType         = oapi.AccessKeyResourceResourceTypeBucket
	testIAMAccessKeyRestrictedOperations = []string{
		new(testSuite).randomString(10),
		new(testSuite).randomString(10),
	}
	testIAMAccessKeyRestrictedTags = []string{
		new(testSuite).randomString(10),
		new(testSuite).randomString(10),
	}
	testIAMAccessKeySecret  = new(testSuite).randomString(10)
	testIAMAccessKeyType    = oapi.AccessKeyTypeRestricted
	testIAMAccessKeyVersion = oapi.AccessKeyVersionV2
)

func (ts *testSuite) TestClient_CreateIAMAccessKey() {
	ts.mock().
		On(
			"CreateAccessKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateAccessKeyJSONRequestBody{
					Name:       &testIAMAccessKeyName,
					Operations: &testIAMAccessKeyRestrictedOperations,
					Resources: &[]oapi.AccessKeyResource{{
						Domain:       &testIAMAccessKeyResourceDomain,
						ResourceName: &testIAMAccessKeyResourceName,
						ResourceType: &testIAMAccessKeyResourceType,
					}},
					Tags: &testIAMAccessKeyRestrictedTags,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateAccessKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.AccessKey{
					Key:        &testIAMAccessKeyKey,
					Name:       &testIAMAccessKeyName,
					Operations: &testIAMAccessKeyRestrictedOperations,
					Resources: &[]oapi.AccessKeyResource{{
						Domain:       &testIAMAccessKeyResourceDomain,
						ResourceName: &testIAMAccessKeyResourceName,
						ResourceType: &testIAMAccessKeyResourceType,
					}},
					Secret:  &testIAMAccessKeySecret,
					Tags:    &testIAMAccessKeyRestrictedTags,
					Type:    &testIAMAccessKeyType,
					Version: &testIAMAccessKeyVersion,
				},
			},
			nil,
		)

	expected := &IAMAccessKey{
		Key:        &testIAMAccessKeyKey,
		Name:       &testIAMAccessKeyName,
		Operations: &testIAMAccessKeyRestrictedOperations,
		Resources: &[]IAMAccessKeyResource{{
			Domain:       string(testIAMAccessKeyResourceDomain),
			ResourceName: testIAMAccessKeyResourceName,
			ResourceType: string(testIAMAccessKeyResourceType),
		}},
		Secret:  &testIAMAccessKeySecret,
		Tags:    &testIAMAccessKeyRestrictedTags,
		Type:    (*string)(&testIAMAccessKeyType),
		Version: (*string)(&testIAMAccessKeyVersion),
	}

	actual, err := ts.client.CreateIAMAccessKey(
		context.Background(),
		testZone,
		testIAMAccessKeyName,
		CreateIAMAccessKeyWithOperations(testIAMAccessKeyRestrictedOperations),
		CreateIAMAccessKeyWithTags(testIAMAccessKeyRestrictedTags),
		CreateIAMAccessKeyWithResources([]IAMAccessKeyResource{{
			Domain:       string(testIAMAccessKeyResourceDomain),
			ResourceName: testIAMAccessKeyResourceName,
			ResourceType: string(testIAMAccessKeyResourceType),
		}}),
	)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetIAMAccessKey() {
	ts.mock().
		On("GetAccessKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // key
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testIAMAccessKeyKey, args.Get(1))
		}).
		Return(&oapi.GetAccessKeyResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.AccessKey{
				Key:        &testIAMAccessKeyKey,
				Name:       &testIAMAccessKeyName,
				Operations: &testIAMAccessKeyRestrictedOperations,
				Resources: &[]oapi.AccessKeyResource{{
					Domain:       &testIAMAccessKeyResourceDomain,
					ResourceName: &testIAMAccessKeyResourceName,
					ResourceType: &testIAMAccessKeyResourceType,
				}},
				Tags:    &testIAMAccessKeyRestrictedTags,
				Type:    &testIAMAccessKeyType,
				Version: &testIAMAccessKeyVersion,
			},
		}, nil)

	expected := &IAMAccessKey{
		Key:        &testIAMAccessKeyKey,
		Name:       &testIAMAccessKeyName,
		Operations: &testIAMAccessKeyRestrictedOperations,
		Resources: &[]IAMAccessKeyResource{{
			Domain:       string(testIAMAccessKeyResourceDomain),
			ResourceName: testIAMAccessKeyResourceName,
			ResourceType: string(testIAMAccessKeyResourceType),
		}},
		Tags:    &testIAMAccessKeyRestrictedTags,
		Type:    (*string)(&testIAMAccessKeyType),
		Version: (*string)(&testIAMAccessKeyVersion),
	}

	actual, err := ts.client.GetIAMAccessKey(context.Background(), testZone, *expected.Key)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListIAMAccessKeyOperations() {
	ts.mock().
		On("ListAccessKeyKnownOperationsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListAccessKeyKnownOperationsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				AccessKeyOperations *[]oapi.AccessKeyOperation `json:"access-key-operations,omitempty"` // nolint:revive
			}{
				AccessKeyOperations: &[]oapi.AccessKeyOperation{{
					Operation: &testIAMAccessKeyOperationName,
					Tags:      &testIAMAccessKeyRestrictedTags,
				}},
			},
		}, nil)

	expected := []*IAMAccessKeyOperation{{
		Name: testIAMAccessKeyOperationName,
		Tags: testIAMAccessKeyRestrictedTags,
	}}

	actual, err := ts.client.ListIAMAccessKeyOperations(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListIAMAccessKeys() {
	ts.mock().
		On("ListAccessKeysWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListAccessKeysResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				AccessKeys *[]oapi.AccessKey `json:"access-keys,omitempty"` // nolint:revive
			}{
				AccessKeys: &[]oapi.AccessKey{{
					Key:        &testIAMAccessKeyKey,
					Name:       &testIAMAccessKeyName,
					Operations: &testIAMAccessKeyRestrictedOperations,
					Resources: &[]oapi.AccessKeyResource{{
						Domain:       &testIAMAccessKeyResourceDomain,
						ResourceName: &testIAMAccessKeyResourceName,
						ResourceType: &testIAMAccessKeyResourceType,
					}},
					Tags:    &testIAMAccessKeyRestrictedTags,
					Type:    &testIAMAccessKeyType,
					Version: &testIAMAccessKeyVersion,
				}},
			},
		}, nil)

	expected := []*IAMAccessKey{{
		Key:        &testIAMAccessKeyKey,
		Name:       &testIAMAccessKeyName,
		Operations: &testIAMAccessKeyRestrictedOperations,
		Resources: &[]IAMAccessKeyResource{{
			Domain:       string(testIAMAccessKeyResourceDomain),
			ResourceName: testIAMAccessKeyResourceName,
			ResourceType: string(testIAMAccessKeyResourceType),
		}},
		Tags:    &testIAMAccessKeyRestrictedTags,
		Type:    (*string)(&testIAMAccessKeyType),
		Version: (*string)(&testIAMAccessKeyVersion),
	}}

	actual, err := ts.client.ListIAMAccessKeys(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListMyIAMAccessKeyOperations() {
	ts.mock().
		On("ListAccessKeyOperationsWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListAccessKeyOperationsResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				AccessKeyOperations *[]oapi.AccessKeyOperation `json:"access-key-operations,omitempty"` // nolint:revive
			}{
				AccessKeyOperations: &[]oapi.AccessKeyOperation{{
					Operation: &testIAMAccessKeyOperationName,
					Tags:      &testIAMAccessKeyRestrictedTags,
				}},
			},
		}, nil)

	expected := []*IAMAccessKeyOperation{{
		Name: testIAMAccessKeyOperationName,
		Tags: testIAMAccessKeyRestrictedTags,
	}}

	actual, err := ts.client.ListMyIAMAccessKeyOperations(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RevokeIAMAccessKey() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"RevokeAccessKeyWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // key
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testIAMAccessKeyKey, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.RevokeAccessKeyResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testIAMAccessKeyKey, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testIAMAccessKeyName, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.RevokeIAMAccessKey(
		context.Background(),
		testZone,
		&IAMAccessKey{Key: &testIAMAccessKeyKey},
	))
	ts.Require().True(deleted)
}
