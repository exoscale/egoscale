package cmd

import (
	"github.com/exoscale/egoscale"
)

type cmd struct {
	command egoscale.Command
	hidden  bool
	name    string
}

var methods = map[string][]cmd{
	"network": {
		{&egoscale.CreateNetwork{}, false, "create"},
		{&egoscale.DeleteNetwork{}, false, "delete"},
		{&egoscale.ListNetworkOfferings{}, false, ""},
		{&egoscale.ListNetworks{}, false, "list"},
		{&egoscale.RestartNetwork{}, true, "restart"},
		{&egoscale.UpdateNetwork{}, false, "update"},
	},
	"vm": {
		{&egoscale.AddNicToVirtualMachine{}, false, "addNic"},
		{&egoscale.ChangeServiceForVirtualMachine{}, false, "changeService"},
		{&egoscale.DeployVirtualMachine{}, false, "deploy"},
		{&egoscale.DestroyVirtualMachine{}, false, "destroy"},
		{&egoscale.ExpungeVirtualMachine{}, false, "expunge"},
		{&egoscale.GetVMPassword{}, false, "getPassword"},
		{&egoscale.GetVirtualMachineUserData{}, false, "getUserData"},
		{&egoscale.ListVirtualMachines{}, false, "list"},
		{&egoscale.MigrateVirtualMachine{}, true, ""},
		{&egoscale.RebootVirtualMachine{}, false, "reboot"},
		{&egoscale.RecoverVirtualMachine{}, false, "recover"},
		{&egoscale.RemoveNicFromVirtualMachine{}, false, "removeNic"},
		{&egoscale.ResetPasswordForVirtualMachine{}, false, "resetPassword"},
		{&egoscale.RestoreVirtualMachine{}, false, "restore"},
		{&egoscale.ScaleVirtualMachine{}, false, "scale"},
		{&egoscale.StartVirtualMachine{}, false, "start"},
		{&egoscale.StopVirtualMachine{}, false, "stop"},
		{&egoscale.UpdateDefaultNicForVirtualMachine{}, true, ""},
		{&egoscale.UpdateVirtualMachine{}, false, "update"},
	},
	"volume": {
		{&egoscale.ListVolumes{}, false, "list"},
		{&egoscale.ResizeVolume{}, false, "resize"},
	},
	"template": {
		{&egoscale.CopyTemplate{}, true, ""},
		{&egoscale.CreateTemplate{}, true, ""},
		{&egoscale.ListTemplates{}, false, "list"},
		{&egoscale.PrepareTemplate{}, true, ""},
		{&egoscale.RegisterTemplate{}, true, ""},
		{&egoscale.ListOSCategories{}, true, ""},
	},
	"account": {
		{&egoscale.EnableAccount{}, true, ""},
		{&egoscale.DisableAccount{}, true, ""},
		{&egoscale.ListAccounts{}, false, "list"},
	},
	"zone": {
		{&egoscale.ListZones{}, false, "list"},
	},
	"snapshot": {
		{&egoscale.CreateSnapshot{}, false, "create"},
		{&egoscale.DeleteSnapshot{}, false, "delete"},
		{&egoscale.ListSnapshots{}, false, "list"},
		{&egoscale.RevertSnapshot{}, false, "revert"},
	},
	"user": {
		{&egoscale.CreateUser{}, true, ""},
		{&egoscale.DeleteUser{}, true, ""},
		//{&egoscale.DisableUser{}, true},
		//{&egoscale.GetUser{}, true},
		{&egoscale.UpdateUser{}, true, ""},
		{&egoscale.ListUsers{}, false, "list"},
		{&egoscale.RegisterUserKeys{}, false, ""},
	},
	"security-group,sg": {
		{&egoscale.AuthorizeSecurityGroupEgress{}, false, "authorizeEgress"},
		{&egoscale.AuthorizeSecurityGroupIngress{}, false, "authorizeIngress"},
		{&egoscale.CreateSecurityGroup{}, false, "create"},
		{&egoscale.DeleteSecurityGroup{}, false, "delete"},
		{&egoscale.ListSecurityGroups{}, false, "list"},
		{&egoscale.RevokeSecurityGroupEgress{}, false, "RevokeEgress"},
		{&egoscale.RevokeSecurityGroupIngress{}, false, "RevokeIngress"},
	},
	"ssh": {
		{&egoscale.RegisterSSHKeyPair{}, false, "register"},
		{&egoscale.ListSSHKeyPairs{}, false, "list"},
		{&egoscale.CreateSSHKeyPair{}, false, "create"},
		{&egoscale.DeleteSSHKeyPair{}, false, "delete"},
		{&egoscale.ResetSSHKeyForVirtualMachine{}, false, "reset"},
	},
	"affinity-group,ag": {
		{&egoscale.CreateAffinityGroup{}, false, "create"},
		{&egoscale.DeleteAffinityGroup{}, false, "delete"},
		{&egoscale.ListAffinityGroups{}, false, "list"},
		{&egoscale.UpdateVMAffinityGroup{}, false, ""},
	},
	"vm-group,vg": {
		{&egoscale.CreateInstanceGroup{}, false, "create"},
		{&egoscale.ListInstanceGroups{}, false, "list"},
	},
	"tags": {
		{&egoscale.CreateTags{}, false, "create"},
		{&egoscale.DeleteTags{}, false, "delete"},
		{&egoscale.ListTags{}, false, "list"},
	},
	"nic": {
		{&egoscale.ActivateIP6{}, false, ""},
		{&egoscale.AddIPToNic{}, false, ""},
		{&egoscale.ListNics{}, false, "list"},
		{&egoscale.RemoveIPFromNic{}, false, ""},
	},
	"address": {
		{&egoscale.AssociateIPAddress{}, false, "associate"},
		{&egoscale.DisassociateIPAddress{}, false, "disassociate"},
		{&egoscale.ListPublicIPAddresses{}, false, "list"},
		{&egoscale.UpdateIPAddress{}, false, "update"},
	},
	"async-job,aj": {
		{&egoscale.QueryAsyncJobResult{}, false, ""},
		{&egoscale.ListAsyncJobs{}, false, ""},
	},
	"apis": {
		{&egoscale.ListAPIs{}, false, "list"},
	},
	"event": {
		{&egoscale.ListEventTypes{}, false, "listType"},
		{&egoscale.ListEvents{}, false, "list"},
	},
	"offerings": {
		{&egoscale.ListResourceDetails{}, false, "listDetails"},
		{&egoscale.ListResourceLimits{}, false, "listLimits"},
		{&egoscale.ListServiceOfferings{}, false, "list"},
	},
	"host": {
		{&egoscale.ListHosts{}, true, ""},
		{&egoscale.UpdateHost{}, true, ""},
	},
	"reversedns": {
		{&egoscale.DeleteReverseDNSFromPublicIPAddress{}, false, "deleteFromAddress"},
		{&egoscale.QueryReverseDNSForPublicIPAddress{}, false, "queryForAddress"},
		{&egoscale.UpdateReverseDNSForPublicIPAddress{}, false, "updateForAddress"},
		{&egoscale.DeleteReverseDNSFromVirtualMachine{}, false, "deleteFromVM"},
		{&egoscale.QueryReverseDNSForVirtualMachine{}, false, "queryForVM"},
		{&egoscale.UpdateReverseDNSForVirtualMachine{}, false, "updateForVM"},
	},
}
