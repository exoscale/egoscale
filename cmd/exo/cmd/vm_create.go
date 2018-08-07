package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/utils"
	"github.com/spf13/cobra"
)

// vmCreateCmd represents the create command
var vmCreateCmd = &cobra.Command{
	Use:     "create <vm name>",
	Short:   "Create and deploy a virtual machine",
	Aliases: gCreateAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		userDataPath, err := cmd.Flags().GetString("cloud-init-file")
		if err != nil {
			return err
		}

		userData := ""

		if userDataPath != "" {
			userData, err = getUserData(userDataPath)
			if err != nil {
				return err
			}
		}

		zone, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}

		if zone == "" {
			zone = gCurrentAccount.DefaultZone
		}

		zone, err = getZoneIDByName(zone)
		if err != nil {
			return err
		}

		templateName, err := cmd.Flags().GetString("template")
		if err != nil {
			return err
		}

		diskSize, err := cmd.Flags().GetInt64("disk")
		if err != nil {
			return err
		}

		template, err := getTemplateByName(zone, templateName)
		if err != nil {
			return err
		}

		keypair, err := cmd.Flags().GetString("keypair")
		if err != nil {
			return err
		}

		sg, err := cmd.Flags().GetString("security-group")
		if err != nil {
			return err
		}

		sgs, err := getSecurityGroups(cs, sg)
		if err != nil {
			return err
		}

		ipv6, err := cmd.Flags().GetBool("ipv6")
		if err != nil {
			return err
		}

		privnet, err := cmd.Flags().GetString("privnet")
		if err != nil {
			return err
		}

		pvs, err := getPrivnetList(cs, privnet, zone)
		if err != nil {
			return err
		}

		so, err := cmd.Flags().GetString("service-offering")
		if err != nil {
			return err
		}

		servOffering, err := getServiceOfferingByName(cs, so)
		if err != nil {
			return err
		}

		affinitygroup, err := cmd.Flags().GetString("anti-affinity-group")
		if err != nil {
			return err
		}

		affinitygroups, err := getAffinityGroup(cs, affinitygroup)
		if err != nil {
			return err
		}

		vmInfo := &egoscale.DeployVirtualMachine{
			Name:              args[0],
			UserData:          userData,
			ZoneID:            zone,
			TemplateID:        template.ID,
			RootDiskSize:      diskSize,
			KeyPair:           keypair,
			SecurityGroupIDs:  sgs,
			IP6:               &ipv6,
			NetworkIDs:        pvs,
			ServiceOfferingID: servOffering.ID,
			AffinityGroupIDs:  affinitygroups,
		}

		r, err := createVM(vmInfo)
		if err != nil {
			return err
		}

		sshinfo, err := getSSHInfo(r.ID, ipv6)
		if err != nil {
			return err
		}

		fmt.Printf(`
What to do now?

1. Connect to the machine

> exo ssh %s
`, r.Name)

		printSSHConnectSTR(sshinfo)

		fmt.Printf(`
2. Put the SSH configuration into ".ssh/config"

> exo ssh %s --info
`, r.Name)

		printSSHInfo(sshinfo)

		fmt.Print(`
Tip of the day:
	You're the sole owner of the private key.
	Be cautious with it.
`)

		return nil
	},
}

func getCommaflag(p string) []string {
	if p == "" {
		return nil
	}

	p = strings.Trim(p, ",")
	args := strings.Split(p, ",")

	res := []string{}

	for _, arg := range args {
		if arg == "" {
			continue
		}
		res = append(res, strings.TrimSpace(arg))
	}

	return res
}

func getSecurityGroups(cs *egoscale.Client, commaParameter string) ([]string, error) {

	sgs := getCommaflag(commaParameter)

	for i, sg := range sgs {
		s, err := getSecurityGroupByNameOrID(cs, sg)
		if err != nil {
			return nil, err
		}
		sgs[i] = s.ID
	}

	return sgs, nil
}

func getPrivnetList(cs *egoscale.Client, commaParameter, zoneID string) ([]string, error) {

	sgs := getCommaflag(commaParameter)

	for i, sg := range sgs {
		s, err := getNetworkIDByName(cs, sg)
		if err != nil {
			return nil, err
		}
		sgs[i] = s.ID
	}

	return sgs, nil
}

func getAffinityGroup(cs *egoscale.Client, commaParameter string) ([]string, error) {
	affs := getCommaflag(commaParameter)

	for i, aff := range affs {
		s, err := getAffinityGroupIDByName(cs, aff)
		if err != nil {
			return nil, err
		}
		affs[i] = s
	}

	return affs, nil
}

func getUserData(userDataPath string) (string, error) {
	buff, err := ioutil.ReadFile(userDataPath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buff), nil
}

func createVM(vmInfos *egoscale.DeployVirtualMachine) (*egoscale.VirtualMachine, error) {
	isDefaultKeyPair := false
	var keyPairs *egoscale.SSHKeyPair

	if vmInfos.KeyPair == "" {
		isDefaultKeyPair = true
		fmt.Println("Creating private SSH key")
		sshKeyName, err := utils.RandStringBytes(64)
		if err != nil {
			return nil, err
		}
		keyPairs, err = createSSHKey(sshKeyName)
		if err != nil {
			r := err.(*egoscale.ErrorResponse)
			if r.ErrorCode != egoscale.ParamError && r.CSErrorCode != egoscale.InvalidParameterValueException {
				return nil, err
			}
			return nil, fmt.Errorf("an SSH key with that name %q already exists, please choose a different name", sshKeyName)
		}
		defer deleteSSHKey(keyPairs.Name) // nolint: errcheck

		vmInfos.KeyPair = keyPairs.Name

	}

	virtualMachine := &egoscale.VirtualMachine{}

	if err := asyncRequest(vmInfos, fmt.Sprintf("Deploying %q ", vmInfos.Name), virtualMachine); err != nil {
		return nil, err
	}

	if isDefaultKeyPair {
		saveKeyPair(keyPairs, virtualMachine.ID)
	}

	return virtualMachine, nil
}

func saveKeyPair(keyPairs *egoscale.SSHKeyPair, vmID string) {
	filePath := path.Join(gConfigFolder, "instances", vmID)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			log.Fatal(err)
		}
	}

	filePath = path.Join(filePath, "id_rsa")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filePath, []byte(keyPairs.PrivateKey), 0600); err != nil {
			log.Fatalf("SSH private key could not be written. %s", err)
		}
	}
}

func init() {
	vmCreateCmd.Flags().StringP("cloud-init-file", "f", "", "Deploy instance with a cloud-init file")
	vmCreateCmd.Flags().StringP("zone", "z", "", "<zone name | id | keyword> (ch-dk-2|ch-gva-2|at-vie-1|de-fra-1)")
	vmCreateCmd.Flags().StringP("template", "t", "Linux Ubuntu 18.04", "<template name | id>")
	vmCreateCmd.Flags().Int64P("disk", "d", 50, "<disk size>")
	vmCreateCmd.Flags().StringP("keypair", "k", "", "<ssh keys name>")
	vmCreateCmd.Flags().StringP("security-group", "s", "", "<name | id, name | id, ...>")
	vmCreateCmd.Flags().StringP("privnet", "p", "", "<name | id, name | id, ...>")
	vmCreateCmd.Flags().StringP("anti-affinity-group", "A", "", "<name | id, name | id, ...>")
	vmCreateCmd.Flags().BoolP("ipv6", "6", false, "enable ipv6")
	vmCreateCmd.Flags().StringP("service-offering", "o", "Small", "<name | id> (micro|tiny|small|medium|large|extra-large|huge|mega|titan")
	vmCmd.AddCommand(vmCreateCmd)
}
