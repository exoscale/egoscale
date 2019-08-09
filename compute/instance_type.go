package compute

import (
	"fmt"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// InstanceType represents an Exoscale Compute instance type.
type InstanceType struct {
	api.Resource

	ID     string
	Name   string
	CPU    int
	Memory int

	c *Client
}

func (t *InstanceType) String() string {
	return fmt.Sprintf("InstanceType(ID=%q, Name=%q)", t.ID, t.Name)
}

// ListInstanceTypes returns the list of available Exoscale Compute instance types.
func (c *Client) ListInstanceTypes() ([]*InstanceType, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.ServiceOffering{})
	if err != nil {
		return nil, err
	}

	instanceTypes := make([]*InstanceType, 0)
	for _, i := range res {
		instanceTypes = append(instanceTypes, c.instanceTypeFromAPI(i.(*egoapi.ServiceOffering)))
	}

	return instanceTypes, nil
}

// GetInstanceType returns an Exoscale Compute instance type by its name.
func (c *Client) GetInstanceTypeByName(name string) (*InstanceType, error) {
	return c.getInstanceType(nil, name)
}

// GetInstanceType returns an Exoscale Compute instance type by its unique identifier.
func (c *Client) GetInstanceTypeByID(id string) (*InstanceType, error) {
	instanceTypeID, err := egoapi.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return c.getInstanceType(instanceTypeID, "")
}

func (c *Client) getInstanceType(id *egoapi.UUID, name string) (*InstanceType, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.ServiceOffering{
		ID:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.instanceTypeFromAPI(res[0].(*egoapi.ServiceOffering)), nil
}

func (c *Client) instanceTypeFromAPI(so *egoapi.ServiceOffering) *InstanceType {
	return &InstanceType{
		Resource: api.MarshalResource(so),
		ID:       so.ID.String(),
		Name:     so.Name,
		CPU:      so.CPUNumber,
		Memory:   so.Memory,
		c:        c,
	}
}
