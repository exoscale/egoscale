package cmd

import (
	"fmt"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// dissociateCmd represents the dissociate command
var dissociateCmd = &cobra.Command{
	Use:   "dissociate <privnet name | id> <vm name | vm id> [vm name | vm id] [...]",
	Short: "Dissociate a private network from instance(s)",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		for _, vm := range args[1:] {
			resp, err := dissociatePrivNet(args[0], vm)
			if err != nil {
				return err
			}
			println(resp)
		}
		return nil
	},
}

func dissociatePrivNet(privnetName, vmName string) (string, error) {
	vm, err := getVMWithNameOrID(cs, vmName)
	if err != nil {
		return "", err
	}

	network, err := getNetworkIDByName(cs, privnetName, vm.ZoneID)
	if err != nil {
		return "", err
	}

	nic, err := containNetID(network, vm.Nic)
	if err != nil {
		return "", err
	}

	_, err = cs.Request(&egoscale.RemoveNicFromVirtualMachine{NicID: nic.ID, VirtualMachineID: vm.ID})
	if err != nil {
		return "", err
	}

	return nic.ID, nil
}

func containNetID(network *egoscale.Network, vmNics []egoscale.Nic) (egoscale.Nic, error) {

	for _, nic := range vmNics {
		if nic.NetworkID == network.ID {
			return nic, nil
		}
	}
	return egoscale.Nic{}, fmt.Errorf("NIC not found")
}

func init() {
	privnetCmd.AddCommand(dissociateCmd)
}
