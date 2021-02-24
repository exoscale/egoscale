package v2

import (
	"context"

	apiv2 "github.com/exoscale/egoscale/v2/api"
	papi "github.com/exoscale/egoscale/v2/internal/public-api"
)

// InstancePool represents an Instance Pool.
type InstancePool struct {
	AntiAffinityGroupIDs []string
	Description          string
	DiskSize             int64
	ElasticIPIDs         []string
	ID                   string
	IPv6Enabled          bool
	InstanceIDs          []string
	InstanceTypeID       string
	ManagerID            string
	Name                 string
	PrivateNetworkIDs    []string
	SSHKey               string
	SecurityGroupIDs     []string
	Size                 int64
	State                string
	TemplateID           string
	UserData             string

	c    *Client
	zone string
}

func instancePoolFromAPI(i *papi.InstancePool) *InstancePool {
	return &InstancePool{
		AntiAffinityGroupIDs: func() []string {
			ids := make([]string, 0)

			if i.AntiAffinityGroups != nil {
				for _, item := range *i.AntiAffinityGroups {
					item := item
					ids = append(ids, *item.Id)
				}
			}

			return ids
		}(),
		Description: papi.OptionalString(i.Description),
		DiskSize:    papi.OptionalInt64(i.DiskSize),
		ElasticIPIDs: func() []string {
			ids := make([]string, 0)

			if i.ElasticIps != nil {
				for _, item := range *i.ElasticIps {
					item := item
					ids = append(ids, *item.Id)
				}
			}

			return ids
		}(),
		ID:          papi.OptionalString(i.Id),
		IPv6Enabled: papi.OptionalBool(i.Ipv6Enabled),
		InstanceIDs: func() []string {
			ids := make([]string, 0)

			if i.Instances != nil {
				for _, item := range *i.Instances {
					item := item
					ids = append(ids, *item.Id)
				}
			}

			return ids
		}(),
		InstanceTypeID: papi.OptionalString(i.InstanceType.Id),
		ManagerID:      papi.OptionalString(i.Manager.Id),
		Name:           papi.OptionalString(i.Name),
		PrivateNetworkIDs: func() []string {
			ids := make([]string, 0)

			if i.PrivateNetworks != nil {
				for _, item := range *i.PrivateNetworks {
					item := item
					ids = append(ids, *item.Id)
				}
			}

			return ids
		}(),
		SSHKey: func() string {
			key := ""
			if i.SshKey != nil {
				key = papi.OptionalString(i.SshKey.Name)
			}
			return key
		}(),
		SecurityGroupIDs: func() []string {
			ids := make([]string, 0)

			if i.SecurityGroups != nil {
				for _, item := range *i.SecurityGroups {
					item := item
					ids = append(ids, *item.Id)
				}
			}

			return ids
		}(),
		Size:       papi.OptionalInt64(i.Size),
		State:      papi.OptionalString(i.State),
		TemplateID: papi.OptionalString(i.Template.Id),
		UserData:   papi.OptionalString(i.UserData),
	}
}

// Scale scales the Instance Pool to the specified number of instances.
func (i *InstancePool) Scale(ctx context.Context, instances int64) error {
	resp, err := i.c.ScaleInstancePoolWithResponse(
		apiv2.WithZone(ctx, i.zone),
		i.ID,
		papi.ScaleInstancePoolJSONRequestBody{Size: instances},
	)
	if err != nil {
		return err
	}

	_, err = papi.NewPoller().
		WithTimeout(i.c.timeout).
		Poll(ctx, i.c.OperationPoller(i.zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}

// EvictNodepoolMembers evicts the specified members (identified by their Compute instance ID) from the
// Instance Pool.
func (i *InstancePool) EvictMembers(ctx context.Context, members []string) error {
	resp, err := i.c.EvictInstancePoolMembersWithResponse(
		apiv2.WithZone(ctx, i.zone),
		i.ID,
		papi.EvictInstancePoolMembersJSONRequestBody{Instances: &members},
	)
	if err != nil {
		return err
	}

	_, err = papi.NewPoller().
		WithTimeout(i.c.timeout).
		Poll(ctx, i.c.OperationPoller(i.zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}

// CreateInstancePool creates an Instance Pool in the specified zone.
func (c *Client) CreateInstancePool(ctx context.Context, zone string, instancePool *InstancePool) (*InstancePool, error) {
	resp, err := c.CreateInstancePoolWithResponse(
		apiv2.WithZone(ctx, zone),
		papi.CreateInstancePoolJSONRequestBody{
			Name:        instancePool.Name,
			Description: &instancePool.Description,
		})
	if err != nil {
		return nil, err
	}

	res, err := papi.NewPoller().
		WithTimeout(c.timeout).
		Poll(ctx, c.OperationPoller(zone, *resp.JSON200.Id))
	if err != nil {
		return nil, err
	}

	return c.GetInstancePool(ctx, zone, *res.(*papi.Reference).Id)
}

// ListInstancePools returns the list of existing Instance Pools in the specified zone.
func (c *Client) ListInstancePools(ctx context.Context, zone string) ([]*InstancePool, error) {
	list := make([]*InstancePool, 0)

	resp, err := c.ListInstancePoolsWithResponse(apiv2.WithZone(ctx, zone))
	if err != nil {
		return nil, err
	}

	if resp.JSON200.InstancePools != nil {
		for i := range *resp.JSON200.InstancePools {
			cluster := instancePoolFromAPI(&(*resp.JSON200.InstancePools)[i])
			cluster.c = c
			cluster.zone = zone

			list = append(list, cluster)
		}
	}

	return list, nil
}

// GetInstancePool returns the Instance Pool corresponding to the specified ID in the specified zone.
func (c *Client) GetInstancePool(ctx context.Context, zone, id string) (*InstancePool, error) {
	resp, err := c.GetInstancePoolWithResponse(apiv2.WithZone(ctx, zone), id)
	if err != nil {
		return nil, err
	}

	instancePool := instancePoolFromAPI(resp.JSON200)
	instancePool.c = c
	instancePool.zone = zone

	return instancePool, nil
}

// UpdateInstancePool updates the specified Instance Pool in the specified zone.
func (c *Client) UpdateInstancePool(ctx context.Context, zone string, instancePool *InstancePool) error {
	resp, err := c.UpdateInstancePoolWithResponse(
		apiv2.WithZone(ctx, zone),
		instancePool.ID,
		papi.UpdateInstancePoolJSONRequestBody{
			AntiAffinityGroups: func() *[]papi.AntiAffinityGroup {
				var list []papi.AntiAffinityGroup
				if l := len(instancePool.AntiAffinityGroupIDs); l > 0 {
					list = make([]papi.AntiAffinityGroup, l)
					for i, v := range instancePool.AntiAffinityGroupIDs {
						v := v
						list[i] = papi.AntiAffinityGroup{Id: &v}
					}
				}
				return &list
			}(),

			Description: &instancePool.Description,
			DiskSize:    &instancePool.DiskSize,

			ElasticIps: func() *[]papi.ElasticIp {
				var list []papi.ElasticIp
				if l := len(instancePool.ElasticIPIDs); l > 0 {
					list = make([]papi.ElasticIp, l)
					for i, v := range instancePool.ElasticIPIDs {
						v := v
						list[i] = papi.ElasticIp{Id: &v}
					}
				}
				return &list
			}(),

			InstanceType: &papi.InstanceType{Id: &instancePool.InstanceTypeID},
			Ipv6Enabled:  &instancePool.IPv6Enabled,
			Name:         &instancePool.Name,

			PrivateNetworks: func() *[]papi.PrivateNetwork {
				var list []papi.PrivateNetwork
				if l := len(instancePool.PrivateNetworkIDs); l > 0 {
					list = make([]papi.PrivateNetwork, l)
					for i, v := range instancePool.PrivateNetworkIDs {
						v := v
						list[i] = papi.PrivateNetwork{Id: &v}
					}
				}
				return &list
			}(),

			SecurityGroups: func() *[]papi.SecurityGroup {
				var list []papi.SecurityGroup
				if l := len(instancePool.SecurityGroupIDs); l > 0 {
					list = make([]papi.SecurityGroup, l)
					for i, v := range instancePool.SecurityGroupIDs {
						v := v
						list[i] = papi.SecurityGroup{Id: &v}
					}
				}
				return &list
			}(),

			SshKey: func() *papi.SshKey {
				if instancePool.SSHKey != "" {
					return &papi.SshKey{Name: &instancePool.SSHKey}
				}
				return nil
			}(),

			Template: func() *papi.Template {
				if instancePool.TemplateID != "" {
					return &papi.Template{Id: &instancePool.TemplateID}
				}
				return nil
			}(),

			UserData: func() *string {
				if instancePool.UserData != "" {
					return &instancePool.UserData
				}
				return nil
			}(),
		})
	if err != nil {
		return err
	}

	_, err = papi.NewPoller().
		WithTimeout(c.timeout).
		Poll(ctx, c.OperationPoller(zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}

// DeleteInstancePool deletes the specified Instance Pool in the specified zone.
func (c *Client) DeleteInstancePool(ctx context.Context, zone, id string) error {
	resp, err := c.DeleteInstancePoolWithResponse(apiv2.WithZone(ctx, zone), id)
	if err != nil {
		return err
	}

	_, err = papi.NewPoller().
		WithTimeout(c.timeout).
		Poll(ctx, c.OperationPoller(zone, *resp.JSON200.Id))
	if err != nil {
		return err
	}

	return nil
}
