package v2

import (
	"context"
	"net"
	"net/http"

	"github.com/stretchr/testify/mock"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testPrivateNetworkDescription     = new(testSuite).randomString(10)
	testPrivateNetworkEndIP           = "192.168.0.254"
	testPrivateNetworkEndIPP          = net.ParseIP(testPrivateNetworkEndIP)
	testPrivateNetworkID              = new(testSuite).randomID()
	testPrivateNetworkLabels          = map[string]string{"k1": "v1", "k2": "v2"}
	testPrivateNetworkName            = new(testSuite).randomString(10)
	testPrivateNetworkNetmask         = "255.255.255.0"
	testPrivateNetworkNetmaskP        = net.ParseIP(testPrivateNetworkNetmask)
	testPrivateNetworkStartIP         = "192.168.0.0"
	testPrivateNetworkStartIPP        = net.ParseIP(testPrivateNetworkStartIP)
	testPrivateNetworkLeaseInstanceID = new(testSuite).randomID()
	testPrivateNetworkLeaseIPAddress  = "192.168.0.1"
	testPrivateNetworkLeaseIPAddressP = net.ParseIP(testPrivateNetworkLeaseIPAddress)
)

func (ts *testSuite) TestClient_CreatePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreatePrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreatePrivateNetworkJSONRequestBody{
					Description: &testPrivateNetworkDescription,
					EndIp:       &testPrivateNetworkEndIP,
					Labels:      &oapi.Labels{AdditionalProperties: testPrivateNetworkLabels},
					Name:        testPrivateNetworkName,
					Netmask:     &testPrivateNetworkNetmask,
					StartIp:     &testPrivateNetworkStartIP,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreatePrivateNetworkResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetPrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetPrivateNetworkResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.PrivateNetwork{
				Description: &testPrivateNetworkDescription,
				EndIp:       &testPrivateNetworkEndIP,
				Id:          &testPrivateNetworkID,
				Labels:      &oapi.Labels{AdditionalProperties: testPrivateNetworkLabels},
				Name:        &testPrivateNetworkName,
				Netmask:     &testPrivateNetworkNetmask,
				StartIp:     &testPrivateNetworkStartIP,
			},
		}, nil)

	expected := &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Labels:      &testPrivateNetworkLabels,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmaskP,
		StartIP:     &testPrivateNetworkStartIPP,
		Zone:        &testZone,
	}

	actual, err := ts.client.CreatePrivateNetwork(context.Background(), testZone, &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Labels:      &testPrivateNetworkLabels,
		Name:        &testPrivateNetworkName,
		Netmask:     &testPrivateNetworkNetmaskP,
		StartIP:     &testPrivateNetworkStartIPP,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeletePrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeletePrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testPrivateNetworkID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeletePrivateNetworkResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeletePrivateNetwork(
		context.Background(),
		testZone,
		&PrivateNetwork{ID: &testPrivateNetworkID},
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_FindPrivateNetwork() {
	ts.mock().
		On("ListPrivateNetworksWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListPrivateNetworksResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				PrivateNetworks *[]oapi.PrivateNetwork `json:"private-networks,omitempty"`
			}{
				PrivateNetworks: &[]oapi.PrivateNetwork{
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
			},
		}, nil)

	ts.mock().
		On("GetPrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testPrivateNetworkID, args.Get(1))
		}).
		Return(&oapi.GetPrivateNetworkResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.PrivateNetwork{
				Id:   &testPrivateNetworkID,
				Name: &testPrivateNetworkName,
			},
		}, nil)

	expected := &PrivateNetwork{
		ID:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,
		Zone: &testZone,
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

func (ts *testSuite) TestClient_GetPrivateNetwork() {
	ts.mock().
		On("GetPrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testPrivateNetworkID, args.Get(1))
		}).
		Return(&oapi.GetPrivateNetworkResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.PrivateNetwork{
				Description: &testPrivateNetworkDescription,
				EndIp:       &testPrivateNetworkEndIP,
				Id:          &testPrivateNetworkID,
				Labels:      &oapi.Labels{AdditionalProperties: testPrivateNetworkLabels},
				Name:        &testPrivateNetworkName,
				Netmask:     &testPrivateNetworkNetmask,
				StartIp:     &testPrivateNetworkStartIP,
				Leases: &[]oapi.PrivateNetworkLease{{
					InstanceId: &testPrivateNetworkLeaseInstanceID,
					Ip:         &testPrivateNetworkLeaseIPAddress,
				}},
			},
		}, nil)

	expected := &PrivateNetwork{
		Description: &testPrivateNetworkDescription,
		EndIP:       &testPrivateNetworkEndIPP,
		ID:          &testPrivateNetworkID,
		Labels:      &testPrivateNetworkLabels,
		Leases: []*PrivateNetworkLease{{
			InstanceID: &testPrivateNetworkLeaseInstanceID,
			IPAddress:  &testPrivateNetworkLeaseIPAddressP,
		}},
		Name:    &testPrivateNetworkName,
		Netmask: &testPrivateNetworkNetmaskP,
		StartIP: &testPrivateNetworkStartIPP,
		Zone:    &testZone,
	}

	actual, err := ts.client.GetPrivateNetwork(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListPrivateNetworks() {
	ts.mock().
		On("ListPrivateNetworksWithResponse",
			mock.Anything,                 // ctx
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListPrivateNetworksResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				PrivateNetworks *[]oapi.PrivateNetwork `json:"private-networks,omitempty"`
			}{
				PrivateNetworks: &[]oapi.PrivateNetwork{{
					Description: &testPrivateNetworkDescription,
					EndIp:       &testPrivateNetworkEndIP,
					Id:          &testPrivateNetworkID,
					Leases: &[]oapi.PrivateNetworkLease{{
						InstanceId: &testPrivateNetworkLeaseInstanceID,
						Ip:         &testPrivateNetworkLeaseIPAddress,
					}},
					Name:    &testPrivateNetworkName,
					Netmask: &testPrivateNetworkNetmask,
					StartIp: &testPrivateNetworkStartIP,
				}},
			},
		}, nil)

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
		Zone:    &testZone,
	}}

	actual, err := ts.client.ListPrivateNetworks(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_UpdatePrivateNetwork() {
	var (
		testPrivateNetworkLabelsUpdated      = map[string]string{"k3": "v3"}
		testPrivateNetworkDescriptionUpdated = testPrivateNetworkDescription + "-updated"
		testPrivateNetworkEndIPUpdated       = "172.16.254.254"
		testPrivateNetworkEndIPPUpdated      = net.ParseIP(testPrivateNetworkEndIPUpdated)
		testPrivateNetworkNameUpdated        = testPrivateNetworkName + "-updated"
		testPrivateNetworkNetmaskUpdated     = "255.255.0.0"
		testPrivateNetworkNetmaskPUpdated    = net.ParseIP(testPrivateNetworkNetmaskUpdated)
		testPrivateNetworkStartIPUpdated     = "172.16.0.0"
		testPrivateNetworkStartIPPUpdated    = net.ParseIP(testPrivateNetworkStartIPUpdated)
		testOperationID                      = ts.randomID()
		testOperationState                   = oapi.OperationStateSuccess
		updated                              = false
	)

	ts.mock().
		On(
			"UpdatePrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdatePrivateNetworkJSONRequestBody{
					Description: &testPrivateNetworkDescriptionUpdated,
					EndIp:       &testPrivateNetworkEndIPUpdated,
					Name:        &testPrivateNetworkNameUpdated,
					Labels:      &oapi.Labels{AdditionalProperties: testPrivateNetworkLabelsUpdated},
					Netmask:     &testPrivateNetworkNetmaskUpdated,
					StartIp:     &testPrivateNetworkStartIPUpdated,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdatePrivateNetworkResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdatePrivateNetwork(context.Background(), testZone, &PrivateNetwork{
		Description: &testPrivateNetworkDescriptionUpdated,
		EndIP:       &testPrivateNetworkEndIPPUpdated,
		ID:          &testPrivateNetworkID,
		Labels:      &testPrivateNetworkLabelsUpdated,
		Name:        &testPrivateNetworkNameUpdated,
		Netmask:     &testPrivateNetworkNetmaskPUpdated,
		StartIP:     &testPrivateNetworkStartIPPUpdated,
	}))
	ts.Require().True(updated)
}

func (ts *testSuite) TestClient_UpdatePrivateNetworkInstanceIPAddress() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		updated            = false
	)

	ts.mock().
		On(
			"UpdatePrivateNetworkInstanceIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.UpdatePrivateNetworkInstanceIpJSONRequestBody{
					Instance: oapi.Instance{Id: &testPrivateNetworkLeaseInstanceID},
					Ip:       &testPrivateNetworkLeaseIPAddress,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdatePrivateNetworkInstanceIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testPrivateNetworkID, nil),
		State:     &testOperationState,
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
