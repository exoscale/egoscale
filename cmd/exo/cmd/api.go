package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"reflect"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

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

func buildCommands(out *egoscale.Command, methods map[string][]cmd) {
	for category, ms := range methods {
		cmd := cobra.Command{
			Use: strings.Replace(category, " ", "-", -1),
		}

		apiCmd.AddCommand(&cmd)

		for i := range ms {
			s := ms[i]

			name := cs.APIName(s.command)
			description := cs.APIDescription(s.command)

			subCMD := cobra.Command{
				Use:  name,
				Long: description,
			}
			buildFlags(s.command, &subCMD)

			subCMD.RunE = func(cmd *cobra.Command, args []string) error {
				resp, err := cs.RequestWithContext(gContext, s.command)
				if err != nil {
					return err
				}

				data, err := json.MarshalIndent(&resp, "", "  ")
				if err != nil {
					return err
				}

				fmt.Println(string(data))

				return nil
			}

			cmd.AddCommand(&subCMD)
		}
	}
}

func buildFlags(method egoscale.Command, cmd *cobra.Command) {
	val := reflect.ValueOf(method)
	// we've got a pointer
	value := val.Elem()

	if value.Kind() != reflect.Struct {
		log.Fatalf("struct was expected")
		return
	}

	ty := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := ty.Field(i)

		if field.Name == "_" {
			continue
		}

		// XXX refactor with request.go
		var argName string
		required := false
		if json, ok := field.Tag.Lookup("json"); ok {
			tags := strings.Split(json, ",")
			argName = tags[0]
			required = true
			for _, tag := range tags {
				if tag == "omitempty" {
					required = false
				}
			}
			if argName == "" || argName == "omitempty" {
				continue
			}
		}

		description := ""
		if required {
			description = "required"
		}

		if doc, ok := field.Tag.Lookup("doc"); ok {
			if description != "" {
				description = fmt.Sprintf("[%s] %s", description, doc)
			} else {
				description = doc
			}
		}

		val := value.Field(i)
		addr := val.Addr().Interface()
		switch val.Kind() {
		case reflect.Bool:
			cmd.Flags().BoolVarP(addr.(*bool), argName, "", false, description)
		case reflect.Int:
			cmd.Flags().IntVarP(addr.(*int), argName, "", 0, description)
		case reflect.Int64:
			cmd.Flags().Int64VarP(addr.(*int64), argName, "", 0, description)
		case reflect.Uint:
			cmd.Flags().UintVarP(addr.(*uint), argName, "", 0, description)
		case reflect.Uint64:
			cmd.Flags().Uint64VarP(addr.(*uint64), argName, "", 0, description)
		case reflect.Float64:
			cmd.Flags().Float64VarP(addr.(*float64), argName, "", 0, description)
		case reflect.Int16:
			typeName := field.Type.Name()
			if typeName != "int16" {
				cmd.Flags().VarP(&intTypeGeneric{addr: addr, base: 10, bitSize: 16, typ: field.Type}, argName, "", description)
			} else {
				cmd.Flags().Int16VarP(addr.(*int16), argName, "", 0, description)
			}
		case reflect.Uint8:
			cmd.Flags().Uint8VarP(addr.(*uint8), argName, "", 0, description)
		case reflect.Uint16:
			cmd.Flags().Uint16VarP(addr.(*uint16), argName, "", 0, description)
		case reflect.String:
			typeName := field.Type.Name()
			if typeName != "string" {
				cmd.Flags().VarP(&stringerTypeGeneric{addr: addr, typ: field.Type}, argName, "", description)
			} else {
				cmd.Flags().StringVarP(addr.(*string), argName, "", "", description)
			}
		case reflect.Slice:
			switch field.Type.Elem().Kind() {
			case reflect.Uint8:
				ip := addr.(*net.IP)
				if *ip == nil || (*ip).Equal(net.IPv4zero) || (*ip).Equal(net.IPv6zero) {
					cmd.Flags().IPP(argName, "", *ip, description)
				}
			case reflect.String:
				cmd.Flags().StringSliceP(argName, "", *addr.(*[]string), description)
			default:
				switch field.Type.Elem() {
				case reflect.TypeOf(egoscale.ResourceTag{}):
					cmd.Flags().VarP(&tagGeneric{addr.(*[]egoscale.ResourceTag)}, argName, "", description)
				case reflect.TypeOf(egoscale.CIDR{}):
					cmd.Flags().VarP(&cidrListGeneric{addr.(*[]egoscale.CIDR)}, argName, "", description)
				case reflect.TypeOf(egoscale.UUID{}):
					cmd.Flags().VarP(&uuidListGeneric{addr.(*[]egoscale.UUID)}, argName, "", description)
				default:
					//log.Printf("[SKIP] Slice of %s is not supported!", field.Name)
				}
			}
		case reflect.Map:
			key := reflect.TypeOf(val.Interface()).Key()
			switch key.Kind() {
			case reflect.String:
				cmd.Flags().VarP(&mapGeneric{addr.(*map[string]string)}, argName, "", description)
			default:
				log.Printf("[SKIP] Type map for %s is not supported!", field.Name)
			}
		case reflect.Ptr:
			switch field.Type.Elem() {
			case reflect.TypeOf(true):
				cmd.Flags().VarP(&boolFlag{(addr.(**bool))}, argName, "", description)
			case reflect.TypeOf(egoscale.CIDR{}):
				cmd.Flags().VarP(&cidr{addr.(**egoscale.CIDR)}, argName, "", description)
			case reflect.TypeOf(egoscale.UUID{}):
				cmd.Flags().VarP(&uuid{addr.(**egoscale.UUID)}, argName, "", description)

			default:
				log.Printf("[SKIP] Ptr type of %s is not supported!", field.Name)
			}
		default:
			log.Printf("[SKIP] Type of %s is not supported! %v", field.Name, val.Kind())
		}
	}
}

type cmd struct {
	command egoscale.Command
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
	"virtual-machine": {
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
	"security-group": {
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
	"affinity-group": {
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
