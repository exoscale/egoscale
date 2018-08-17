package cmd

import (
	"github.com/exoscale/egoscale"
)

type cmd struct {
	command egoscale.Command
	name    string
}

var methods = map[string][]cmd{
	"network": {
		{&egoscale.CreateNetwork{}, "create"},
		{&egoscale.DeleteNetwork{}, "delete"},
		{&egoscale.ListNetworkOfferings{}, ""},
		{&egoscale.ListNetworks{}, "list"},
		{&egoscale.RestartNetwork{}, "restart"},
		{&egoscale.UpdateNetwork{}, "update"},
	},
	"vm": {
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
		{&egoscale.UpdateDefaultNicForVirtualMachine{}, ""},
		{&egoscale.UpdateVirtualMachine{}, "update"},
	},
	"volume": {
		{&egoscale.ListVolumes{}, "list"},
		{&egoscale.ResizeVolume{}, "resize"},
	},
	"template": {
		{&egoscale.CopyTemplate{}, ""},
		{&egoscale.CreateTemplate{}, ""},
		{&egoscale.ListTemplates{}, "list"},
		{&egoscale.PrepareTemplate{}, ""},
		{&egoscale.RegisterTemplate{}, ""},
		{&egoscale.ListOSCategories{}, ""},
	},
	"account": {
		{&egoscale.EnableAccount{}, ""},
		{&egoscale.DisableAccount{}, ""},
		{&egoscale.ListAccounts{}, "list"},
	},
	"zone": {
		{&egoscale.ListZones{}, "list"},
	},
	"snapshot": {
		{&egoscale.CreateSnapshot{}, "create"},
		{&egoscale.DeleteSnapshot{}, "delete"},
		{&egoscale.ListSnapshots{}, "list"},
		{&egoscale.RevertSnapshot{}, "revert"},
	},
	"user": {
		{&egoscale.CreateUser{}, ""},
		{&egoscale.DeleteUser{}, ""},
		//{&egoscale.DisableUser{}, ""},
		//{&egoscale.GetUser{}, ""},
		{&egoscale.UpdateUser{}, ""},
		{&egoscale.ListUsers{}, "list"},
		{&egoscale.RegisterUserKeys{}, ""},
	},
	"security-group,sg": {
		{&egoscale.AuthorizeSecurityGroupEgress{}, "authorizeEgress"},
		{&egoscale.AuthorizeSecurityGroupIngress{}, "authorizeIngress"},
		{&egoscale.CreateSecurityGroup{}, "create"},
		{&egoscale.DeleteSecurityGroup{}, "delete"},
		{&egoscale.ListSecurityGroups{}, "list"},
		{&egoscale.RevokeSecurityGroupEgress{}, "revokeEgress"},
		{&egoscale.RevokeSecurityGroupIngress{}, "revokeIngress"},
	},
	"ssh": {
		{&egoscale.RegisterSSHKeyPair{}, "register"},
		{&egoscale.ListSSHKeyPairs{}, "list"},
		{&egoscale.CreateSSHKeyPair{}, "create"},
		{&egoscale.DeleteSSHKeyPair{}, "delete"},
		{&egoscale.ResetSSHKeyForVirtualMachine{}, "reset"},
	},
	"affinity-group,ag": {
		{&egoscale.CreateAffinityGroup{}, "create"},
		{&egoscale.DeleteAffinityGroup{}, "delete"},
		{&egoscale.ListAffinityGroups{}, "list"},
		{&egoscale.UpdateVMAffinityGroup{}, ""},
	},
	"vm-group,vg": {
		{&egoscale.CreateInstanceGroup{}, "create"},
		{&egoscale.ListInstanceGroups{}, "list"},
	},
	"tags": {
		{&egoscale.CreateTags{}, "create"},
		{&egoscale.DeleteTags{}, "delete"},
		{&egoscale.ListTags{}, "list"},
	},
	"nic": {
		{&egoscale.ActivateIP6{}, ""},
		{&egoscale.AddIPToNic{}, ""},
		{&egoscale.ListNics{}, "list"},
		{&egoscale.RemoveIPFromNic{}, ""},
	},
	"address": {
		{&egoscale.AssociateIPAddress{}, "associate"},
		{&egoscale.DisassociateIPAddress{}, "disassociate"},
		{&egoscale.ListPublicIPAddresses{}, "list"},
		{&egoscale.UpdateIPAddress{}, "update"},
	},
	"async-job,aj": {
		{&egoscale.QueryAsyncJobResult{}, ""},
		{&egoscale.ListAsyncJobs{}, ""},
	},
	"apis": {
		{&egoscale.ListAPIs{}, "list"},
	},
	"event": {
		{&egoscale.ListEventTypes{}, "listType"},
		{&egoscale.ListEvents{}, "list"},
	},
	"offerings": {
		{&egoscale.ListResourceDetails{}, "listDetails"},
		{&egoscale.ListResourceLimits{}, "listLimits"},
		{&egoscale.ListServiceOfferings{}, "list"},
	},
	"host": {
		{&egoscale.ListHosts{}, ""},
		{&egoscale.UpdateHost{}, ""},
	},
	"reversedns": {
		{&egoscale.DeleteReverseDNSFromPublicIPAddress{}, "deleteFromAddress"},
		{&egoscale.QueryReverseDNSForPublicIPAddress{}, "queryForAddress"},
		{&egoscale.UpdateReverseDNSForPublicIPAddress{}, "updateForAddress"},
		{&egoscale.DeleteReverseDNSFromVirtualMachine{}, "deleteFromVM"},
		{&egoscale.QueryReverseDNSForVirtualMachine{}, "queryForVM"},
		{&egoscale.UpdateReverseDNSForVirtualMachine{}, "updateForVM"},
	},
}
