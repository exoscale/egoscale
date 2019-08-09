// +build testacc

package compute

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	egoerr "github.com/exoscale/egoscale/error"
	"github.com/exoscale/egoscale/internal/egoscale"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type instanceFixture struct {
	c   *Client
	req *egoapi.DeployVirtualMachine
	res *egoapi.VirtualMachine
}

func newInstanceFixture(c *Client, opts ...instanceFixtureOpt) *instanceFixture {
	var fixture = &instanceFixture{
		c:   c,
		req: &egoapi.DeployVirtualMachine{},
	}

	// Fixture default options
	for _, opt := range []instanceFixtureOpt{
		instanceFixtureOptZone(testZoneID),
		instanceFixtureOptName(testPrefix + "-" + testRandomString()),
		instanceFixtureOptType(testInstanceServiceOfferingID),
		instanceFixtureOptTemplate(testInstanceTemplateID),
	} {
		opt(fixture)
	}

	for _, opt := range opts {
		opt(fixture)
	}

	return fixture
}

func (f *instanceFixture) setup() (*instanceFixture, error) { // nolint:unused,deadcode
	res, err := f.c.c.RequestWithContext(f.c.ctx, f.req)
	if err != nil {
		return nil, f.c.csError(err)
	}
	f.res = res.(*egoscale.VirtualMachine)

	return f, nil
}

func (f *instanceFixture) teardown() error { // nolint:unused,deadcode
	_, err := f.c.c.RequestWithContext(f.c.ctx, &egoapi.DestroyVirtualMachine{ID: f.res.ID})
	return f.c.csError(err)
}

type instanceFixtureOpt func(*instanceFixture)

func instanceFixtureOptZone(id string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) { f.req.ZoneID = egoapi.MustParseUUID(id) }
}

func instanceFixtureOptName(name string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) { f.req.Name = name; f.req.DisplayName = name }
}

func instanceFixtureOptType(typeID string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) { f.req.ServiceOfferingID = egoapi.MustParseUUID(typeID) }
}

func instanceFixtureOptTemplate(templateID string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) { f.req.TemplateID = egoapi.MustParseUUID(templateID) }
}

func instanceFixtureOptAntiAffinityGroups(ids []string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) {
		reqIDs := make([]egoapi.UUID, len(ids))
		for i := range ids {
			reqIDs[i] = *(egoapi.MustParseUUID(ids[i]))
		}
		f.req.AffinityGroupIDs = reqIDs
	}
}

func instanceFixtureOptSecurityGroups(ids []string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) {
		reqIDs := make([]egoapi.UUID, len(ids))
		for i := range ids {
			reqIDs[i] = *(egoapi.MustParseUUID(ids[i]))
		}
		f.req.SecurityGroupIDs = reqIDs
	}
}

func instanceFixtureOptPrivateNetworks(ids []string) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) {
		reqIDs := make([]egoapi.UUID, len(ids))
		for i := range ids {
			reqIDs[i] = *(egoapi.MustParseUUID(ids[i]))
		}
		f.req.NetworkIDs = reqIDs
	}
}

func instanceFixtureOptStart(flag bool) instanceFixtureOpt { // nolint:unused,deadcode
	return func(f *instanceFixture) { f.req.StartVM = &flag }
}

func (t *accTestSuite) withInstanceFixture(f func(*instanceFixture), opts ...instanceFixtureOpt) {
	instanceFixture, err := newInstanceFixture(t.client, opts...).setup()
	if err != nil {
		t.FailNow("Compute instance fixture setup failed", err)
	}

	f(instanceFixture)
}

type instanceTestSuite struct {
	suite.Suite
	*accTestSuite
}

func (t *instanceTestSuite) SetupTest() {
	client, err := testClientFromEnv()
	if err != nil {
		t.FailNow("unable to initialize API client", err)
	}

	t.accTestSuite = &accTestSuite{
		Suite:  t.Suite,
		client: client,
	}
}

func (t *instanceTestSuite) TestCreateInstance() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withInstanceTemplateFixture(func(instanceTemplateFixture *instanceTemplateFixture) {

			t.withInstanceTypeFixture(func(instanceTypeFixture *instanceTypeFixture) {

				t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
					defer antiAffinityGroupFixture.teardown() // nolint:errcheck

					t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
						defer securityGroupFixture.teardown() // nolint:errcheck

						t.withPrivateNetworkFixture(func(privateNetworkFixture *privateNetworkFixture) {
							defer privateNetworkFixture.teardown() // nolint:errcheck

							t.withSSHKeyFixture(func(sshKeyFixture *sshKeyFixture) {
								defer sshKeyFixture.teardown() // nolint:errcheck

								zone := t.client.zoneFromAPI(zoneFixture.res)
								instanceName := testPrefix + "-" + testRandomString()
								instanceType := t.client.instanceTypeFromAPI(instanceTypeFixture.res)
								instanceTemplate, err := t.client.instanceTemplateFromAPI(instanceTemplateFixture.res)
								if err != nil {
									t.FailNow("instance template fixture setup failed", err)
								}
								antiAffinityGroup := t.client.antiAffinityGroupFromAPI(antiAffinityGroupFixture.res)
								securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)
								privateNetwork, err := t.client.privateNetworkFromAPI(privateNetworkFixture.res)
								if err != nil {
									t.FailNow("Private Network fixture setup failed", err)
								}
								sshKey := t.client.sshKeyFromAPI(sshKeyFixture.res)

								instance, err := t.client.CreateInstance(
									zone,
									&InstanceCreateOpts{
										Name:               instanceName,
										Type:               instanceType,
										Template:           instanceTemplate,
										VolumeSize:         20,
										AntiAffinityGroups: []*AntiAffinityGroup{antiAffinityGroup},
										SecurityGroups:     []*SecurityGroup{securityGroup},
										PrivateNetworks:    []*PrivateNetwork{privateNetwork},
										SSHKey:             sshKey,
										EnableIPv6:         true,
										UserData:           fmt.Sprintln("#cloud-config\npackage_upgrade: true"),
									})
								if err != nil {
									t.FailNow("instance creation failed", err)
								}

								actualInstance := egoapi.VirtualMachine{}
								if err := json.Unmarshal(instance.Raw(), &actualInstance); err != nil {
									t.FailNow("unable to unmarshal raw resource", err)
								}
								res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.ListVolumes{
									VirtualMachineID: actualInstance.ID,
								})
								if err != nil {
									t.FailNow("instance volumes listing failed", err)
								}
								actualInstanceVolumeSize := int64(res.(*egoapi.ListVolumesResponse).Volume[0].Size)

								assert.Equal(t.T(), actualInstance.ID.String(), instance.ID)
								assert.Equal(t.T(), zone.ID, instance.Zone.ID)
								assert.Equal(t.T(), instanceName, actualInstance.Name)
								assert.Equal(t.T(), instanceName, instance.Name)
								assert.Equal(t.T(), instanceType.ID, actualInstance.ServiceOfferingID.String())
								assert.Equal(t.T(), instanceType.ID, instance.Type.ID)
								assert.Equal(t.T(), instanceTemplate.ID, actualInstance.TemplateID.String())
								assert.Equal(t.T(), instanceTemplate.ID, instance.Template.ID)
								assert.Equal(t.T(), int64(21474836480), actualInstanceVolumeSize) // 20 GB
								assert.Equal(t.T(), int64(21474836480), instance.VolumeSize)
								assert.Equal(t.T(), sshKey.Name, instance.SSHKey.Name)
								assert.True(t.T(), instance.IPv4Address.Equal(actualInstance.DefaultNic().IPAddress))
								assert.True(t.T(), instance.IPv6Address.Equal(actualInstance.DefaultNic().IP6Address))
								assert.Len(t.T(), actualInstance.AffinityGroup, 1)
								assert.True(t.T(),
									antiAffinityGroupFixture.res.ID.Equal(*actualInstance.AffinityGroup[0].ID))
								assert.Len(t.T(), actualInstance.SecurityGroup, 1)
								assert.True(t.T(),
									securityGroupFixture.res.ID.Equal(*actualInstance.SecurityGroup[0].ID))
								assert.NotNil(t.T(), actualInstance.NicByNetworkID(*privateNetworkFixture.res.ID))

								if _, err = t.client.c.RequestWithContext(t.client.ctx, &egoapi.DestroyVirtualMachine{
									ID: egoapi.MustParseUUID(instance.ID),
								}); err != nil {
									t.FailNow("instance deletion failed", err)
								}
							})
						})
					})
				})
			})
		})
	})
}

func (t *instanceTestSuite) TestListInstances() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)

			instances, err := t.client.ListInstances(zone)
			if err != nil {
				t.FailNow("Private Networks listing failed", err)
			}

			// We cannot guarantee that there will be only our resources in the
			// testing environment, so we ensure we get at least our fixture instance
			assert.GreaterOrEqual(t.T(), len(instances), 1)
		})
	})
}

func (t *instanceTestSuite) TestGetInstanceByID() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)

			instance, err := t.client.GetInstanceByID(zone, instanceFixture.res.ID.String())
			if err != nil {
				t.FailNow("instance retrieval failed", err)
			}
			assert.Equal(t.T(), instanceFixture.res.ID.String(), instance.ID)

			instance, err = t.client.GetInstanceByID(zone, "00000000-0000-0000-0000-000000000000")
			assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
			assert.Empty(t.T(), instance)
		})
	})
}

func (t *instanceTestSuite) TestGetInstanceByAddress() {
	t.withZoneFixture(func(zoneFixture *zoneFixture) {

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			zone := t.client.zoneFromAPI(zoneFixture.res)

			instance, err := t.client.GetInstanceByAddress(zone, instanceFixture.res.IP().String())
			if err != nil {
				t.FailNow("instance retrieval failed", err)
			}
			assert.Equal(t.T(), instanceFixture.res.ID.String(), instance.ID)

			instance, err = t.client.GetInstanceByAddress(zone, "1.2.3.4")
			assert.EqualError(t.T(), err, egoerr.ErrResourceNotFound.Error())
			assert.Empty(t.T(), instance)
		})
	})
}

func (t *instanceTestSuite) TestInstanceAntiAffinityGroups() {
	t.withAntiAffinityGroupFixture(func(antiAffinityGroupFixture *antiAffinityGroupFixture) {
		defer antiAffinityGroupFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			instanceAntiAffinityGroups, err := instance.AntiAffinityGroups()
			if err != nil {
				t.FailNow("Compute instance Anti-Affinity Groups retrieval failed", err)
			}
			assert.Len(t.T(), instanceAntiAffinityGroups, 1)
			assert.Equal(t.T(), antiAffinityGroupFixture.res.ID.String(), instanceAntiAffinityGroups[0].ID)
		}, instanceFixtureOptAntiAffinityGroups([]string{antiAffinityGroupFixture.res.ID.String()}))
	})
}

func (t *instanceTestSuite) TestInstanceElasticIPs() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.AddIPToNic{
				NicID:     instanceFixture.res.DefaultNic().ID,
				IPAddress: elasticIPFixture.res.IPAddress,
			}); err != nil {
				t.FailNow("unable to attach Elastic IP to Compute instance fixture", err)
			}

			instanceElasticIPs, err := instance.ElasticIPs()
			if err != nil {
				t.FailNow("Compute instance Elastic IPs retrieval failed", err)
			}
			assert.Len(t.T(), instanceElasticIPs, 1)
			assert.Equal(t.T(), elasticIPFixture.res.ID.String(), instanceElasticIPs[0].ID)
		})
	})
}

func (t *instanceTestSuite) TestInstanceAttachElasticIP() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}
			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if err := instance.AttachElasticIP(elasticIP); err != nil {
				t.FailNow("Compute instance Elastic IP attachment failed", err)
			}

			res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.ListNics{
				NicID: instanceFixture.res.DefaultNic().ID,
			})
			if err != nil {
				t.FailNow("unable to retrieve Compute instance fixture NIC", err)
			}
			instanceNICs := res.(*egoapi.ListNicsResponse).Nic
			assert.Len(t.T(), instanceNICs, 1)
			assert.Len(t.T(), instanceNICs[0].SecondaryIP, 1)
			assert.True(t.T(), instanceNICs[0].SecondaryIP[0].IPAddress.Equal(elasticIPFixture.res.IPAddress))
		})
	})
}

func (t *instanceTestSuite) TestInstanceDetachElasticIP() {
	t.withElasticIPFixture(func(elasticIPFixture *elasticIPFixture) {
		defer elasticIPFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			elasticIP, err := t.client.elasticIPFromAPI(elasticIPFixture.res)
			if err != nil {
				t.FailNow("Elastic IP fixture setup failed", err)
			}
			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.AddIPToNic{
				NicID:     instanceFixture.res.DefaultNic().ID,
				IPAddress: elasticIPFixture.res.IPAddress,
			}); err != nil {
				t.FailNow("unable to attach Elastic IP to Compute instance fixture", err)
			}

			if err := instance.DetachElasticIP(elasticIP); err != nil {
				t.FailNow("Compute instance Elastic IP attachment failed", err)
			}

			res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.ListNics{
				NicID: instanceFixture.res.DefaultNic().ID,
			})
			if err != nil {
				t.FailNow("unable to retrieve Compute instance fixture NIC", err)
			}
			instanceNICs := res.(*egoapi.ListNicsResponse).Nic
			assert.Len(t.T(), instanceNICs, 1)
			assert.Len(t.T(), instanceNICs[0].SecondaryIP, 0)
		})
	})
}

// func (t *instanceTestSuite) TestInstancePrivateNetworks() {
// }

// func (t *instanceTestSuite) TestInstanceAttachPrivateNetwork() {
// }

// func (t *instanceTestSuite) TestInstanceDetachPrivateNetwork() {
// }

func (t *instanceTestSuite) TestInstanceSecurityGroups() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			instanceSecurityGroups, err := instance.SecurityGroups()
			if err != nil {
				t.FailNow("Compute instance Security Groups retrieval failed", err)
			}
			assert.Len(t.T(), instanceSecurityGroups, 1)
			assert.Equal(t.T(), securityGroupFixture.res.ID.String(), instanceSecurityGroups[0].ID)
		}, instanceFixtureOptSecurityGroups([]string{securityGroupFixture.res.ID.String()}))
	})
}

func (t *instanceTestSuite) TestInstanceReverseDNS() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.UpdateReverseDNSForVirtualMachine{
			ID:         instanceFixture.res.ID,
			DomainName: testReverseDNS,
		}); err != nil {
			t.FailNow("Compute instance fixture reverse DNS setting failed", err)
		}

		reverseDNS, err := instance.ReverseDNS()
		if err != nil {
			t.FailNow("instance reverse DNS retrieval failed", err)
		}

		assert.Equal(t.T(), testReverseDNS, reverseDNS)
	})
}

func (t *instanceTestSuite) TestInstanceSetReverseDNS() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if err = instance.SetReverseDNS(testReverseDNS); err != nil {
			t.FailNow("instance reverse DNS setting failed", err)
		}

		res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.QueryReverseDNSForVirtualMachine{
			ID: instanceFixture.res.ID,
		})
		if err != nil {
			t.FailNow("Compute instance fixture retrieval failed", err)
		}
		reverseDNS := res.(*egoapi.VirtualMachine).DefaultNic().ReverseDNS

		assert.Len(t.T(), reverseDNS, 1)
		assert.Equal(t.T(), testReverseDNS, reverseDNS[0].DomainName)
	})
}

func (t *instanceTestSuite) TestInstanceUnsetReverseDNS() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if _, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.UpdateReverseDNSForVirtualMachine{
			ID:         instanceFixture.res.ID,
			DomainName: testReverseDNS,
		}); err != nil {
			t.FailNow("Compute instance fixture reverse DNS setting failed", err)
		}

		if err := instance.UnsetReverseDNS(); err != nil {
			t.FailNow("instance reverse DNS deletion failed", err)
		}

		res, err := t.client.c.RequestWithContext(t.client.ctx, &egoapi.QueryReverseDNSForVirtualMachine{
			ID: instanceFixture.res.ID,
		})
		if err != nil {
			t.FailNow("Compute instance fixture retrieval failed", err)
		}
		reverseDNS := res.(*egoapi.VirtualMachine).DefaultNic().ReverseDNS

		assert.Len(t.T(), reverseDNS, 0)
	})
}

// func (t *instanceTestSuite) TestInstanceResizeVolume() {
// }

// func (t *instanceTestSuite) TestInstanceSnapshotVolume() {
// }

// func (t *instanceTestSuite) TestInstanceVolumeSnapshots() {
// }

func (t *instanceTestSuite) TestInstanceState() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		instanceState, err := instance.State()
		if err != nil {
			t.FailNow("instance state retrieval failed", err)
		}
		assert.Equal(t.T(), "running", instanceState)
	})
}

func (t *instanceTestSuite) TestInstanceUpdate() {
	t.withSecurityGroupFixture(func(securityGroupFixture *securityGroupFixture) {
		defer securityGroupFixture.teardown() // nolint:errcheck

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			instanceNameEdited := instance.Name + "-edited"
			securityGroup := t.client.securityGroupFromAPI(securityGroupFixture.res)

			if err := instance.Update(
				&InstanceUpdateOpts{
					Name:           instanceNameEdited,
					SecurityGroups: []*SecurityGroup{securityGroup},
					UserData:       fmt.Sprintln("#cloud-config\npackage_upgrade: true"),
				}); err != nil {
				t.FailNow("instance update failed", err)
			}

			res, err := t.client.c.GetWithContext(t.client.ctx, &egoapi.VirtualMachine{ID: instanceFixture.res.ID})
			if err != nil {
				t.FailNow("Compute instance fixture retrieval failed", err)
			}
			actualInstance := res.(*egoapi.VirtualMachine)

			assert.Equal(t.T(), instanceNameEdited, actualInstance.Name)
			assert.Equal(t.T(), instanceNameEdited, instance.Name)
			assert.True(t.T(), securityGroupFixture.res.ID.Equal(*actualInstance.SecurityGroup[0].ID))
		}, instanceFixtureOptStart(false))
	})
}

func (t *instanceTestSuite) TestInstanceScale() {
	t.withInstanceTypeFixture(func(instanceTypeFixture *instanceTypeFixture) {

		t.withInstanceFixture(func(instanceFixture *instanceFixture) {
			defer instanceFixture.teardown() // nolint:errcheck

			instance, err := t.client.instanceFromAPI(instanceFixture.res)
			if err != nil {
				t.FailNow("Compute instance fixture setup failed", err)
			}

			instanceTypeSmall := t.client.instanceTypeFromAPI(instanceTypeFixture.res)

			if err := instance.Scale(instanceTypeSmall); err != nil {
				t.FailNow("instance scaling failed", err)
			}

			res, err := t.client.c.GetWithContext(t.client.ctx, &egoapi.VirtualMachine{ID: instanceFixture.res.ID})
			if err != nil {
				t.FailNow("Compute instance fixture retrieval failed", err)
			}
			assert.Equal(t.T(), instanceTypeSmall.ID, res.(*egoapi.VirtualMachine).ServiceOfferingID.String())
		}, instanceFixtureOptStart(false))
	}, instanceTypeFixtureOptName("small"))
}

func (t *instanceTestSuite) TestInstanceStart() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if err := instance.Start(nil); err != nil {
			t.FailNow("instance start failed", err)
		}

		res, err := t.client.c.GetWithContext(t.client.ctx, &egoapi.VirtualMachine{ID: instanceFixture.res.ID})
		if err != nil {
			t.FailNow("Compute instance fixture retrieval failed", err)
		}
		actualInstanceState := strings.ToLower(res.(*egoapi.VirtualMachine).State)
		assert.True(t.T(), actualInstanceState == "starting" || actualInstanceState == "running")
	}, instanceFixtureOptStart(false))
}

func (t *instanceTestSuite) TestInstanceStop() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if err := instance.Stop(); err != nil {
			t.FailNow("instance stop failed", err)
		}

		res, err := t.client.c.GetWithContext(t.client.ctx, &egoapi.VirtualMachine{ID: instanceFixture.res.ID})
		if err != nil {
			t.FailNow("Compute instance fixture retrieval failed", err)
		}
		actualInstanceState := strings.ToLower(res.(*egoapi.VirtualMachine).State)
		assert.True(t.T(), actualInstanceState == "stopping" || actualInstanceState == "stopped")
	})
}

func (t *instanceTestSuite) TestInstanceReboot() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if err := instance.Reboot(); err != nil {
			t.FailNow("instance reboot failed", err)
		}

		res, err := t.client.c.GetWithContext(t.client.ctx, &egoapi.VirtualMachine{ID: instanceFixture.res.ID})
		if err != nil {
			t.FailNow("Compute instance fixture retrieval failed", err)
		}
		actualInstanceState := strings.ToLower(res.(*egoapi.VirtualMachine).State)
		assert.True(t.T(), actualInstanceState == "starting" || actualInstanceState == "running")
	})
}

func (t *instanceTestSuite) TestInstanceDelete() {
	t.withInstanceFixture(func(instanceFixture *instanceFixture) {
		defer instanceFixture.teardown() // nolint:errcheck

		instance, err := t.client.instanceFromAPI(instanceFixture.res)
		if err != nil {
			t.FailNow("Compute instance fixture setup failed", err)
		}

		if err = instance.Delete(); err != nil {
			t.FailNow("Compute instance deletion failed", err)
		}
		assert.Empty(t.T(), instance.ID)
		assert.Empty(t.T(), instance.Zone)

		r, _ := t.client.c.ListWithContext(t.client.ctx, &egoapi.VirtualMachine{
			ZoneID: instanceFixture.res.ZoneID,
			ID:     instanceFixture.res.ID,
		})
		assert.Len(t.T(), r, 0)
	})
}

func TestAccComputeInstanceTestSuite(t *testing.T) {
	suite.Run(t, new(instanceTestSuite))
}
