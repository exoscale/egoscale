package v2

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/jarcoal/httpmock"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	papi "github.com/exoscale/egoscale/v2/internal/public-api"
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
	testInstanceManagerType               = papi.ManagerTypeInstancePool
	testInstanceName                      = new(clientTestSuite).randomString(10)
	testInstancePrivateNetworkID          = new(clientTestSuite).randomID()
	testInstancePublicIP                  = "1.2.3.4"
	testInstancePublicIPP                 = net.ParseIP(testInstancePublicIP)
	testInstanceSSHKey                    = new(clientTestSuite).randomString(10)
	testInstanceSecurityGroupID           = new(clientTestSuite).randomID()
	testInstanceSnapshotID                = new(clientTestSuite).randomID()
	testInstanceState                     = papi.InstanceStateRunning
	testInstanceTemplateID                = new(clientTestSuite).randomID()
	testInstanceUserData                  = "I2Nsb3VkLWNvbmZpZwphcHRfdXBncmFkZTogdHJ1ZQ=="
)

func (ts *clientTestSuite) TestClient_AttachInstanceToElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s:attach", testInstanceElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual papi.AttachInstanceToElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AttachInstanceToElasticIpJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState   = papi.OperationStateSuccess
		testPrivateIPAddress = net.ParseIP("10.0.0.1")
		attached             = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/private-network/%s:attach", testInstancePrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual papi.AttachInstanceToPrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AttachInstanceToPrivateNetworkJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
				Ip:       func() *string { ip := testPrivateIPAddress.String(); return &ip }(),
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState = papi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/security-group/%s:attach", testInstanceSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual papi.AttachInstanceToSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.AttachInstanceToSecurityGroupJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState = papi.OperationStateSuccess
	)

	httpmock.RegisterResponder("POST", "/instance",
		func(req *http.Request) (*http.Response, error) {
			var actual papi.CreateInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.CreateInstanceJSONRequestBody{
				AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
				DiskSize:           testInstanceDiskSize,
				InstanceType:       papi.InstanceType{Id: &testInstanceInstanceTypeID},
				Ipv6Enabled:        &testInstanceIPv6Enabled,
				Labels:             &papi.Labels{AdditionalProperties: testInstanceLabels},
				Name:               &testInstanceName,
				SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
				SshKey:             &papi.SshKey{Name: &testInstanceSSHKey},
				Template:           papi.Template{Id: &testInstanceTemplateID},
				UserData:           &testInstanceUserData,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), papi.Instance{
		AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
		CreatedAt:          &testInstanceCreatedAt,
		DiskSize:           &testInstanceDiskSize,
		ElasticIps:         &[]papi.ElasticIp{{Id: &testInstanceElasticIPID}},
		Id:                 &testInstanceID,
		InstanceType:       &papi.InstanceType{Id: &testInstanceInstanceTypeID},
		Ipv6Address:        &testInstanceIPv6Address,
		Labels:             &papi.Labels{AdditionalProperties: testInstanceLabels},
		Manager:            &papi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
		Name:               &testInstanceName,
		PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
		PublicIp:           &testInstancePublicIP,
		SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
		Snapshots:          &[]papi.Snapshot{{Id: &testInstanceSnapshotID}},
		SshKey:             &papi.SshKey{Name: &testInstanceSSHKey},
		State:              &testInstanceState,
		Template:           &papi.Template{Id: &testInstanceTemplateID},
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
		testOperationState = papi.OperationStateSuccess
	)

	ts.mockAPIRequest("POST", fmt.Sprintf("/instance/%s:create-snapshot", testInstanceID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/snapshot/%s", testSnapshotID), papi.Snapshot{
		CreatedAt: &testSnapshotCreatedAt,
		Id:        &testSnapshotID,
		Instance:  &papi.Instance{Id: &testInstanceID},
		Name:      &testSnapshotName,
		State:     &testSnapshotState,
	})

	expected := &Snapshot{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testInstanceID,
		Name:       &testSnapshotName,
		State:      (*string)(&testSnapshotState),
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
		testOperationState = papi.OperationStateSuccess
		deleted            = false
	)

	httpmock.RegisterResponder("DELETE", fmt.Sprintf("/instance/%s", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			deleted = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.DeleteInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(deleted)
}

func (ts *clientTestSuite) TestClient_DetachInstanceFromElasticIP() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		detached           = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/elastic-ip/%s:detach", testInstanceElasticIPID),
		func(req *http.Request) (*http.Response, error) {
			detached = true

			var actual papi.DetachInstanceFromElasticIpJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.DetachInstanceFromElasticIpJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState = papi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/private-network/%s:detach", testInstancePrivateNetworkID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual papi.DetachInstanceFromPrivateNetworkJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.DetachInstanceFromPrivateNetworkJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState = papi.OperationStateSuccess
		attached           = false
	)

	httpmock.RegisterResponder("PUT",
		fmt.Sprintf("/security-group/%s:detach", testInstanceSecurityGroupID),
		func(req *http.Request) (*http.Response, error) {
			attached = true

			var actual papi.DetachInstanceFromSecurityGroupJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.DetachInstanceFromSecurityGroupJSONRequestBody{
				Instance: papi.Instance{Id: &testInstanceID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), papi.Instance{
		AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
		CreatedAt:          &testInstanceCreatedAt,
		DiskSize:           &testInstanceDiskSize,
		ElasticIps:         &[]papi.ElasticIp{{Id: &testInstanceElasticIPID}},
		Id:                 &testInstanceID,
		InstanceType:       &papi.InstanceType{Id: &testInstanceInstanceTypeID},
		Ipv6Address:        &testInstanceIPv6Address,
		Manager:            &papi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
		Name:               &testInstanceName,
		PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
		PublicIp:           &testInstancePublicIP,
		SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
		Snapshots:          &[]papi.Snapshot{{Id: &testInstanceSnapshotID}},
		SshKey:             &papi.SshKey{Name: &testInstanceSSHKey},
		State:              &testInstanceState,
		Template:           &papi.Template{Id: &testInstanceTemplateID},
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
	}

	actual, err := ts.client.GetInstance(context.Background(), testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_FindInstance() {
	ts.mockAPIRequest("GET", "/instance", struct {
		Instances *[]papi.Instance `json:"instances,omitempty"`
	}{
		Instances: &[]papi.Instance{
			{
				CreatedAt:    &testInstanceCreatedAt,
				DiskSize:     &testInstanceDiskSize,
				Id:           &testInstanceID,
				InstanceType: &papi.InstanceType{Id: &testInstanceInstanceTypeID},
				Name:         &testInstanceName,
				State:        &testInstanceState,
				Template:     &papi.Template{Id: &testInstanceTemplateID},
			},
			{
				CreatedAt:    &testInstanceCreatedAt,
				DiskSize:     &testInstanceDiskSize,
				Id:           func() *string { id := ts.randomID(); return &id }(),
				InstanceType: &papi.InstanceType{Id: &testInstanceInstanceTypeID},
				Name:         func() *string { name := "dup"; return &name }(),
				State:        &testInstanceState,
				Template:     &papi.Template{Id: &testInstanceTemplateID},
			},
			{
				CreatedAt:    &testInstanceCreatedAt,
				DiskSize:     &testInstanceDiskSize,
				Id:           func() *string { id := ts.randomID(); return &id }(),
				InstanceType: &papi.InstanceType{Id: &testInstanceInstanceTypeID},
				Name:         func() *string { name := "dup"; return &name }(),
				State:        &testInstanceState,
				Template:     &papi.Template{Id: &testInstanceTemplateID},
			},
		},
	})

	ts.mockAPIRequest("GET", fmt.Sprintf("/instance/%s", testInstanceID), papi.Instance{
		CreatedAt:    &testInstanceCreatedAt,
		DiskSize:     &testInstanceDiskSize,
		Id:           &testInstanceID,
		InstanceType: &papi.InstanceType{Id: &testInstanceInstanceTypeID},
		Name:         &testInstanceName,
		State:        &testInstanceState,
		Template:     &papi.Template{Id: &testInstanceTemplateID},
	})

	expected := &Instance{
		CreatedAt:      &testInstanceCreatedAt,
		DiskSize:       &testInstanceDiskSize,
		ID:             &testInstanceID,
		InstanceTypeID: &testInstanceInstanceTypeID,
		Name:           &testInstanceName,
		State:          (*string)(&testInstanceState),
		TemplateID:     &testInstanceTemplateID,
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
		Instances *[]papi.Instance `json:"instances,omitempty"`
	}{
		Instances: &[]papi.Instance{{
			AntiAffinityGroups: &[]papi.AntiAffinityGroup{{Id: &testInstanceAntiAffinityGroupID}},
			CreatedAt:          &testInstanceCreatedAt,
			DiskSize:           &testInstanceDiskSize,
			ElasticIps:         &[]papi.ElasticIp{{Id: &testInstanceElasticIPID}},
			Id:                 &testInstanceID,
			InstanceType:       &papi.InstanceType{Id: &testInstanceInstanceTypeID},
			Ipv6Address:        &testInstanceIPv6Address,
			Manager:            &papi.Manager{Id: &testInstanceManagerID, Type: &testInstanceManagerType},
			Name:               &testInstanceName,
			PrivateNetworks:    &[]papi.PrivateNetwork{{Id: &testInstancePrivateNetworkID}},
			PublicIp:           &testInstancePublicIP,
			SecurityGroups:     &[]papi.SecurityGroup{{Id: &testInstanceSecurityGroupID}},
			Snapshots:          &[]papi.Snapshot{{Id: &testInstanceSnapshotID}},
			SshKey:             &papi.SshKey{Name: &testInstanceSSHKey},
			State:              &testInstanceState,
			Template:           &papi.Template{Id: &testInstanceTemplateID},
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
	}}

	actual, err := ts.client.ListInstances(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestClient_RebootInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		started            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:reboot", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			started = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.RebootInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestClient_ResetInstance() {
	var (
		testResetDiskSize   int64 = 50
		testResetTemplateID       = ts.randomID()
		testOperationID           = ts.randomID()
		testOperationState        = papi.OperationStateSuccess
		reset                     = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:reset", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			reset = true

			var actual papi.ResetInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.ResetInstanceJSONRequestBody{
				DiskSize: &testResetDiskSize,
				Template: &papi.Template{Id: &testResetTemplateID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
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
		testOperationState       = papi.OperationStateSuccess
		resized                  = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:resize-disk", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			resized = true

			var actual papi.ResizeInstanceDiskJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.ResizeInstanceDiskJSONRequestBody{
				DiskSize: testResizeDiskSize,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
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
		testOperationState = papi.OperationStateSuccess
		reverted           = false
	)

	httpmock.RegisterResponder("POST", fmt.Sprintf("/instance/%s:revert-snapshot", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			reverted = true

			var actual papi.RevertInstanceToSnapshotJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.RevertInstanceToSnapshotJSONRequestBody{Id: testSnapshotID}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
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
		testOperationState      = papi.OperationStateSuccess
		scaled                  = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:scale", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			scaled = true

			var actual papi.ScaleInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.ScaleInstanceJSONRequestBody{
				InstanceType: papi.InstanceType{Id: &testScaleInstanceTypeID},
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testSnapshotID},
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
		testOperationState = papi.OperationStateSuccess
		started            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:start", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			started = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.StartInstance(context.Background(), testZone, &Instance{ID: &testInstanceID}))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestClient_StopInstance() {
	var (
		testOperationID    = ts.randomID()
		testOperationState = papi.OperationStateSuccess
		stopped            = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s:stop", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			stopped = true

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
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
		testOperationState          = papi.OperationStateSuccess
		updated                     = false
	)

	httpmock.RegisterResponder("PUT", fmt.Sprintf("/instance/%s", testInstanceID),
		func(req *http.Request) (*http.Response, error) {
			updated = true

			var actual papi.UpdateInstanceJSONRequestBody
			ts.unmarshalJSONRequestBody(req, &actual)

			expected := papi.UpdateInstanceJSONRequestBody{
				Labels:   &papi.Labels{AdditionalProperties: testInstanceLabelsUpdated},
				Name:     &testInstanceNameUpdated,
				UserData: &testInstanceUserDataUpdated,
			}
			ts.Require().Equal(expected, actual)

			resp, err := httpmock.NewJsonResponse(http.StatusOK, papi.Operation{
				Id:        &testOperationID,
				State:     &testOperationState,
				Reference: &papi.Reference{Id: &testInstanceID},
			})
			if err != nil {
				ts.T().Fatalf("error initializing mock HTTP responder: %s", err)
			}

			return resp, nil
		})

	ts.mockAPIRequest("GET", fmt.Sprintf("/operation/%s", testOperationID), papi.Operation{
		Id:        &testOperationID,
		State:     &testOperationState,
		Reference: &papi.Reference{Id: &testInstanceID},
	})

	ts.Require().NoError(ts.client.UpdateInstance(context.Background(), testZone, &Instance{
		ID:       &testInstanceID,
		Labels:   &testInstanceLabelsUpdated,
		Name:     &testInstanceNameUpdated,
		UserData: &testInstanceUserDataUpdated,
	}))
	ts.Require().True(updated)
}
