package cmd

import (
	"github.com/exoscale/egoscale"
)

type category struct {
	name  string
	alias []string
	doc   string
	cmd   []cmd
}

type cmd struct {
	command egoscale.Command
	name    string
}

var methods = []category{
	category{
		"network",
		[]string{"net"},
		"doc net",
		[]cmd{
			{&egoscale.CreateNetwork{}, "create"},
			{&egoscale.DeleteNetwork{}, "delete"},
			{&egoscale.ListNetworkOfferings{}, ""},
			{&egoscale.ListNetworks{}, "list"},
			{&egoscale.UpdateNetwork{}, "update"},
		},
	},
	category{
		"vm",
		[]string{"virtual-machine"},
		"doc vm",
		[]cmd{
			{&egoscale.AddNicToVirtualMachine{}, "addNic"},
			{&egoscale.ChangeServiceForVirtualMachine{}, "changeService"},
			{&egoscale.DeployVirtualMachine{}, "deploy"},
			{&egoscale.DestroyVirtualMachine{}, "destroy"},
			{&egoscale.ExpungeVirtualMachine{}, "expunge"},
			{&egoscale.GetVMPassword{}, "getPassword"},
			{&egoscale.GetVirtualMachineUserData{}, "getUserData"},
			{&egoscale.ListVirtualMachines{}, "list"},
			{&egoscale.MigrateVirtualMachine{}, ""},
			{&egoscale.RebootVirtualMachine{}, "reboot"},
			{&egoscale.RecoverVirtualMachine{}, "recover"},
			{&egoscale.RemoveNicFromVirtualMachine{}, "removeNic"},
			{&egoscale.ResetPasswordForVirtualMachine{}, "resetPassword"},
			{&egoscale.RestoreVirtualMachine{}, "restore"},
			{&egoscale.ScaleVirtualMachine{}, "scale"},
			{&egoscale.StartVirtualMachine{}, "start"},
			{&egoscale.StopVirtualMachine{}, "stop"},
			{&egoscale.UpdateVirtualMachine{}, "update"},
			{&egoscale.CreateAffinityGroup{}, "createAffinityGroup"},
			{&egoscale.DeleteAffinityGroup{}, "deleteAffinityGroup"},
			{&egoscale.ListAffinityGroups{}, "listAffinityGroup"},
			{&egoscale.UpdateVMAffinityGroup{}, ""},
			{&egoscale.DeleteReverseDNSFromVirtualMachine{}, "deleteReverseDNSFromVM"},
			{&egoscale.QueryReverseDNSForVirtualMachine{}, "queryReverseDNSForVM"},
			{&egoscale.UpdateReverseDNSForVirtualMachine{}, "updateReverseDNSForVM"},
		},
	},
	category{
		"volume",
		[]string{"vol"},
		"doc volume",
		[]cmd{
			{&egoscale.ListVolumes{}, "list"},
			{&egoscale.ResizeVolume{}, "resize"},
		},
	},
	category{
		"template",
		[]string{"temp"},
		"doc template",
		[]cmd{
			{&egoscale.ListTemplates{}, "list"},
		},
	},
	category{
		"account",
		[]string{"acc"},
		"doc account",
		[]cmd{
			{&egoscale.ListAccounts{}, "list"},
		},
	},
	category{
		"zone",
		nil,
		"doc zone",
		[]cmd{
			{&egoscale.ListZones{}, "list"},
		},
	},
	category{
		"snapshot",
		[]string{"snap"},
		"doc snapshot",
		[]cmd{
			{&egoscale.CreateSnapshot{}, "create"},
			{&egoscale.DeleteSnapshot{}, "delete"},
			{&egoscale.ListSnapshots{}, "list"},
			{&egoscale.RevertSnapshot{}, "revert"},
		},
	},
	category{
		"user",
		[]string{"usr"},
		"doc user",
		[]cmd{
			{&egoscale.ListUsers{}, "list"},
			{&egoscale.RegisterUserKeys{}, ""},
		},
	},
	category{
		"security-group",
		[]string{"sg"},
		"doc security-group",
		[]cmd{
			{&egoscale.AuthorizeSecurityGroupEgress{}, "authorizeEgress"},
			{&egoscale.AuthorizeSecurityGroupIngress{}, "authorizeIngress"},
			{&egoscale.CreateSecurityGroup{}, "create"},
			{&egoscale.DeleteSecurityGroup{}, "delete"},
			{&egoscale.ListSecurityGroups{}, "list"},
			{&egoscale.RevokeSecurityGroupEgress{}, "revokeEgress"},
			{&egoscale.RevokeSecurityGroupIngress{}, "revokeIngress"},
		},
	},
	category{
		"ssh",
		nil,
		"doc ssh",
		[]cmd{
			{&egoscale.RegisterSSHKeyPair{}, "register"},
			{&egoscale.ListSSHKeyPairs{}, "list"},
			{&egoscale.CreateSSHKeyPair{}, "create"},
			{&egoscale.DeleteSSHKeyPair{}, "delete"},
			{&egoscale.ResetSSHKeyForVirtualMachine{}, "reset"},
		},
	},
	category{
		"vm-group",
		[]string{"vg"},
		"doc vm-group",
		[]cmd{
			{&egoscale.CreateInstanceGroup{}, "create"},
			{&egoscale.ListInstanceGroups{}, "list"},
		},
	},
	category{
		"tags",
		nil,
		"doc tags",
		[]cmd{
			{&egoscale.CreateTags{}, "create"},
			{&egoscale.DeleteTags{}, "delete"},
			{&egoscale.ListTags{}, "list"},
		},
	},
	category{
		"nic",
		nil,
		"doc nic",
		[]cmd{
			{&egoscale.ActivateIP6{}, ""},
			{&egoscale.AddIPToNic{}, ""},
			{&egoscale.ListNics{}, "list"},
			{&egoscale.RemoveIPFromNic{}, ""},
		},
	},
	category{
		"address",
		[]string{"addr"},
		"doc address",
		[]cmd{
			{&egoscale.AssociateIPAddress{}, "associate"},
			{&egoscale.DisassociateIPAddress{}, "disassociate"},
			{&egoscale.ListPublicIPAddresses{}, "list"},
			{&egoscale.UpdateIPAddress{}, "update"},
			{&egoscale.DeleteReverseDNSFromPublicIPAddress{}, "deleteReverseDNSFromAddress"},
			{&egoscale.QueryReverseDNSForPublicIPAddress{}, "queryReverseDNSForAddress"},
			{&egoscale.UpdateReverseDNSForPublicIPAddress{}, "updateReverseDNSForAddress"},
		},
	},
	category{
		"async-job",
		[]string{"aj"},
		"doc async-job",
		[]cmd{
			{&egoscale.QueryAsyncJobResult{}, ""},
			{&egoscale.ListAsyncJobs{}, ""},
		},
	},
	category{
		"apis",
		nil,
		"doc apis",
		[]cmd{
			{&egoscale.ListAPIs{}, "list"},
		},
	},
	category{
		"event",
		nil,
		"doc event",
		[]cmd{
			{&egoscale.ListEventTypes{}, "listType"},
			{&egoscale.ListEvents{}, "list"},
		},
	},
	category{
		"offerings",
		nil,
		"doc offerings",
		[]cmd{
			{&egoscale.ListResourceDetails{}, "listDetails"},
			{&egoscale.ListResourceLimits{}, "listLimits"},
			{&egoscale.ListServiceOfferings{}, "list"},
		},
	},
}
