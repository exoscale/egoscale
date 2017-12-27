package egoscale

func (exo *Client) GetAllSecurityGroups() (map[string]SecurityGroup, error) {
	var sgs map[string]SecurityGroup
	securityGroups, err := exo.ListSecurityGroups(&ListSecurityGroupsRequest{})

	if err != nil {
		return nil, err
	}

	sgs = make(map[string]SecurityGroup)
	for _, sg := range securityGroups {
		sgs[sg.Name] = *sg
	}
	return sgs, nil
}

func (exo *Client) GetSecurityGroupId(name string) (string, error) {
	securityGroups, err := exo.ListSecurityGroups(&ListSecurityGroupsRequest{
		SecurityGroupName: name,
	})
	if err != nil {
		return "", err
	}

	for _, sg := range securityGroups {
		if sg.Name == name {
			return sg.Id, nil
		}
	}

	return "", nil
}

func (exo *Client) GetKeypairs() ([]SshKeyPair, error) {
	var keypairs []SshKeyPair

	keys, err := exo.ListSshKeyPairs(&ListSshKeyPairsRequest{})
	if err != nil {
		return keypairs, err
	}
	keypairs = make([]SshKeyPair, len(keys))
	for i, keypair := range keys {
		keypairs[i] = *keypair
	}
	return keypairs, nil
}

func (exo *Client) GetAffinityGroups() (map[string]string, error) {
	var affinitygroups map[string]string
	groups, err := exo.ListAffinityGroups(&ListAffinityGroupsRequest{})
	if err != nil {
		return affinitygroups, err
	}

	affinitygroups = make(map[string]string)
	for _, affinitygroup := range groups {
		affinitygroups[affinitygroup.Name] = affinitygroup.Id
	}
	return affinitygroups, nil
}

func (exo *Client) GetTopology() (*Topology, error) {
	zones, err := exo.GetAllZones()
	if err != nil {
		return nil, err
	}
	images, err := exo.GetImages()
	if err != nil {
		return nil, err
	}
	securityGroups, err := exo.GetAllSecurityGroups()
	if err != nil {
		return nil, err
	}
	groups := make(map[string]string)
	for k, v := range securityGroups {
		groups[k] = v.Id
	}

	keypairs, err := exo.GetKeypairs()
	if err != nil {
		return nil, err
	}

	/* Convert the ssh keypair to contain just the name */
	keynames := make([]string, len(keypairs))
	for i, k := range keypairs {
		keynames[i] = k.Name
	}

	affinitygroups, err := exo.GetAffinityGroups()
	if err != nil {
		return nil, err
	}

	profiles, err := exo.GetProfiles()
	if err != nil {
		return nil, err
	}

	topo := &Topology{
		Zones:          zones,
		Images:         images,
		Keypairs:       keynames,
		Profiles:       profiles,
		AffinityGroups: affinitygroups,
		SecurityGroups: groups,
	}

	return topo, nil
}
