package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/jarcoal/httpmock"

	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testPrivateNetworkDescription = "Test Private Network description"
	testPrivateNetworkEndIP       = "192.168.0.254"
	testPrivateNetworkID          = new(clientTestSuite).randomID()
	testPrivateNetworkName        = "test-private-network"
	testPrivateNetworkNetmask     = "255.255.255.0"
	testPrivateNetworkStartIP     = "192.168.0.0"
)

func (ts *clientTestSuite) TestClient_CreatePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	ts.mockAPIRequest("POST", "/private-network", papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testPrivateNetworkID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testPrivateNetworkID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIp:       &testPrivateNetworkEndIP,
		Id:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmask,
		StartIp:     &testPrivateNetworkStartIP,
	})

	expected := &PrivateNetwork{
		Description: testPrivateNetworkDescription,
		EndIP:       net.ParseIP(testPrivateNetworkEndIP),
		ID:          testPrivateNetworkID,
		Name:        testPrivateNetworkName,
		Netmask:     net.ParseIP(testPrivateNetworkNetmask),
		StartIP:     net.ParseIP(testPrivateNetworkStartIP),
	}

	actual, err := ts.client.CreatePrivateNetwork(context.Background(), testZone, &PrivateNetwork{
		Description: testPrivateNetworkDescription,
		EndIP:       net.ParseIP(testPrivateNetworkEndIP),
		ID:          testPrivateNetworkID,
		Name:        testPrivateNetworkName,
		Netmask:     net.ParseIP(testPrivateNetworkNetmask),
		StartIP:     net.ParseIP(testPrivateNetworkStartIP),
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_ListPrivateNetworks() {
	ts.mockAPIRequest("GET", "/private-network", struct {
		PrivateNetworks *[]papi.PrivateNetwork `json:"private-networks,omitempty"`
	}{
		PrivateNetworks: &[]papi.PrivateNetwork{{
			Description: &testPrivateNetworkDescription,
			EndIp:       &testPrivateNetworkEndIP,
			Id:          &testPrivateNetworkID,
			Name:        &testPrivateNetworkName,
			Netmask:     &testPrivateNetworkNetmask,
			StartIp:     &testPrivateNetworkStartIP,
		}},
	})

	expected := []*PrivateNetwork{{
		Description: testPrivateNetworkDescription,
		EndIP:       net.ParseIP(testPrivateNetworkEndIP),
		ID:          testPrivateNetworkID,
		Name:        testPrivateNetworkName,
		Netmask:     net.ParseIP(testPrivateNetworkNetmask),
		StartIP:     net.ParseIP(testPrivateNetworkStartIP),
	}}

	actual, err := ts.client.ListPrivateNetworks(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_GetPrivateNetwork() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIp:       &testPrivateNetworkEndIP,
		Id:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmask,
		StartIp:     &testPrivateNetworkStartIP,
	})

	expected := &PrivateNetwork{
		Description: testPrivateNetworkDescription,
		EndIP:       net.ParseIP(testPrivateNetworkEndIP),
		ID:          testPrivateNetworkID,
		Name:        testPrivateNetworkName,
		Netmask:     net.ParseIP(testPrivateNetworkNetmask),
		StartIP:     net.ParseIP(testPrivateNetworkStartIP),
	}

	actual, err := ts.client.GetPrivateNetwork(context.Background(), testZone, expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdatePrivateNetwork() {
	var (
		testPrivateNetworkDescriptionUpdated = testPrivateNetworkDescription + "-updated"
		testPrivateNetworkEndIPUpdated       = "172.16.254.254"
		testPrivateNetworkNameUpdated        = testPrivateNetworkName + "-updated"
		testPrivateNetworkNetmaskUpdated     = "255.255.0.0"
		testPrivateNetworkStartIPUpdated     = "172.16.0.0"
		testOperationID                      = ts.randomID()
		testOperationState                   = "success"
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/private-network/%s", testPrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			var actual papi.UpdatePrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdatePrivateNetworkJSONRequestBody{
				Description: &testPrivateNetworkDescriptionUpdated,
				EndIp:       &testPrivateNetworkEndIPUpdated,
				Name:        &testPrivateNetworkNameUpdated,
				Netmask:     &testPrivateNetworkNetmaskUpdated,
				StartIp:     &testPrivateNetworkStartIPUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testPrivateNetworkID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testPrivateNetworkID},
	})

	ts.Require().NoError(ts.client.UpdatePrivateNetwork(context.Background(), testZone, &PrivateNetwork{
		Description: testPrivateNetworkDescriptionUpdated,
		EndIP:       net.ParseIP(testPrivateNetworkEndIPUpdated),
		ID:          testPrivateNetworkID,
		Name:        testPrivateNetworkNameUpdated,
		Netmask:     net.ParseIP(testPrivateNetworkNetmaskUpdated),
		StartIP:     net.ParseIP(testPrivateNetworkStartIPUpdated),
	}))
}

func (ts *clientTestSuite) TestClient_DeletePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = "success"
	)

	httpmock.RegisterResponder("DELETE", "=~^/private-network/.*",
		func(req *http.Request) (*http.Response, error) {
			ts.Require().Equal(fmt.Sprintf("/private-network/%s", testPrivateNetworkID), req.URL.String())

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testPrivateNetworkID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testPrivateNetworkID},
	})

	ts.Require().NoError(ts.client.DeletePrivateNetwork(context.Background(), testZone, testPrivateNetworkID))
}
