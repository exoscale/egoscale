package cmd

import (
	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

const userDocumentationURL string = "http://cloudstack.apache.org/api/apidocs-4.4/user/%s.html"
const rootDocumentationURL string = "http://cloudstack.apache.org/api/apidocs-4.4/root_admin/%s.html"

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api <command>",
	Short: "Exoscale api",
}

func init() {
	RootCmd.AddCommand(apiCmd)
	var method egoscale.Command
	buildCommands(&method, methods)
	apiCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode on")
	apiCmd.PersistentFlags().BoolP("dry-run", "D", false, "produce a cURL ready URL")
	apiCmd.PersistentFlags().BoolP("dry-json", "j", false, "produce a JSON preview of the query")
	apiCmd.PersistentFlags().StringP("theme", "t", "", "syntax highlighting theme, see: https://xyproto.github.io/splash/docs/")

}

func buildCommands(out *egoscale.Command, methods map[string][]cmd) []cobra.Command {
	commands := make([]cobra.Command, 0)

	for category, ms := range methods {

		cmd := cobra.Command{
			Use: category,
		}

		apiCmd.AddCommand(&cmd)

		for i := range ms {
			s := ms[i]

			name := cs.APIName(s.command)
			description := cs.APIDescription(s.command)

			subCMD := cobra.Command{
				Use:   name,
				Short: description,
				RunE: func(cmd *cobra.Command, args []string) error {
					return cmd.Usage()
				},
			}
			// report back the current command
			cmd.AddCommand(&subCMD)
			commands = append(commands, cmd)
		}
	}

	return commands
}

type cmd struct {
	command egoscale.Command
	hidden  bool
}

var methods = map[string][]cmd{
	"network": {
		{&egoscale.CreateNetwork{}, false},
		{&egoscale.DeleteNetwork{}, false},
		{&egoscale.ListNetworkOfferings{}, false},
		{&egoscale.ListNetworks{}, false},
		{&egoscale.RestartNetwork{}, true},
		{&egoscale.UpdateNetwork{}, false},
	},
	"virtual machine": {
		{&egoscale.AddNicToVirtualMachine{}, false},
		{&egoscale.ChangeServiceForVirtualMachine{}, false},
		{&egoscale.DeployVirtualMachine{}, false},
		{&egoscale.DestroyVirtualMachine{}, false},
		{&egoscale.ExpungeVirtualMachine{}, false},
		{&egoscale.GetVMPassword{}, false},
		{&egoscale.GetVirtualMachineUserData{}, false},
		{&egoscale.ListVirtualMachines{}, false},
		{&egoscale.MigrateVirtualMachine{}, true},
		{&egoscale.RebootVirtualMachine{}, false},
		{&egoscale.RecoverVirtualMachine{}, false},
		{&egoscale.RemoveNicFromVirtualMachine{}, false},
		{&egoscale.ResetPasswordForVirtualMachine{}, false},
		{&egoscale.RestoreVirtualMachine{}, false},
		{&egoscale.ScaleVirtualMachine{}, false},
		{&egoscale.StartVirtualMachine{}, false},
		{&egoscale.StopVirtualMachine{}, false},
		{&egoscale.UpdateDefaultNicForVirtualMachine{}, true},
		{&egoscale.UpdateVirtualMachine{}, false},
	},
	"volume": {
		{&egoscale.ListVolumes{}, false},
		{&egoscale.ResizeVolume{}, false},
	},
	"template": {
		{&egoscale.CopyTemplate{}, true},
		{&egoscale.CreateTemplate{}, true},
		{&egoscale.ListTemplates{}, false},
		{&egoscale.PrepareTemplate{}, true},
		{&egoscale.RegisterTemplate{}, true},
		{&egoscale.ListOSCategories{}, true},
	},
	"account": {
		{&egoscale.EnableAccount{}, true},
		{&egoscale.DisableAccount{}, true},
		{&egoscale.ListAccounts{}, false},
	},
	"zone": {
		{&egoscale.ListZones{}, false},
	},
	"snapshot": {
		{&egoscale.CreateSnapshot{}, false},
		{&egoscale.DeleteSnapshot{}, false},
		{&egoscale.ListSnapshots{}, false},
		{&egoscale.RevertSnapshot{}, false},
	},
	"user": {
		{&egoscale.CreateUser{}, true},
		{&egoscale.DeleteUser{}, true},
		//{&egoscale.DisableUser{}, true},
		//{&egoscale.GetUser{}, true},
		{&egoscale.UpdateUser{}, true},
		{&egoscale.ListUsers{}, false},
		{&egoscale.RegisterUserKeys{}, false},
	},
	"security group": {
		{&egoscale.AuthorizeSecurityGroupEgress{}, false},
		{&egoscale.AuthorizeSecurityGroupIngress{}, false},
		{&egoscale.CreateSecurityGroup{}, false},
		{&egoscale.DeleteSecurityGroup{}, false},
		{&egoscale.ListSecurityGroups{}, false},
		{&egoscale.RevokeSecurityGroupEgress{}, false},
		{&egoscale.RevokeSecurityGroupIngress{}, false},
	},
	"ssh": {
		{&egoscale.RegisterSSHKeyPair{}, false},
		{&egoscale.ListSSHKeyPairs{}, false},
		{&egoscale.CreateSSHKeyPair{}, false},
		{&egoscale.DeleteSSHKeyPair{}, false},
		{&egoscale.ResetSSHKeyForVirtualMachine{}, false},
	},
	"affinity group": {
		{&egoscale.CreateAffinityGroup{}, false},
		{&egoscale.DeleteAffinityGroup{}, false},
		{&egoscale.ListAffinityGroups{}, false},
		{&egoscale.UpdateVMAffinityGroup{}, false},
	},
	"vm group": {
		{&egoscale.CreateInstanceGroup{}, false},
		{&egoscale.ListInstanceGroups{}, false},
	},
	"tags": {
		{&egoscale.CreateTags{}, false},
		{&egoscale.DeleteTags{}, false},
		{&egoscale.ListTags{}, false},
	},
	"nic": {
		{&egoscale.ActivateIP6{}, false},
		{&egoscale.AddIPToNic{}, false},
		{&egoscale.ListNics{}, false},
		{&egoscale.RemoveIPFromNic{}, false},
	},
	"address": {
		{&egoscale.AssociateIPAddress{}, false},
		{&egoscale.DisassociateIPAddress{}, false},
		{&egoscale.ListPublicIPAddresses{}, false},
		{&egoscale.UpdateIPAddress{}, false},
	},
	"async job": {
		{&egoscale.QueryAsyncJobResult{}, false},
		{&egoscale.ListAsyncJobs{}, false},
	},
	"apis": {
		{&egoscale.ListAPIs{}, false},
	},
	"event": {
		{&egoscale.ListEventTypes{}, false},
		{&egoscale.ListEvents{}, false},
	},
	"offerings": {
		{&egoscale.ListResourceDetails{}, false},
		{&egoscale.ListResourceLimits{}, false},
		{&egoscale.ListServiceOfferings{}, false},
	},
	"host": {
		{&egoscale.ListHosts{}, true},
		{&egoscale.UpdateHost{}, true},
	},
	"reversedns": {
		{&egoscale.DeleteReverseDNSFromPublicIPAddress{}, false},
		{&egoscale.QueryReverseDNSForPublicIPAddress{}, false},
		{&egoscale.UpdateReverseDNSForPublicIPAddress{}, false},
		{&egoscale.DeleteReverseDNSFromVirtualMachine{}, false},
		{&egoscale.QueryReverseDNSForVirtualMachine{}, false},
		{&egoscale.UpdateReverseDNSForVirtualMachine{}, false},
	},
}
