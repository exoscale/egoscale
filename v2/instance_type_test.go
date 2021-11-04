package v2

import (
	"context"
	"net/http"

	"github.com/stretchr/testify/mock"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstanceTypeAuthorized       = true
	testInstanceTypeCPUs       int64 = 16
	testInstanceTypeGPUs       int64 = 2
	testInstanceTypeFamily           = oapi.InstanceTypeFamilyGpu2
	testInstanceTypeID               = new(testSuite).randomID()
	testInstanceTypeMemory     int64 = 96636764160
	testInstanceTypeSize             = oapi.InstanceTypeSizeMedium
)

func (ts *testSuite) TestClient_FindInstanceType() {
	ts.mock().
		On("ListInstanceTypesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListInstanceTypesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				InstanceTypes *[]oapi.InstanceType `json:"instance-types,omitempty"`
			}{
				InstanceTypes: &[]oapi.InstanceType{{
					Authorized: &testInstanceTypeAuthorized,
					Cpus:       &testInstanceTypeCPUs,
					Family:     &testInstanceTypeFamily,
					Id:         &testInstanceTypeID,
					Memory:     &testInstanceTypeMemory,
					Size:       &testInstanceTypeSize,
				}},
			},
		}, nil)

	ts.mock().
		On("GetInstanceTypeWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceTypeID, args.Get(1))
		}).
		Return(&oapi.GetInstanceTypeResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.InstanceType{
				Authorized: &testInstanceTypeAuthorized,
				Cpus:       &testInstanceTypeCPUs,
				Family:     &testInstanceTypeFamily,
				Id:         &testInstanceTypeID,
				Memory:     &testInstanceTypeMemory,
				Size:       &testInstanceTypeSize,
			},
		}, nil)

	expected := &InstanceType{
		Authorized: &testInstanceTypeAuthorized,
		CPUs:       &testInstanceTypeCPUs,
		Family:     (*string)(&testInstanceTypeFamily),
		ID:         &testInstanceTypeID,
		Memory:     &testInstanceTypeMemory,
		Size:       (*string)(&testInstanceTypeSize),
	}

	actual, err := ts.client.FindInstanceType(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindInstanceType(context.Background(), testZone, *expected.Family+"."+*expected.Size)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_GetInstanceType() {
	ts.mock().
		On("GetInstanceTypeWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceTypeID, args.Get(1))
		}).
		Return(&oapi.GetInstanceTypeResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.InstanceType{
				Authorized: &testInstanceTypeAuthorized,
				Cpus:       &testInstanceTypeCPUs,
				Family:     &testInstanceTypeFamily,
				Gpus:       &testInstanceTypeGPUs,
				Id:         &testInstanceTypeID,
				Memory:     &testInstanceTypeMemory,
				Size:       &testInstanceTypeSize,
			},
		}, nil)

	expected := &InstanceType{
		Authorized: &testInstanceTypeAuthorized,
		CPUs:       &testInstanceTypeCPUs,
		Family:     (*string)(&testInstanceTypeFamily),
		GPUs:       &testInstanceTypeGPUs,
		ID:         &testInstanceTypeID,
		Memory:     &testInstanceTypeMemory,
		Size:       (*string)(&testInstanceTypeSize),
	}

	actual, err := ts.client.GetInstanceType(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListInstanceTypes() {
	ts.mock().
		On("ListInstanceTypesWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListInstanceTypesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				InstanceTypes *[]oapi.InstanceType `json:"instance-types,omitempty"`
			}{
				InstanceTypes: &[]oapi.InstanceType{{
					Authorized: &testInstanceTypeAuthorized,
					Cpus:       &testInstanceTypeCPUs,
					Family:     &testInstanceTypeFamily,
					Gpus:       &testInstanceTypeGPUs,
					Id:         &testInstanceTypeID,
					Memory:     &testInstanceTypeMemory,
					Size:       &testInstanceTypeSize,
				}},
			},
		}, nil)

	expected := []*InstanceType{{
		Authorized: &testInstanceTypeAuthorized,
		CPUs:       &testInstanceTypeCPUs,
		Family:     (*string)(&testInstanceTypeFamily),
		GPUs:       &testInstanceTypeGPUs,
		ID:         &testInstanceTypeID,
		Memory:     &testInstanceTypeMemory,
		Size:       (*string)(&testInstanceTypeSize),
	}}

	actual, err := ts.client.ListInstanceTypes(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}
