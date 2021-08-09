package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/jarcoal/httpmock"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

var (
	testPrivateNetworkDescription     = new(clientTestSuite).randomString(10)
	testPrivateNetworkEndIP           = "192.168.0.254"
	testPrivateNetworkEndIPP          = net.ParseIP(testPrivateNetworkEndIP)
	testPrivateNetworkID              = new(clientTestSuite).randomID()
	testPrivateNetworkName            = new(clientTestSuite).randomString(10)
	testPrivateNetworkNetmask         = "255.255.255.0"
	testPrivateNetworkNetmaskP        = net.ParseIP(testPrivateNetworkNetmask)
	testPrivateNetworkStartIP         = "192.168.0.0"
	testPrivateNetworkStartIPP        = net.ParseIP(testPrivateNetworkStartIP)
	testPrivateNetworkLeaseInstanceID = new(clientTestSuite).randomID()
	testPrivateNetworkLeaseIPAddress  = "192.168.0.1"
	testPrivateNetworkLeaseIPAddressP = net.ParseIP(testPrivateNetworkLeaseIPAddress)
)

func (ts *clientTestSuite) TestClient_CreatePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/private-network",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreatePrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreatePrivateNetworkJSONRequestBody{
				Description: &testPrivateNetworkDescription,
				EndIp:       &testPrivateNetworkEndIP,
				Name:        testPrivateNetworkName,
				Netmask:     &testPrivateNetworkNetmask,
				StartIp:     &testPrivateNetworkStartIP,
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

	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIp:       &testPrivateNetworkEndIP,
		Id:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmask,
		StartIp:     &testPrivateNetworkStartIP,
	})

	expected := &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmaskP,
		StartIP:     &testPrivateNetworkStartIPP,
	}

	actual, err := ts.client.CreatePrivateNetwork(context.Background(), testZone, &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmaskP,
		StartIP:     &testPrivateNetworkStartIPP,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_DeletePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/private-network/%s", testPrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

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

	ts.Require().NoError(ts.client.DeletePrivateNetwork(
		context.Background(),
		testZone,
		&PrivateNetwork{ID: &testPrivateNetworkID},
	))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_FindPrivateNetwork() {
	ts.mockAPIRequest("GET", "/private-network", struct {
		PrivateNetworks *[]papi.PrivateNetwork `json:"private-networks,omitempty"`
	}{
		PrivateNetworks: &[]papi.PrivateNetwork{
			{
				Id:   &testPrivateNetworkID,
				Name: &testPrivateNetworkName,
			},
			{
				Id:   func() *string { id := ts.randomID(); return &id }(),
				Name: func() *string { name := "dup"; return &name }(),
			},
			{
				Id:   func() *string { id := ts.randomID(); return &id }(),
				Name: func() *string { name := "dup"; return &name }(),
			},
		},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Id:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,
	})

	expected := &PrivateNetwork{
		ID:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,
	}

	actual, err := ts.client.FindPrivateNetwork(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindPrivateNetwork(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	_, err = ts.client.FindPrivateNetwork(context.Background(), testZone, "dup")
	ts.Require().EqualError(err, apiv2.ErrTooManyFound.Error())
}

func (ts *clientTestSuite) TestClient_GetPrivateNetwork() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/private-network/%s", testPrivateNetworkID), papi.PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIp:       &testPrivateNetworkEndIP,
		Id:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmask,
		StartIp:     &testPrivateNetworkStartIP,
		Leases: &[]papi.PrivateNetworkLease{{
			InstanceId: &testPrivateNetworkLeaseInstanceID,
			Ip:         &testPrivateNetworkLeaseIPAddress,
		}},
	})

	expected := &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Leases: []*PrivateNetworkLease{{
			InstanceID: &testPrivateNetworkLeaseInstanceID,
			IPAddress:  &testPrivateNetworkLeaseIPAddressP,
		}},
		Name:    &testPrivateNetworkName,
		Netmask: &testPrivateNetworkNetmaskP,
		StartIP: &testPrivateNetworkStartIPP,
	}

	actual, err := ts.client.GetPrivateNetwork(context.Background(), testZone, *expected.ID)
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
			Leases: &[]papi.PrivateNetworkLease{{
				InstanceId: &testPrivateNetworkLeaseInstanceID,
				Ip:         &testPrivateNetworkLeaseIPAddress,
			}},
			Name:    &testPrivateNetworkName,
			Netmask: &testPrivateNetworkNetmask,
			StartIp: &testPrivateNetworkStartIP,
		}},
	})

	expected := []*PrivateNetwork{{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Leases: []*PrivateNetworkLease{{
			InstanceID: &testPrivateNetworkLeaseInstanceID,
			IPAddress:  &testPrivateNetworkLeaseIPAddressP,
		}},
		Name:    &testPrivateNetworkName,
		Netmask: &testPrivateNetworkNetmaskP,
		StartIP: &testPrivateNetworkStartIPP,
	}}

	actual, err := ts.client.ListPrivateNetworks(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_UpdatePrivateNetwork() {
	var (
		testPrivateNetworkDescriptionUpdated = testPrivateNetworkDescription + "-updated"
		testPrivateNetworkEndIPUpdated       = "172.16.254.254"
		testPrivateNetworkEndIPPUpdated      = net.ParseIP(testPrivateNetworkEndIPUpdated)
		testPrivateNetworkNameUpdated        = testPrivateNetworkName + "-updated"
		testPrivateNetworkNetmaskUpdated     = "255.255.0.0"
		testPrivateNetworkNetmaskPUpdated    = net.ParseIP(testPrivateNetworkNetmaskUpdated)
		testPrivateNetworkStartIPUpdated     = "172.16.0.0"
		testPrivateNetworkStartIPPUpdated    = net.ParseIP(testPrivateNetworkStartIPUpdated)
		testOperationID                      = ts.randomID()
		testOperationState                   = papi.OperationStateSuccess
		updated                              = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/private-network/%s", testPrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

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
		Description: &testPrivateNetworkDescriptionUpdated,
		EndIP:       &testPrivateNetworkEndIPPUpdated,
		ID:          &testPrivateNetworkID,
		Name:        &testPrivateNetworkNameUpdated,
		Netmask:     &testPrivateNetworkNetmaskPUpdated,
		StartIP:     &testPrivateNetworkStartIPPUpdated,
	}))
	ts.Require().True(updated)
}

func (ts *clientTestSuite) TestClient_UpdatePrivateNetworkInstanceIPAddress() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		updated            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/private-network/%s:update-ip", testPrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdatePrivateNetworkInstanceIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdatePrivateNetworkInstanceIpJSONRequestBody{
				Instance: papi.Instance{Id: &testPrivateNetworkLeaseInstanceID},
				Ip:       &testPrivateNetworkLeaseIPAddress,
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

	ts.Require().NoError(ts.client.UpdatePrivateNetworkInstanceIPAddress(
		context.Background(),
		testZone,
		&Instance{ID: &testPrivateNetworkLeaseInstanceID},
		&PrivateNetwork{ID: &testPrivateNetworkID},
		testPrivateNetworkLeaseIPAddressP),
	)
	ts.Require().True(updated)
}
