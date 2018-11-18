package egoscale

import (
	"net"
	"net/url"
	"testing"
)

func TestVirtualMachine(t *testing.T) {
	instance := &VirtualMachine{}
	if instance.ResourceType() != "UserVM" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestDeployVirtualMachine(t *testing.T) {
	req := &DeployVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestDestroyVirtualMachine(t *testing.T) {
	req := &DestroyVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestRebootVirtualMachine(t *testing.T) {
	req := &RebootVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestStartVirtualMachine(t *testing.T) {
	req := &StartVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestStopVirtualMachine(t *testing.T) {
	req := &StopVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestResetPasswordForVirtualMachine(t *testing.T) {
	req := &ResetPasswordForVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestUpdateVirtualMachine(t *testing.T) {
	req := &UpdateVirtualMachine{}
	_ = req.response().(*VirtualMachine)
}

func TestListVirtualMachines(t *testing.T) {
	req := &ListVirtualMachines{}
	_ = req.response().(*ListVirtualMachinesResponse)
}

func TestGetVMPassword(t *testing.T) {
	req := &GetVMPassword{}
	_ = req.response().(*Password)
}

func TestRestoreVirtualMachine(t *testing.T) {
	req := &RestoreVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestChangeServiceForVirtualMachine(t *testing.T) {
	req := &ChangeServiceForVirtualMachine{}
	_ = req.response().(*VirtualMachine)
}

func TestScaleVirtualMachine(t *testing.T) {
	req := &ScaleVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestRecoverVirtualMachine(t *testing.T) {
	req := &RecoverVirtualMachine{}
	_ = req.response().(*VirtualMachine)
}

func TestExpungeVirtualMachine(t *testing.T) {
	req := &ExpungeVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestGetVirtualMachineUserData(t *testing.T) {
	req := &GetVirtualMachineUserData{}
	_ = req.response().(*VirtualMachineUserData)
}

func TestAddNicToVirtualMachine(t *testing.T) {
	req := &AddNicToVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestRemoveNicFromVirtualMachine(t *testing.T) {
	req := &RemoveNicFromVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestUpdateDefaultNicForVirtualMachine(t *testing.T) {
	req := &UpdateDefaultNicForVirtualMachine{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestUpdateVMNicIP(t *testing.T) {
	req := &UpdateVMNicIP{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*VirtualMachine)
}

func TestDeployOnBeforeSend(t *testing.T) {
	req := &DeployVirtualMachine{
		SecurityGroupNames: []string{"default"},
	}
	params := url.Values{}

	if err := req.onBeforeSend(params); err != nil {
		t.Error(err)
	}
}

func TestDeployOnBeforeSendNoSG(t *testing.T) {
	req := &DeployVirtualMachine{}
	params := url.Values{}

	// CS will pick the default oiine
	if err := req.onBeforeSend(params); err != nil {
		t.Error(err)
	}
}

func TestDeployOnBeforeSendBothSG(t *testing.T) {
	req := &DeployVirtualMachine{
		SecurityGroupIDs:   []UUID{*MustParseUUID("f2b4e439-2b23-441c-ba66-0e25cdfe1b2b")},
		SecurityGroupNames: []string{"foo"},
	}
	params := url.Values{}

	if err := req.onBeforeSend(params); err == nil {
		t.Errorf("DeployVM should only accept SG ids or names")
	}
}

func TestDeployOnBeforeSendBothAG(t *testing.T) {
	req := &DeployVirtualMachine{
		AffinityGroupIDs:   []UUID{*MustParseUUID("f2b4e439-2b23-441c-ba66-0e25cdfe1b2b")},
		AffinityGroupNames: []string{"foo"},
	}
	params := url.Values{}

	if err := req.onBeforeSend(params); err == nil {
		t.Errorf("DeployVM should only accept SG ids or names")
	}
}

func TestNicHelpers(t *testing.T) {
	vm := &VirtualMachine{
		ID: MustParseUUID("25ce0763-f34d-435a-8b84-08466908355a"),
		Nic: []Nic{
			{
				ID:           MustParseUUID("e3b9c165-f3c3-4672-be54-08bfa6bac6fe"),
				IsDefault:    true,
				MACAddress:   MustParseMAC("06:aa:14:00:00:18"),
				IPAddress:    net.ParseIP("192.168.0.10"),
				Gateway:      net.ParseIP("192.168.0.1"),
				Netmask:      net.ParseIP("255.255.255.0"),
				NetworkID:    MustParseUUID("d48bfccc-c11f-438f-8177-9cf6a40dc4d8"),
				NetworkName:  "defaultGuestNetwork",
				BroadcastURI: "vlan://untagged",
				TrafficType:  "Guest",
				Type:         "Shared",
			}, {
				BroadcastURI: "vxlan://001",
				ID:           MustParseUUID("10b8ffc8-62b3-4b87-82d0-fb7f31bc99b6"),
				IsDefault:    false,
				MACAddress:   MustParseMAC("0a:7b:5e:00:25:fa"),
				NetworkID:    MustParseUUID("5f1033fe-2abd-4dda-80b6-c946e21a78ec"),
				NetworkName:  "privNetForBasicZone1",
				TrafficType:  "Guest",
				Type:         "Isolated",
			}, {
				BroadcastURI: "vxlan://002",
				ID:           MustParseUUID("10b8ffc8-62b3-4b87-82d0-fb7f31bc99b7"),
				IsDefault:    false,
				MACAddress:   MustParseMAC("0a:7b:5e:00:25:ff"),
				NetworkID:    MustParseUUID("5f1033fe-2abd-4dda-80b6-c946e21a72ec"),
				NetworkName:  "privNetForBasicZone2",
				TrafficType:  "Guest",
				Type:         "Isolated",
			},
		},
	}

	nic := vm.DefaultNic()
	if nic.IPAddress.String() != "192.168.0.10" {
		t.Errorf("Default NIC doesn't match")
	}

	ip := vm.IP()
	if ip.String() != "192.168.0.10" {
		t.Errorf("IP Address doesn't match")
	}

	nic1 := vm.NicByID(*MustParseUUID("e3b9c165-f3c3-4672-be54-08bfa6bac6fe"))
	if nic1.ID != nil && !nic.ID.Equal(*nic1.ID) {
		t.Errorf("NicByID does not match %#v %#v", nic, nic1)
	}

	if len(vm.NicsByType("Isolated")) != 2 {
		t.Errorf("Isolated nics count does not match")
	}

	if len(vm.NicsByType("Shared")) != 1 {
		t.Errorf("Shared nics count does not match")
	}

	if len(vm.NicsByType("Dummy")) != 0 {
		t.Errorf("Dummy nics count does not match")
	}

	if vm.NicByNetworkID(*MustParseUUID("5f1033fe-2abd-4dda-80b6-c946e21a78ec")) == nil {
		t.Errorf("NetworkID nic wasn't found")
	}

	if vm.NicByNetworkID(*MustParseUUID("5f1033fe-2abd-4dda-80b6-c946e21a78ed")) != nil {
		t.Errorf("NetworkID nic was found??")
	}
}

func TestNicNoDefault(t *testing.T) {
	vm := &VirtualMachine{
		Nic: []Nic{},
	}

	// code coverage...
	nic := vm.DefaultNic()
	if nic != nil {
		t.Errorf("Default NIC wasn't nil?")
	}
}

func TestUserDataDecode(t *testing.T) {
	userDatas := []VirtualMachineUserData{{
		UserData: "aGVsbG8h",
	}, {
		UserData: "H4sIAEd08VsC/8tIzcnJVwQAYMmGmgYAAAA=",
	}}

	expected := "hello!"
	for _, tt := range userDatas {
		if output, _ := tt.Decode(); output != "hello!" {
			t.Errorf("bad userdata decoding, want: %q, got: %q", expected, output)
		}
	}
}
