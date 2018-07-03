package cmd

import (
	"fmt"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var privnetShowCmd = &cobra.Command{
	Use:   "show <privnet name | id>",
	Short: "Show a private network details",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		network, vms, err := privnetDetails(args[0])
		if err != nil {
			return err
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Zone", "Name", "Virtual Machine", "Virtual Machine ID"})
		zone := network.ZoneName
		name := network.Name
		if len(vms) > 0 {
			for _, vm := range vms {
				table.Append([]string{zone, name, vm.Name, vm.ID})
				zone = ""
				name = ""
			}
		} else {
			table.Append([]string{zone, network.Name, "", ""})
		}
		table.Render()

		return nil
	},
}

func privnetDetails(privnetName string) (*egoscale.Network, []egoscale.VirtualMachine, error) {

	resp, err := cs.List(&egoscale.Zone{})
	if err != nil {
		return nil, nil, err
	}

	network, err := searchPrivnet(privnetName, resp)
	if err != nil {
		return nil, nil, err
	}

	vms, err := cs.List(&egoscale.VirtualMachine{ZoneID: network.ZoneID})
	if err != nil {
		return nil, nil, err
	}

	var vmsRes []egoscale.VirtualMachine
	for _, v := range vms {
		vm := v.(*egoscale.VirtualMachine)

		if _, err := containNetID(network, vm.Nic); err == nil {
			vmsRes = append(vmsRes, *vm)
		}
	}

	return network, vmsRes, nil
}

func searchPrivnet(privName string, zones []interface{}) (*egoscale.Network, error) {

	for _, z := range zones {
		zone := z.(*egoscale.Zone)
		net, err := getNetworkIDByName(cs, privName, zone.ID)
		if err != nil {
			continue
		}
		return net, nil
	}
	return nil, fmt.Errorf("Private Network %q not found", privName)
}

func init() {
	privnetCmd.AddCommand(privnetShowCmd)
}
