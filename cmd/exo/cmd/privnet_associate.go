package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// privnetAssociateCmd represents the associate command
var privnetAssociateCmd = &cobra.Command{
	Use:     "associate <privnet name | id> <vm name | vm id> [<ip>] [<vm name | vm id> [<ip>]] [...]",
	Short:   "Associate a private network to instance(s)",
	Aliases: gAssociateAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		network, err := getNetworkByName(args[0])
		if err != nil {
			return err
		}
		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Nic ID", "Virtual Machine ID", "IP Address"})
		for i := 1; i < len(args); i++ {
			name := args[i]
			if i != len(args)-1 {
				ip := net.ParseIP(args[i+1])
				if ip != nil {
					// the next param is an ip
					nic, err := associatePrivNet(network, name, ip)
					if err != nil {
						return err
					}
					table.Append([]string{
						nic.ID.String(),
						nic.VirtualMachineID.String(),
						nicIP(nic)})
					i = i + 1
					continue
				}
			}
			nic, err := associatePrivNet(network, name, nil)
			if err != nil {
				return err
			}
			table.Append([]string{
				nic.ID.String(),
				nic.VirtualMachineID.String(),
				nicIP(nic)})
		}
		table.Render()
		return nil
	},
}

func associatePrivNet(privnet *egoscale.Network, vmName string, ip net.IP) (*egoscale.Nic, error) {
	vm, err := getVMWithNameOrID(vmName)
	if err != nil {
		return nil, err
	}

	req := &egoscale.AddNicToVirtualMachine{NetworkID: privnet.ID, VirtualMachineID: vm.ID, IPAddress: ip}
	resp, err := cs.RequestWithContext(gContext, req)
	if err != nil {
		return nil, err
	}

	nic := resp.(*egoscale.VirtualMachine).NicByNetworkID(*privnet.ID)
	if nic == nil {
		return nil, fmt.Errorf("no nics found for network %q", privnet.ID)
	}

	return nic, nil

}

func init() {
	privnetCmd.AddCommand(privnetAssociateCmd)
}
