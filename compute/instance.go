package compute

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/exoscale/egoscale/api"
	egoerr "github.com/exoscale/egoscale/error"
	egoapi "github.com/exoscale/egoscale/internal/egoscale"
)

// InstanceCreateOpts represents the Compute instance resource creation options.
type InstanceCreateOpts struct {
	// Name represents the Compute instance name.
	Name string
	// Type represents the Compute instance type.
	Type *InstanceType
	// Template represents the Compute instance template.
	Template *InstanceTemplate
	// VolumeSize represents the Compute instance storage volume size in GB.
	VolumeSize int64
	// AntiAffinityGroups represents the Compute instance Anti-Affinity Groups.
	AntiAffinityGroups []*AntiAffinityGroup
	// SecurityGroups represents the Compute instance Security Groups.
	SecurityGroups []*SecurityGroup
	// PrivateNetworks represents the Compute instance Private Networks.
	PrivateNetworks []*PrivateNetwork
	// SSHKey represents the Compute instance SSH key.
	SSHKey *SSHKey
	// EnableIPv6 indicates whether to enable IPv6 on the Compute instance public network interface.
	EnableIPv6 bool
	// UserData represents the Compute instance cloud-init User-Data content.
	UserData string
}

// InstanceUpdateOpts represents the Compute instance resource update options.
type InstanceUpdateOpts struct {
	// Name represents the Compute instance name.
	Name string
	// SecurityGroups represents the Compute instance Security Groups.
	SecurityGroups []*SecurityGroup
	// UserData represents the Compute instance cloud-init User-Data content.
	UserData string
}

// InstanceStartOpts represents the Compute instance start options.
type InstanceStartOpts struct {
	// RescueProfile represents the rescue profile to use when starting the Compute instance.
	RescueProfile string
}

// Instance represents a Compute instance resource.
type Instance struct {
	api.Resource

	ID          string
	Name        string
	Type        *InstanceType
	Template    *InstanceTemplate
	VolumeSize  int64
	IPv4Address net.IP
	IPv6Address net.IP
	SSHKey      *SSHKey
	Zone        *Zone

	volumeID string

	c *Client
}

func (i *Instance) String() string {
	return fmt.Sprintf("Instance(ID=%q, Name=%q)", i.ID, i.Name)
}

// AntiAffinityGroups returns the list of Anti-Affinity Groups the Compute instance is member of.
func (i *Instance) AntiAffinityGroups() ([]*AntiAffinityGroup, error) {
	res, err := i.c.c.ListWithContext(i.c.ctx, &egoapi.ListAffinityGroups{
		VirtualMachineID: egoapi.MustParseUUID(i.ID),
	})
	if err != nil {
		return nil, err
	}

	antiAffinityGroups := make([]*AntiAffinityGroup, 0)
	for _, r := range res {
		antiAffinityGroups = append(antiAffinityGroups, i.c.antiAffinityGroupFromAPI(r.(*egoapi.AffinityGroup)))
	}

	return antiAffinityGroups, nil
}

// ElasticIPs returns the list of Elastic IPs attached to the Compute instance.
func (i *Instance) ElasticIPs() ([]*ElasticIP, error) {
	res, err := i.c.c.ListWithContext(i.c.ctx, &egoapi.IPAddress{
		VirtualMachineID: egoapi.MustParseUUID(i.ID),
		IsElastic:        true,
	})
	if err != nil {
		return nil, err
	}

	elasticIPs := make([]*ElasticIP, 0)
	for _, r := range res {
		eip, err := i.c.elasticIPFromAPI(r.(*egoapi.IPAddress))
		if err != nil {
			return nil, err
		}

		elasticIPs = append(elasticIPs, eip)
	}

	return elasticIPs, nil
}

// AttachElasticIP attaches an Elastic IP to the Compute instance.
func (i *Instance) AttachElasticIP(elasticIP *ElasticIP) error {
	if elasticIP == nil {
		return errors.New("missing Elastic IP")
	}

	instanceNIC, err := i.defaultNIC()
	if err != nil {
		return err
	}

	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.AddIPToNic{
		NicID:     instanceNIC.ID,
		IPAddress: elasticIP.Address,
	}); err != nil {
		return i.c.csError(err)
	}

	return nil
}

// DetachElasticIP detaches an Elastic IP from the Compute instance.
func (i *Instance) DetachElasticIP(elasticIP *ElasticIP) error {
	if elasticIP == nil {
		return errors.New("missing Elastic IP")
	}

	instanceNIC, err := i.defaultNIC()
	if err != nil {
		return err
	}

	res, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.ListNics{NicID: instanceNIC.ID})
	if err != nil {
		i.c.csError(err)
	}

	for _, sip := range res.(*egoapi.ListNicsResponse).Nic[0].SecondaryIP {
		if sip.IPAddress.Equal(elasticIP.Address) {
			if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.RemoveIPFromNic{ID: sip.ID}); err != nil {
				return i.c.csError(err)
			}

			return nil
		}
	}

	return fmt.Errorf("Elastic IP address %q not attached to this Compute instance", elasticIP.Address)
}

// TODO func(i *Instance) PrivateNetworks()

// TODO func(i *Instance) AttachPrivateNetwork()

// TODO func(i *Instance) DetachPrivateNetwork()

// SecurityGroups returns the list of Security Groups the Compute instance is member of.
func (i *Instance) SecurityGroups() ([]*SecurityGroup, error) {
	res, err := i.c.c.ListWithContext(i.c.ctx, &egoapi.ListSecurityGroups{
		VirtualMachineID: egoapi.MustParseUUID(i.ID),
	})
	if err != nil {
		return nil, err
	}

	securityGroups := make([]*SecurityGroup, 0)
	for _, r := range res {
		securityGroups = append(securityGroups, i.c.securityGroupFromAPI(r.(*egoapi.SecurityGroup)))
	}

	return securityGroups, nil
}

// ReverseDNS returns the reverse DNS record currently set on the Compute instance public network interface IP address.
func (i *Instance) ReverseDNS() (string, error) {
	res, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.QueryReverseDNSForVirtualMachine{
		ID: egoapi.MustParseUUID(i.ID),
	})
	if err != nil {
		return "", i.c.csError(err)
	}

	reverseDNS := res.(*egoapi.VirtualMachine).DefaultNic().ReverseDNS
	if len(reverseDNS) == 0 {
		return "", nil
	}

	return reverseDNS[0].DomainName, nil
}

// SetReverseDNS sets the Compute instance public network interface IP address reverse DNS record.
func (i *Instance) SetReverseDNS(record string) error {
	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.UpdateReverseDNSForVirtualMachine{
		ID:         egoapi.MustParseUUID(i.ID),
		DomainName: record,
	}); err != nil {
		return i.c.csError(err)
	}

	return nil
}

// UnsetReverseDNS unsets the reverse DNS record currently set on the Compute instance public network interface IP
// address.
func (i *Instance) UnsetReverseDNS() error {
	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.DeleteReverseDNSFromVirtualMachine{
		ID: egoapi.MustParseUUID(i.ID),
	}); err != nil {
		return i.c.csError(err)
	}

	return nil
}

// TODO func(i *Instance) ResizeVolume()

// TODO func(i *Instance) SnapshotVolume()

// TODO func(i *Instance) VolumeSnapshots()

// State returns the Compute instance current state.
func (i *Instance) State() (string, error) {
	res, err := i.c.c.GetWithContext(i.c.ctx, &egoapi.ListVirtualMachines{
		ZoneID: egoapi.MustParseUUID(i.Zone.ID),
		ID:     egoapi.MustParseUUID(i.ID),
	})
	if err != nil {
		return "", err
	}

	return strings.ToLower(res.(*egoapi.VirtualMachine).State), nil
}

// Update updates the Compute instance properties.
func (i *Instance) Update(opts *InstanceUpdateOpts) error {
	var req egoapi.UpdateVirtualMachine

	if opts == nil {
		opts = new(InstanceUpdateOpts)
	}

	req.ID = egoapi.MustParseUUID(i.ID)
	req.Name = opts.Name
	req.DisplayName = opts.Name

	if opts.UserData != "" {
		req.UserData = base64.StdEncoding.EncodeToString([]byte(opts.UserData))
	}

	if opts.SecurityGroups != nil {
		req.SecurityGroupIDs = func(items []*SecurityGroup) []egoapi.UUID {
			ids := make([]egoapi.UUID, len(items))
			for i := range items {
				ids[i] = *egoapi.MustParseUUID(items[i].ID)
			}
			return ids
		}(opts.SecurityGroups)
	}

	if _, err := i.c.c.RequestWithContext(i.c.ctx, &req); err != nil {
		return err
	}

	if opts.Name != "" {
		i.Name = opts.Name
	}

	return nil
}

// Scale changes the Compute instance type.
func (i *Instance) Scale(it *InstanceType) error {
	if it == nil {
		return errors.New("missing instance type parameter")
	}

	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.ScaleVirtualMachine{
		ID:                egoapi.MustParseUUID(i.ID),
		ServiceOfferingID: egoapi.MustParseUUID(it.ID),
	}); err != nil {
		return err
	}

	i.Type = it

	return nil
}

// Start starts a stopped Compute instance.
func (i *Instance) Start(opts *InstanceStartOpts) error {
	if opts == nil {
		opts = new(InstanceStartOpts)
	}

	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.StartVirtualMachine{
		ID:            egoapi.MustParseUUID(i.ID),
		RescueProfile: opts.RescueProfile,
	}); err != nil {
		return err
	}

	return nil
}

// Stop stops a running Compute instance.
func (i *Instance) Stop() error {
	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.StopVirtualMachine{
		ID: egoapi.MustParseUUID(i.ID),
	}); err != nil {
		return err
	}

	return nil
}

// Reboot reboots the Compute instance.
func (i *Instance) Reboot() error {
	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.RebootVirtualMachine{
		ID: egoapi.MustParseUUID(i.ID),
	}); err != nil {
		return err
	}

	return nil
}

// Delete deletes the Compute instance.
func (i *Instance) Delete() error {
	if _, err := i.c.c.RequestWithContext(i.c.ctx, &egoapi.DestroyVirtualMachine{
		ID: egoapi.MustParseUUID(i.ID)},
	); err != nil {
		return i.c.csError(err)
	}

	i.ID = ""
	i.Name = ""
	i.Zone = nil
	i.Type = nil
	i.Template = nil
	i.VolumeSize = 0
	i.IPv4Address = nil
	i.IPv6Address = nil
	i.SSHKey = nil

	return nil
}

func (i *Instance) defaultNIC() (*egoapi.Nic, error) {
	var instance egoapi.VirtualMachine

	if err := json.Unmarshal(i.Raw(), &instance); err != nil {
		return nil, err
	}

	return instance.DefaultNic(), nil
}

// CreateInstance creates a new Compute instance resource.
func (c *Client) CreateInstance(zone *Zone, opts *InstanceCreateOpts) (*Instance, error) {
	var req egoapi.DeployVirtualMachine

	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}
	req.ZoneID = egoapi.MustParseUUID(zone.ID)

	if opts == nil {
		opts = new(InstanceCreateOpts)
	}

	req.Name = opts.Name
	req.DisplayName = opts.Name
	req.RootDiskSize = opts.VolumeSize
	req.IP6 = &opts.EnableIPv6

	if opts.Type != nil {
		req.ServiceOfferingID = egoapi.MustParseUUID(opts.Type.ID)
	}

	if opts.Template != nil {
		req.TemplateID = egoapi.MustParseUUID(opts.Template.ID)
	}

	if opts.AntiAffinityGroups != nil {
		req.AffinityGroupIDs = func(groups []*AntiAffinityGroup) []egoapi.UUID {
			ids := make([]egoapi.UUID, len(groups))
			for i := range groups {
				ids[i] = *egoapi.MustParseUUID(groups[i].ID)
			}
			return ids
		}(opts.AntiAffinityGroups)
	}

	if opts.SecurityGroups != nil {
		req.SecurityGroupIDs = func(items []*SecurityGroup) []egoapi.UUID {
			ids := make([]egoapi.UUID, len(items))
			for i := range items {
				ids[i] = *egoapi.MustParseUUID(items[i].ID)
			}
			return ids
		}(opts.SecurityGroups)
	}

	if opts.PrivateNetworks != nil {
		req.NetworkIDs = func(items []*PrivateNetwork) []egoapi.UUID {
			ids := make([]egoapi.UUID, len(items))
			for i := range items {
				ids[i] = *egoapi.MustParseUUID(items[i].ID)
			}
			return ids
		}(opts.PrivateNetworks)
	}

	if opts.SSHKey != nil {
		req.KeyPair = opts.SSHKey.Name
	}

	if opts.UserData != "" {
		req.UserData = base64.StdEncoding.EncodeToString([]byte(opts.UserData))
	}

	res, err := c.c.RequestWithContext(c.ctx, &req)
	if err != nil {
		return nil, c.csError(err)
	}

	return c.instanceFromAPI(res.(*egoapi.VirtualMachine))
}

// ListInstances returns the list of Compute instances.
func (c *Client) ListInstances(zone *Zone) ([]*Instance, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	res, err := c.c.ListWithContext(c.ctx, &egoapi.VirtualMachine{ZoneID: egoapi.MustParseUUID(zone.ID)})
	if err != nil {
		return nil, err
	}

	instances := make([]*Instance, 0)
	for _, i := range res {
		instance, err := c.instanceFromAPI(i.(*egoapi.VirtualMachine))
		if err != nil {
			return nil, err
		}

		instances = append(instances, instance)
	}

	return instances, nil
}

// GetInstanceByID returns a Compute instance by its unique identifier.
func (c *Client) GetInstanceByID(zone *Zone, id string) (*Instance, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	instanceID, err := egoapi.ParseUUID(id)
	if err != nil {
		return nil, err
	}

	return c.getInstance(egoapi.MustParseUUID(zone.ID), instanceID, nil)
}

// GetInstanceByAddress returns a Compute instance by its primary IP address.
func (c *Client) GetInstanceByAddress(zone *Zone, address string) (*Instance, error) {
	if zone == nil {
		return nil, egoerr.ErrMissingZone
	}

	return c.getInstance(egoapi.MustParseUUID(zone.ID), nil, net.ParseIP(address))
}

func (c *Client) getInstance(zoneID, id *egoapi.UUID, address net.IP) (*Instance, error) {
	res, err := c.c.ListWithContext(c.ctx, &egoapi.ListVirtualMachines{
		ZoneID:    zoneID,
		ID:        id,
		IPAddress: address,
	})
	if err != nil {
		return nil, c.csError(err)
	}

	if len(res) == 0 {
		return nil, egoerr.ErrResourceNotFound
	}

	return c.instanceFromAPI(res[0].(*egoapi.VirtualMachine))
}

func (c *Client) instanceFromAPI(i *egoapi.VirtualMachine) (*Instance, error) {
	var sshKey *SSHKey

	if i.DefaultNic() == nil {
		return nil, fmt.Errorf("default NIC missing")
	}

	zone, err := c.GetZoneByID(i.ZoneID.String())
	if err != nil {
		return nil, err
	}

	typ, err := c.GetInstanceTypeByID(i.ServiceOfferingID.String())
	if err != nil {
		return nil, err
	}

	template, err := c.GetInstanceTemplate(zone, i.TemplateID.String(), "")
	if err != nil {
		return nil, err
	}

	if i.KeyPair != "" {
		if sshKey, err = c.GetSSHKey(i.KeyPair); err != nil {
			return nil, err
		}
	}

	res, err := c.c.RequestWithContext(c.ctx, &egoapi.ListVolumes{VirtualMachineID: i.ID})
	if err != nil {
		return nil, c.csError(err)
	}
	volumeSize := int64(res.(*egoapi.ListVolumesResponse).Volume[0].Size)

	return &Instance{
		Resource:    api.MarshalResource(i),
		ID:          i.ID.String(),
		Name:        i.Name,
		Zone:        zone,
		Type:        typ,
		Template:    template,
		VolumeSize:  int64(volumeSize),
		IPv4Address: i.DefaultNic().IPAddress,
		IPv6Address: i.DefaultNic().IP6Address,
		SSHKey:      sshKey,
		c:           c,
	}, nil
}
