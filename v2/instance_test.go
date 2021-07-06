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

func (ts *clientTestSuite) TestInstance_get() {
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

		c:    ts.client,
		zone: testZone,
	}

	actual, err := new(Instance).get(context.Background(), ts.client, testZone, *expected.ID)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_AntiAffinityGroups() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/anti-affinity-group/%s", testAntiAffinityGroupID),
		papi.AntiAffinityGroup{
			Description: &testAntiAffinityGroupDescription,
			Id:          &testAntiAffinityGroupID,
			Name:        &testAntiAffinityGroupName,
		},
	)

	expected := []*AntiAffinityGroup{{
		Description: &testAntiAffinityGroupDescription,
		ID:          &testAntiAffinityGroupID,
		Name:        &testAntiAffinityGroupName,
	}}

	instance := &Instance{
		AntiAffinityGroupIDs: &[]string{testAntiAffinityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instance.AntiAffinityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_AttachElasticIP() {
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

	elasticIP := &ElasticIP{
		ID: &testInstanceElasticIPID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.AttachElasticIP(context.Background(), elasticIP))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestInstance_AttachPrivateNetwork() {
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

	privateNetwork := &PrivateNetwork{
		ID: &testInstancePrivateNetworkID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.AttachPrivateNetwork(context.Background(), privateNetwork, testPrivateIPAddress))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestInstance_AttachSecurityGroup() {
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

	securityGroup := &SecurityGroup{
		ID: &testInstanceSecurityGroupID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.AttachSecurityGroup(context.Background(), securityGroup))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestInstance_CreateSnapshot() {
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

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	expected := &Snapshot{
		CreatedAt:  &testSnapshotCreatedAt,
		ID:         &testSnapshotID,
		InstanceID: &testInstanceID,
		Name:       &testSnapshotName,
		State:      (*string)(&testSnapshotState),

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instance.CreateSnapshot(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_DetachElasticIP() {
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

	elasticIP := &ElasticIP{
		ID: &testInstanceElasticIPID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.DetachElasticIP(context.Background(), elasticIP))
	ts.Require().True(detached)
}

func (ts *clientTestSuite) TestInstance_DetachPrivateNetwork() {
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

	privateNetwork := &PrivateNetwork{
		ID: &testInstancePrivateNetworkID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.DetachPrivateNetwork(context.Background(), privateNetwork))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestInstance_DetachSecurityGroup() {
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

	securityGroup := &SecurityGroup{
		ID: &testInstanceSecurityGroupID,
	}

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.DetachSecurityGroup(context.Background(), securityGroup))
	ts.Require().True(attached)
}

func (ts *clientTestSuite) TestInstance_ElasticIPs() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/elastic-ip/%s", testElasticIPID),
		papi.ElasticIp{
			Id: &testElasticIPID,
			Ip: &testElasticIPAddress,
		},
	)

	expected := []*ElasticIP{{
		ID:        &testElasticIPID,
		IPAddress: &testElasticIPAddressP,

		c:    ts.client,
		zone: testZone,
	}}

	instance := &Instance{
		ElasticIPIDs: &[]string{testElasticIPID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instance.ElasticIPs(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_PrivateNetworks() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/private-network/%s", testPrivateNetworkID),
		papi.PrivateNetwork{
			Id:   &testPrivateNetworkID,
			Name: &testPrivateNetworkName,
		},
	)

	expected := []*PrivateNetwork{{
		ID:   &testPrivateNetworkID,
		Name: &testPrivateNetworkName,
	}}

	instance := &Instance{
		PrivateNetworkIDs: &[]string{testPrivateNetworkID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instance.PrivateNetworks(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_RevertToSnapshot() {
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

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	snapshot := &Snapshot{
		ID:         &testSnapshotID,
		InstanceID: &testInstanceID,
	}

	ts.Require().NoError(instance.RevertToSnapshot(context.Background(), snapshot))
	ts.Require().True(reverted)
}

func (ts *clientTestSuite) TestInstance_SecurityGroups() {
	ts.mockAPIRequest(
		"GET",
		fmt.Sprintf("/security-group/%s", testSecurityGroupID),
		papi.SecurityGroup{
			Id:   &testSecurityGroupID,
			Name: &testSecurityGroupName,
		},
	)

	expected := []*SecurityGroup{{
		ID:   &testSecurityGroupID,
		Name: &testSecurityGroupName,

		c:    ts.client,
		zone: testZone,
	}}

	instance := &Instance{
		SecurityGroupIDs: &[]string{testSecurityGroupID},

		c:    ts.client,
		zone: testZone,
	}

	actual, err := instance.SecurityGroups(context.Background())
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
}

func (ts *clientTestSuite) TestInstance_Reboot() {
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

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.Reboot(context.Background()))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestInstance_Start() {
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

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.Start(context.Background()))
	ts.Require().True(started)
}

func (ts *clientTestSuite) TestInstance_Stop() {
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

	instance := &Instance{
		ID: &testInstanceID,

		c:    ts.client,
		zone: testZone,
	}

	ts.Require().NoError(instance.Stop(context.Background()))
	ts.Require().True(stopped)
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

		c:    ts.client,
		zone: testZone,
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

		c:    ts.client,
		zone: testZone,
	}}

	actual, err := ts.client.ListInstances(context.Background(), testZone)
	ts.Require().NoError(err)
	ts.Require().Equal(expected, actual)
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

		c:    ts.client,
		zone: testZone,
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

		c:    ts.client,
		zone: testZone,
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

	ts.Require().NoError(ts.client.DeleteInstance(context.Background(), testZone, testInstanceID))
	ts.Require().True(deleted)
}
