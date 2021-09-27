package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	"github.com/exoscale/egoscale/v2/oapi"
)

var (
	testInstanceAntiAffinityGroupID       = new(clientTestSuite).randomID()
	testInstanceCreatedAt, _              = time.Parse(iso8601Format, "2020-05-26T12:09:42Z")
	testInstanceDiskSize            int64 = 10
	testInstanceElasticIPID               = new(clientTestSuite).randomID()
	testInstanceID                        = new(clientTestSuite).randomID()
	testInstanceIPv6Address               = "2001:db8:abcd::1"
	testInstanceIPv6AddressP              = net.ParseIP(testInstanceIPv6Address)
	testInstanceIPv6Enabled               = true
	testInstanceInstanceTypeID            = new(clientTestSuite).randomID()
	testInstanceLabels                    = map[string]string{"k1": "v1", "k2": "v2"}
	testInstanceManagerID                 = new(clientTestSuite).randomID()
	testInstanceManagerType               = oapi.ManagerTypeInstancePool
	testInstanceName                      = new(clientTestSuite).randomString(10)
	testInstancePrivateNetworkID          = new(clientTestSuite).randomID()
	testInstancePublicIP                  = "1.2.3.4"
	testInstancePublicIPP                 = net.ParseIP(testInstancePublicIP)
	testInstanceSSHKey                    = new(clientTestSuite).randomString(10)
	testInstanceSecurityGroupID           = new(clientTestSuite).randomID()
	testInstanceSnapshotID                = new(clientTestSuite).randomID()
	testInstanceState                     = oapi.InstanceStateRunning
	testInstanceTemplateID                = new(clientTestSuite).randomID()
	testInstanceUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *clientTestSuite) TestClient_AttachInstanceToElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s:attach", testInstanceElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual oapi.AttachInstanceToElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.AttachInstanceToElasticIpJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.AttachInstanceToElasticIP(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&ElasticIP{ID: &testInstanceElasticIPID},
	))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestClient_AttachInstanceToPrivateNetwork() {
	var (
		testOperationID      = ts.randomID()
		testOperationState   = oapi.OperationStateSuccess
		testPrivateIPAddress = net.ParseIP("10.0.0.1")
		attached             = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/private-network/%s:attach", testInstancePrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual oapi.AttachInstanceToPrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.AttachInstanceToPrivateNetworkJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
				Ip:       func() *string { ip := testPrivateIPAddress.String(); return &ip }(),
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.AttachInstanceToPrivateNetwork(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&PrivateNetwork{ID: &testInstancePrivateNetworkID},
		testPrivateIPAddress,
	))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestClient_AttachInstanceToSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/security-group/%s:attach", testInstanceSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual oapi.AttachInstanceToSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.AttachInstanceToSecurityGroupJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.AttachInstanceToSecurityGroup(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&SecurityGroup{ID: &testInstanceSecurityGroupID},
	))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestClient_CreateInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/instance",
		func(req *http.Request) (*http.Response, error) {
			var actual oapi.CreateInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.CreateInstanceJSONRequestBody{
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
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), oapi.Instance{
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
	})

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

func (ts *clientTestSuite) TestClient_CreateInstanceSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
	)

	ts.mockAPIRequest("POST", fmt.Sprintf("/instance/%s:create-snapshot", testInstanceID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/snapshot/%s", testSnapshotID), oapi.Snapshot{
		CreatedAt: &testSnapshotCreatedAt,
		Id:        &testSnapshotID,
		Instance:  &oapi.Instance{Id: &testInstanceID},
		Name:      &testSnapshotName,
		State:     &testSnapshotState,
	})

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

func (ts *clientTestSuite) TestClient_DeleteInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/instance/%s", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.DeleteInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_DetachInstanceFromElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		detached           = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s:detach", testInstanceElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			detached = true

			var actual oapi.DetachInstanceFromElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.DetachInstanceFromElasticIpJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.DetachInstanceFromElasticIP(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&ElasticIP{ID: &testInstanceElasticIPID},
	))
	ts.Require().True(detached)
}

func (ts *clientTestSuite) TestClient_DetachInstanceFromPrivateNetwork() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/private-network/%s:detach", testInstancePrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual oapi.DetachInstanceFromPrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.DetachInstanceFromPrivateNetworkJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.DetachInstanceFromPrivateNetwork(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&PrivateNetwork{ID: &testInstancePrivateNetworkID},
	))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestClient_DetachInstanceFromSecurityGroup() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/security-group/%s:detach", testInstanceSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual oapi.DetachInstanceFromSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.DetachInstanceFromSecurityGroupJSONRequestBody{
				Instance: oapi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.DetachInstanceFromSecurityGroup(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&SecurityGroup{ID: &testInstanceSecurityGroupID},
	))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestClient_GetInstance() {
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), oapi.Instance{
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
	})

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

func (ts *clientTestSuite) TestClient_FindInstance() {
	ts.mockAPIRequest("GET", "/instance", struct {
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
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), oapi.Instance{
		CreatedAt:    &testInstanceCreatedAt,
		DiskSize:     &testInstanceDiskSize,
		Id:           &testInstanceID,
		InstanceType: &oapi.InstanceType{Id: &testInstanceInstanceTypeID},
		Name:         &testInstanceName,
		State:        &testInstanceState,
		Template:     &oapi.Template{Id: &testInstanceTemplateID},
	})

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

func (ts *clientTestSuite) TestClient_ListInstances() {
	ts.mockAPIRequest("GET", "/instance", struct {
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
	})

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

	actual, err := ts.client.ListInstances(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_RebootInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		started            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:reboot", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			started = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.RebootInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestClient_ResetInstance() {
	var (
		testResetDiskSize   int64 = 50
		testResetTemplateID       = ts.randomID()
		testOperationID           = ts.randomID()
		testOperationState        = oapi.OperationStateSuccess
		reset                     = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:reset", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			reset = true

			var actual oapi.ResetInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.ResetInstanceJSONRequestBody{
				DiskSize: &testResetDiskSize,
				Template: &oapi.Template{Id: &testResetTemplateID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
	})

	ts.Require().NoError(ts.client.ResetInstance(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&Template{ID: &testResetTemplateID},
		testResetDiskSize,
	))
	ts.Require().True(reset)
}

func (ts *clientTestSuite) TestClient_ResizeInstanceDisk() {
	var (
		testResizeDiskSize int64 = 50
		testOperationID          = ts.randomID()
		testOperationState       = oapi.OperationStateSuccess
		resized                  = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:resize-disk", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			resized = true

			var actual oapi.ResizeInstanceDiskJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.ResizeInstanceDiskJSONRequestBody{
				DiskSize: testResizeDiskSize,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
	})

	ts.Require().NoError(ts.client.ResizeInstanceDisk(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		testResizeDiskSize),
	)
	ts.Require().True(resized)
}

func (ts *clientTestSuite) TestClient_RevertInstanceToSnapshot() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		reverted           = false
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/instance/%s:revert-snapshot", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			reverted = true

			var actual oapi.RevertInstanceToSnapshotJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.RevertInstanceToSnapshotJSONRequestBody{Id: testSnapshotID}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
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

func (ts *clientTestSuite) TestClient_ScaleInstance() {
	var (
		testScaleInstanceTypeID = ts.randomID()
		testOperationID         = ts.randomID()
		testOperationState      = oapi.OperationStateSuccess
		scaled                  = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:scale", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

			var actual oapi.ScaleInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.ScaleInstanceJSONRequestBody{
				InstanceType: oapi.InstanceType{Id: &testScaleInstanceTypeID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testSnapshotID},
	})

	ts.Require().NoError(ts.client.ScaleInstance(
		context.Background(),
		testZone,
		&Instance{ID: &testInstanceID},
		&InstanceType{ID: &testScaleInstanceTypeID},
	))
	ts.Require().True(scaled)
}

func (ts *clientTestSuite) TestClient_StartInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		started            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:start", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			started = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.StartInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestClient_StopInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = oapi.OperationStateSuccess
		stopped            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:stop", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			stopped = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.StopInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(stopped)
}

func (ts *clientTestSuite) TestClient_UpdateInstance() {
	var (
		testInstanceLabelsUpdated   = map[string]string{"k3": "v3"}
		testInstanceNameUpdated     = testInstanceName + "-updated"
		testInstanceUserDataUpdated = testInstanceUserData + "-updated"
		testOperationID             = ts.randomID()
		testOperationState          = oapi.OperationStateSuccess
		updated                     = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual oapi.UpdateInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := oapi.UpdateInstanceJSONRequestBody{
				Labels:   &oapi.Labels{AdditionalProperties: testInstanceLabelsUpdated},
				Name:     &testInstanceNameUpdated,
				UserData: &testInstanceUserDataUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, oapi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &oapi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), oapi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &oapi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.UpdateInstance(context.Background(), testZone, &Instance{
		ID:       &testInstanceID,
		Labels:   &testInstanceLabelsUpdated,
		Name:     &testInstanceNameUpdated,
		UserData: &testInstanceUserDataUpdated,
	}))
	ts.Require().True(updated)
}
