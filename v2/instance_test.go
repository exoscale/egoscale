package v2

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/stretchr/testify/mock"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstanceAntiAffinityGroupID       = new(testSuite).randomID()
	testInstanceCreatedAt, _              = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testInstanceDiskSize            int64 = 10
	testInstanceElasticIPID               = new(testSuite).randomID()
	testInstanceID                        = new(testSuite).randomID()
	testInstanceIPv6Address               = "2001:db8:abcd::1"
	testInstanceIPv6AddressP              = net.ParseIP(testInstanceIPv6Address)
	testInstanceIPv6Enabled               = true
	testInstanceInstanceTypeID            = new(testSuite).randomID()
	testInstanceLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testInstanceManagerID                 = new(testSuite).randomID()
	testInstanceManagerType               = oapi.ManagerTypeInstancePool
	testInstanceName                      = new(testSuite).randomString(10)
	testInstancePrivateNetworkID          = new(testSuite).randomID()
	testInstancePublicIP                  = "1.2.3.4"
	testInstancePublicIPP                 = net.ParseIP(testInstancePublicIP)
	testInstanceSSHKey                    = new(testSuite).randomString(10)
	testInstanceSecurityGroupID           = new(testSuite).randomID()
	testInstanceSnapshotID                = new(testSuite).randomID()
	testInstanceStartRescueProfile        = new(testSuite).randomString(10)
	testInstanceState                     = oapi.InstanceStateRunning
	testInstanceTemplateID                = new(testSuite).randomID()
	testInstanceUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
	testInstanceReverseDNSDomain = "example.net"
)

func (ts *testSuite) TestClient_AttachInstanceToElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	ts.mock().
		On(
			"AttachInstanceToElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.AttachInstanceToElasticIpJSONRequestBody{Instance: oapi.Instance{Id: &testInstanceID}},
				args.Get(2),
			)
			attached = true
		}).
		Return(
			&oapi.AttachInstanceToElasticIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.AttachInstanceToElasticIP(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&ElasticIP{ID: &testInstanceElasticIPID},
	))
	ts.Require().True(attached)
}

func (ts *testSuite) TestClient_AttachInstanceToPrivateNetwork() {
	var (
		testOperationID      = ts.randomID()
		testOperationState   = oapi.OperationStateSuccess
		testPrivateIPAddress = net.ParseIP("10.0.0.1")
		attached             = false
	)

	ts.mock().
		On(
			"AttachInstanceToPrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.AttachInstanceToPrivateNetworkJSONRequestBody{
					Instance: oapi.Instance{Id: &testInstanceID},
					Ip:       func() *string { ip := testPrivateIPAddress.String(); return &ip }(),
				},
				args.Get(2),
			)
			attached = true
		}).
		Return(
			&oapi.AttachInstanceToPrivateNetworkResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.AttachInstanceToPrivateNetwork(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&PrivateNetwork{ID: &testInstancePrivateNetworkID},
		AttachInstanceToPrivateNetworkWithIPAddress(testPrivateIPAddress),
	))
	ts.Require().True(attached)
}

func (ts *testSuite) TestClient_AttachInstanceToSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	ts.mock().
		On(
			"AttachInstanceToSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.AttachInstanceToSecurityGroupJSONRequestBody{Instance: oapi.Instance{Id: &testInstanceID}},
				args.Get(2),
			)
			attached = true
		}).
		Return(
			&oapi.AttachInstanceToSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.AttachInstanceToSecurityGroup(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&SecurityGroup{ID: &testInstanceSecurityGroupID},
	))
	ts.Require().True(attached)
}

func (ts *testSuite) TestClient_CreateInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.CreateInstanceJSONRequestBody{
					AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
					DiskSize:           testInstanceDiskSize,
					InstanceType:       oapi.InstanceType{Id: &testInstanceInstanceTypeID},
					Ipv6Enabled:        &testInstanceIPv6Enabled,
					Labels:             &oapi.Labels{AdditionalProperties: testInstanceLabels},
					Name:               &testInstanceName,
					SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
					SshKey:             &oapi.SshKey{Name: &testInstanceSSHKey},
					Template:           oapi.Template{Id: &testInstanceTemplateID},
					UserData:           &testInstanceUserData,
				},
				args.Get(1),
			)
		}).
		Return(
			&oapi.CreateInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetInstanceResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Instance{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
				CreatedAt:          &testInstanceCreatedAt,
				DiskSize:           &testInstanceDiskSize,
				ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstanceElasticIPID}},
				Id:                 &testInstanceID,
				InstanceType:       &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
				Ipv6Address:        &testInstanceIPv6Address,
				Labels:             &oapi.Labels{AdditionalProperties: testInstanceLabels},
				Manager:            &oapi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
				Name:               &testInstanceName,
				PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
				PublicIp:           &testInstancePublicIP,
				SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
				Snapshots:          &[]oapi.Snapshot{{Id: &testInstanceSnapshotID}},
				SshKey:             &oapi.SshKey{Name: &testInstanceSSHKey},
				State:              &testInstanceState,
				Template:           &oapi.Template{Id: &testInstanceTemplateID},
				UserData:           &testInstanceUserData,
			},
		}, nil)

	expected := &Instance{
		AntiAffinityGroupIDs: &[]string{testInstanceAntiAffinityGroupID},
		CreatedAt:            &testInstanceCreatedAt,
		DiskSize:             &testInstanceDiskSize,
		ElasticIPIDs:         &[]string{testInstanceElasticIPID},
		ID:                   &testInstanceID,
		IPv6Address:          &testInstanceIPv6AddressP,
		IPv6Enabled:          &testInstanceIPv6Enabled,
		InstanceTypeID:       &testInstanceInstanceTypeID,
		Labels:               &testInstanceLabels,
		Manager:              &InstanceManager{ID: testInstanceManagerID, Type: string(testInstanceManagerType)},
		Name:                 &testInstanceName,
		PrivateNetworkIDs:    &[]string{testInstancePrivateNetworkID},
		PublicIPAddress:      &testInstancePublicIPP,
		SSHKey:               &testInstanceSSHKey,
		SecurityGroupIDs:     &[]string{testInstanceSecurityGroupID},
		SnapshotIDs:          &[]string{testInstanceSnapshotID},
		State:                (*string)(&testInstanceState),
		TemplateID:           &testInstanceTemplateID,
		UserData:             &testInstanceUserData,
		Zone:                 &testZone,
	}

	actual, err := ts.client.CreateInstance(context.Background(), testZone, &Instance{
		AntiAffinityGroupIDs: &[]string{testInstanceAntiAffinityGroupID},
		CreatedAt:            &testInstanceCreatedAt,
		DiskSize:             &testInstanceDiskSize,
		ElasticIPIDs:         &[]string{testInstanceElasticIPID},
		ID:                   &testInstanceID,
		IPv6Address:          &testInstanceIPv6AddressP,
		IPv6Enabled:          &testInstanceIPv6Enabled,
		InstanceTypeID:       &testInstanceInstanceTypeID,
		Labels:               &testInstanceLabels,
		Manager:              &InstanceManager{ID: testInstanceManagerID, Type: string(testInstanceManagerType)},
		Name:                 &testInstanceName,
		PrivateNetworkIDs:    &[]string{testInstancePrivateNetworkID},
		PublicIPAddress:      &testInstancePublicIPP,
		SSHKey:               &testInstanceSSHKey,
		SecurityGroupIDs:     &[]string{testInstanceSecurityGroupID},
		SnapshotIDs:          &[]string{testInstanceSnapshotID},
		State:                (*string)(&testInstanceState),
		TemplateID:           &testInstanceTemplateID,
		UserData:             &testInstanceUserData,
	})
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_CreateInstanceSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mock().
		On(
			"CreateSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
		}).
		Return(
			&oapi.CreateSnapshotResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testSnapshotID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testSnapshotID, nil),
		State:     &testOperationState,
	})

	ts.mock().
		On("GetSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.GetSnapshotResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Snapshot{
				CreatedAt: &testSnapshotCreatedAt,
				Id:        &testSnapshotID,
				Instance:  &oapi.Instance{Id: &testInstanceID},
				Name:      &testSnapshotName,
				State:     &testSnapshotState,
			},
		}, nil)

	expected := &Snapshot{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testInstanceID,
		Name:       &testSnapshotName,
		State:      (*string)(&testSnapshotState),
		Zone:       &testZone,
	}

	actual, err := ts.client.CreateInstanceSnapshot(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
	)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On(
			"DeleteInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_DetachInstanceFromElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		detached           = false
	)

	ts.mock().
		On(
			"DetachInstanceFromElasticIpWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.DetachInstanceFromElasticIpJSONRequestBody{Instance: oapi.Instance{Id: &testInstanceID}},
				args.Get(2),
			)
			detached = true
		}).
		Return(
			&oapi.DetachInstanceFromElasticIpResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DetachInstanceFromElasticIP(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&ElasticIP{ID: &testInstanceElasticIPID},
	))
	ts.Require().True(detached)
}

func (ts *testSuite) TestClient_DetachInstanceFromPrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		detached           = false
	)

	ts.mock().
		On(
			"DetachInstanceFromPrivateNetworkWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.DetachInstanceFromPrivateNetworkJSONRequestBody{Instance: oapi.Instance{Id: &testInstanceID}},
				args.Get(2),
			)
			detached = true
		}).
		Return(
			&oapi.DetachInstanceFromPrivateNetworkResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DetachInstanceFromPrivateNetwork(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&PrivateNetwork{ID: &testInstancePrivateNetworkID},
	))
	ts.Require().True(detached)
}

func (ts *testSuite) TestClient_DetachInstanceFromSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		detached           = false
	)

	ts.mock().
		On(
			"DetachInstanceFromSecurityGroupWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				oapi.DetachInstanceFromSecurityGroupJSONRequestBody{Instance: oapi.Instance{Id: &testInstanceID}},
				args.Get(2),
			)
			detached = true
		}).
		Return(
			&oapi.DetachInstanceFromSecurityGroupResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DetachInstanceFromSecurityGroup(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&SecurityGroup{ID: &testInstanceSecurityGroupID},
	))
	ts.Require().True(detached)
}

func (ts *testSuite) TestClient_FindInstance() {
	ts.mock().
		On("ListInstancesWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // params
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Return(&oapi.ListInstancesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Instances *[]oapi.Instance `json:"instances,omitempty"`
			}{
				Instances: &[]oapi.Instance{
					{
						CreatedAt:    &testInstanceCreatedAt,
						DiskSize:     &testInstanceDiskSize,
						Id:           &testInstanceID,
						InstanceType: &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
						Name:         &testInstanceName,
						State:        &testInstanceState,
						Template:     &oapi.Template{Id: &testInstanceTemplateID},
					},
					{
						CreatedAt:    &testInstanceCreatedAt,
						DiskSize:     &testInstanceDiskSize,
						Id:           func() *string { id := ts.randomID(); return &id }(),
						InstanceType: &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
						Name:         func() *string { name := "dup"; return &name }(),
						State:        &testInstanceState,
						Template:     &oapi.Template{Id: &testInstanceTemplateID},
					},
					{
						CreatedAt:    &testInstanceCreatedAt,
						DiskSize:     &testInstanceDiskSize,
						Id:           func() *string { id := ts.randomID(); return &id }(),
						InstanceType: &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
						Name:         func() *string { name := "dup"; return &name }(),
						State:        &testInstanceState,
						Template:     &oapi.Template{Id: &testInstanceTemplateID},
					},
				},
			},
		}, nil)

	ts.mock().
		On("GetInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
		}).
		Return(&oapi.GetInstanceResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Instance{
				CreatedAt:    &testInstanceCreatedAt,
				DiskSize:     &testInstanceDiskSize,
				Id:           &testInstanceID,
				InstanceType: &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
				Name:         &testInstanceName,
				State:        &testInstanceState,
				Template:     &oapi.Template{Id: &testInstanceTemplateID},
			},
		}, nil)

	expected := &Instance{
		CreatedAt:      &testInstanceCreatedAt,
		DiskSize:       &testInstanceDiskSize,
		ID:             &testInstanceID,
		InstanceTypeID: &testInstanceInstanceTypeID,
		Name:           &testInstanceName,
		State:          (*string)(&testInstanceState),
		TemplateID:     &testInstanceTemplateID,
		Zone:           &testZone,
	}

	actual, err := ts.client.FindInstance(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	actual, err = ts.client.FindInstance(context.Background(), testZone, *expected.Name)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)

	_, err = ts.client.FindInstance(context.Background(), testZone, "dup")
	ts.Require().EqualError(err, apiv2.ErrTooManyFound.Error())
}

func (ts *testSuite) TestClient_GetInstance() {
	ts.mock().
		On("GetInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
		}).
		Return(&oapi.GetInstanceResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.Instance{
				AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
				CreatedAt:          &testInstanceCreatedAt,
				DiskSize:           &testInstanceDiskSize,
				ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstanceElasticIPID}},
				Id:                 &testInstanceID,
				InstanceType:       &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
				Ipv6Address:        &testInstanceIPv6Address,
				Manager:            &oapi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
				Name:               &testInstanceName,
				PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
				PublicIp:           &testInstancePublicIP,
				SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
				Snapshots:          &[]oapi.Snapshot{{Id: &testInstanceSnapshotID}},
				SshKey:             &oapi.SshKey{Name: &testInstanceSSHKey},
				State:              &testInstanceState,
				Template:           &oapi.Template{Id: &testInstanceTemplateID},
				UserData:           &testInstanceUserData,
			},
		}, nil)

	expected := &Instance{
		AntiAffinityGroupIDs: &[]string{testInstanceAntiAffinityGroupID},
		CreatedAt:            &testInstanceCreatedAt,
		DiskSize:             &testInstanceDiskSize,
		ElasticIPIDs:         &[]string{testInstanceElasticIPID},
		ID:                   &testInstanceID,
		IPv6Address:          &testInstanceIPv6AddressP,
		IPv6Enabled:          &testInstanceIPv6Enabled,
		InstanceTypeID:       &testInstanceInstanceTypeID,
		Manager:              &InstanceManager{ID: testInstanceManagerID, Type: string(testInstanceManagerType)},
		Name:                 &testInstanceName,
		PrivateNetworkIDs:    &[]string{testInstancePrivateNetworkID},
		PublicIPAddress:      &testInstancePublicIPP,
		SSHKey:               &testInstanceSSHKey,
		SecurityGroupIDs:     &[]string{testInstanceSecurityGroupID},
		SnapshotIDs:          &[]string{testInstanceSnapshotID},
		State:                (*string)(&testInstanceState),
		TemplateID:           &testInstanceTemplateID,
		UserData:             &testInstanceUserData,
		Zone:                 &testZone,
	}

	actual, err := ts.client.GetInstance(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_ListInstances() {
	ts.mock().
		On("ListInstancesWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // params
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(
				&oapi.ListInstancesParams{
					ManagerId:   &testInstanceManagerID,
					ManagerType: (*oapi.ListInstancesParamsManagerType)(&testInstanceManagerType),
				},
				args.Get(1),
			)
		}).
		Return(&oapi.ListInstancesResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &struct {
				Instances *[]oapi.Instance `json:"instances,omitempty"`
			}{
				Instances: &[]oapi.Instance{{
					AntiAffinityGroups: &[]oapi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
					CreatedAt:          &testInstanceCreatedAt,
					DiskSize:           &testInstanceDiskSize,
					ElasticIps:         &[]oapi.ElasticIp{{Id: &testInstanceElasticIPID}},
					Id:                 &testInstanceID,
					InstanceType:       &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
					Ipv6Address:        &testInstanceIPv6Address,
					Manager:            &oapi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
					Name:               &testInstanceName,
					PrivateNetworks:    &[]oapi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
					PublicIp:           &testInstancePublicIP,
					SecurityGroups:     &[]oapi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
					Snapshots:          &[]oapi.Snapshot{{Id: &testInstanceSnapshotID}},
					SshKey:             &oapi.SshKey{Name: &testInstanceSSHKey},
					State:              &testInstanceState,
					Template:           &oapi.Template{Id: &testInstanceTemplateID},
					UserData:           &testInstanceUserData,
				}},
			},
		}, nil)

	expected := []*Instance{{
		AntiAffinityGroupIDs: &[]string{testInstanceAntiAffinityGroupID},
		CreatedAt:            &testInstanceCreatedAt,
		DiskSize:             &testInstanceDiskSize,
		ElasticIPIDs:         &[]string{testInstanceElasticIPID},
		ID:                   &testInstanceID,
		IPv6Address:          &testInstanceIPv6AddressP,
		IPv6Enabled:          &testInstanceIPv6Enabled,
		InstanceTypeID:       &testInstanceInstanceTypeID,
		Manager:              &InstanceManager{ID: testInstanceManagerID, Type: string(testInstanceManagerType)},
		Name:                 &testInstanceName,
		PrivateNetworkIDs:    &[]string{testInstancePrivateNetworkID},
		PublicIPAddress:      &testInstancePublicIPP,
		SSHKey:               &testInstanceSSHKey,
		SecurityGroupIDs:     &[]string{testInstanceSecurityGroupID},
		SnapshotIDs:          &[]string{testInstanceSnapshotID},
		State:                (*string)(&testInstanceState),
		TemplateID:           &testInstanceTemplateID,
		UserData:             &testInstanceUserData,
		Zone:                 &testZone,
	}}

	actual, err := ts.client.ListInstances(
		context.Background(),
		testZone,
		ListInstancesByManagerID(testInstanceManagerID),
		ListInstancesByManagerType(string(testInstanceManagerType)),
	)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_RebootInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		rebooted           = false
	)

	ts.mock().
		On(
			"RebootInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			rebooted = true
		}).
		Return(
			&oapi.RebootInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.RebootInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(rebooted)
}

func (ts *testSuite) TestClient_ResetInstance() {
	var (
		testResetDiskSize   int64 = 50
		testResetTemplateID       = ts.randomID()
		testOperationID           = ts.randomID()
		testOperationState        = oapi.OperationStateSuccess
		reset                     = false
	)

	ts.mock().
		On(
			"ResetInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(oapi.ResetInstanceJSONRequestBody{
				DiskSize: &testResetDiskSize,
				Template: &oapi.Template{Id: &testResetTemplateID},
			}, args.Get(2))
			reset = true
		}).
		Return(
			&oapi.ResetInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.ResetInstance(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		ResetInstanceWithTemplate(&Template{ID: &testResetTemplateID}),
		ResetInstanceWithDiskSize(testResetDiskSize),
	))
	ts.Require().True(reset)
}

func (ts *testSuite) TestClient_ResizeInstanceDisk() {
	var (
		testResizeDiskSize int64 = 50
		testOperationID          = ts.randomID()
		testOperationState       = oapi.OperationStateSuccess
		resized                  = false
	)

	ts.mock().
		On(
			"ResizeInstanceDiskWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(oapi.ResizeInstanceDiskJSONRequestBody{DiskSize: testResizeDiskSize}, args.Get(2))
			resized = true
		}).
		Return(
			&oapi.ResizeInstanceDiskResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.ResizeInstanceDisk(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		testResizeDiskSize),
	)
	ts.Require().True(resized)
}

func (ts *testSuite) TestClient_RevertInstanceToSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		reverted           = false
	)

	ts.mock().
		On(
			"RevertInstanceToSnapshotWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(oapi.RevertInstanceToSnapshotJSONRequestBody{Id: testSnapshotID}, args.Get(2))
			reverted = true
		}).
		Return(
			&oapi.RevertInstanceToSnapshotResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	snapshot := &Snapshot{
		ID:         &testSnapshotID,
		InstanceID: &testInstanceID,
	}

	ts.Require().NoError(ts.client.RevertInstanceToSnapshot(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		snapshot,
	))
	ts.Require().True(reverted)
}

func (ts *testSuite) TestClient_ScaleInstance() {
	var (
		testScaleInstanceTypeID = ts.randomID()
		testOperationID         = ts.randomID()
		testOperationState      = oapi.OperationStateSuccess
		scaled                  = false
	)

	ts.mock().
		On(
			"ScaleInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(oapi.ScaleInstanceJSONRequestBody{
				InstanceType: oapi.InstanceType{Id: &testScaleInstanceTypeID},
			}, args.Get(2))
			scaled = true
		}).
		Return(
			&oapi.ScaleInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.ScaleInstance(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&InstanceType{ID: &testScaleInstanceTypeID},
	))
	ts.Require().True(scaled)
}

func (ts *testSuite) TestClient_StartInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		started            = false
	)

	ts.mock().
		On(
			"StartInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(oapi.StartInstanceJSONRequestBody{
				RescueProfile: (*oapi.StartInstanceJSONBodyRescueProfile)(&testInstanceStartRescueProfile),
			}, args.Get(2))
			started = true
		}).
		Return(
			&oapi.StartInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.StartInstance(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		StartInstanceWithRescueProfile(testInstanceStartRescueProfile),
	))
	ts.Require().True(started)
}

func (ts *testSuite) TestClient_StopInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		stopped            = false
	)

	ts.mock().
		On(
			"StopInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			stopped = true
		}).
		Return(
			&oapi.StopInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.StopInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(stopped)
}

func (ts *testSuite) TestClient_UpdateInstance() {
	var (
		testInstanceLabelsUpdated   = map[string]string{"k3": "v3"}
		testInstanceNameUpdated     = testInstanceName + "-updated"
		testInstanceUserDataUpdated = testInstanceUserData + "-updated"
		testOperationID             = ts.randomID()
		testOperationState          = oapi.OperationStateSuccess
		updated                     = false
	)

	ts.mock().
		On(
			"UpdateInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(
				oapi.UpdateInstanceJSONRequestBody{
					Labels:   &oapi.Labels{AdditionalProperties: testInstanceLabelsUpdated},
					Name:     &testInstanceNameUpdated,
					UserData: &testInstanceUserDataUpdated,
				},
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateInstance(context.Background(), testZone, &Instance{
		ID:       &testInstanceID,
		Labels:   &testInstanceLabelsUpdated,
		Name:     &testInstanceNameUpdated,
		UserData: &testInstanceUserDataUpdated,
	}))
	ts.Require().True(updated)
}

func (ts *testSuite) TestClient_GetInstanceReverseDNS() {
	ts.mock().
		On("GetReverseDnsInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
		}).
		Return(&oapi.GetReverseDnsInstanceResponse{
			HTTPResponse: &http.Response{StatusCode: http.StatusOK},
			JSON200: &oapi.ReverseDnsRecord{
				DomainName: (*oapi.DomainName)(&testInstanceReverseDNSDomain),
			},
		}, nil)

	expected := testInstanceReverseDNSDomain

	actual, err := ts.client.GetInstanceReverseDNS(context.Background(), testZone, testInstanceID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *testSuite) TestClient_DeleteInstanceReverseDNS() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	ts.mock().
		On("DeleteReverseDnsInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			deleted = true
		}).
		Return(
			&oapi.DeleteReverseDnsInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.DeleteInstanceReverseDNS(
		context.Background(), 
		testZone,
		testInstanceID,
	))
	ts.Require().True(deleted)
}

func (ts *testSuite) TestClient_UpdateInstanceReverseDNS() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		updated            = false
	)

	ts.mock().
		On("UpdateReverseDnsInstanceWithResponse",
			mock.Anything,                 // ctx
			mock.Anything,                 // id
			mock.Anything,                 // body
			([]oapi.RequestEditorFn)(nil), // reqEditors
		).
		Run(func(args mock.Arguments) {
			ts.Require().Equal(testInstanceID, args.Get(1))
			ts.Require().Equal(
				oapi.UpdateReverseDnsInstanceJSONRequestBody{
					DomainName: &testInstanceReverseDNSDomain,
				}, 
				args.Get(2),
			)
			updated = true
		}).
		Return(
			&oapi.UpdateReverseDnsInstanceResponse{
				HTTPResponse: &http.Response{StatusCode: http.StatusOK},
				JSON200: &oapi.Operation{
					Id:        &testOperationID,
					Reference: oapi.NewReference(nil, &testInstanceID, nil),
					State:     &testOperationState,
				},
			},
			nil,
		)

	ts.mockGetOperation(&oapi.Operation{
		Id:        &testOperationID,
		Reference: oapi.NewReference(nil, &testInstanceID, nil),
		State:     &testOperationState,
	})

	ts.Require().NoError(ts.client.UpdateInstanceReverseDNS(
		context.Background(),
		testZone,
		testInstanceID,
		testInstanceReverseDNSDomain,
	))
	ts.Require().True(updated)
}
