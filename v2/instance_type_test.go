package v2

import (
	"context"
	"fmt"

	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstanceTypeAuthorized       = true
	testInstanceTypeCPUs       int64 = 16
	testInstanceTypeGPUs       int64 = 2
	testInstanceTypeFamily           = oapi.InstanceTypeFamilyGpu2
	testInstanceTypeID               = new(clientTestSuite).randomID()
	testInstanceTypeMemory     int64 = 96636764160
	testInstanceTypeSize             = oapi.InstanceTypeSizeMedium
)

func (ts *clientTestSuite) TestClient_ListInstanceTypes() {
	ts.mockAPIRequest("GET", "/instance-type", struct {
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
	})

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

func (ts *clientTestSuite) TestClient_GetInstanceType() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-type/%s", testInstanceTypeID), oapi.InstanceType{
		Authorized: &testInstanceTypeAuthorized,
		Cpus:       &testInstanceTypeCPUs,
		Family:     &testInstanceTypeFamily,
		Gpus:       &testInstanceTypeGPUs,
		Id:         &testInstanceTypeID,
		Memory:     &testInstanceTypeMemory,
		Size:       &testInstanceTypeSize,
	})

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

func (ts *clientTestSuite) TestClient_FindInstanceType() {
	ts.mockAPIRequest("GET", "/instance-type", struct {
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
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance-type/%s", testInstanceTypeID), oapi.InstanceType{
		Authorized: &testInstanceTypeAuthorized,
		Cpus:       &testInstanceTypeCPUs,
		Family:     &testInstanceTypeFamily,
		Id:         &testInstanceTypeID,
		Memory:     &testInstanceTypeMemory,
		Size:       &testInstanceTypeSize,
	})

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
