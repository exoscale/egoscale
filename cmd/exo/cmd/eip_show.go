// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"os"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var eipShowCmd = &cobra.Command{
	Use:     "show <ip address | eip id>",
	Short:   "Show an eip details",
	Aliases: gShowAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}
		return eipDetails(args[0])
	},
}

func eipDetails(eip string) error {

	var eipID = eip
	if !isUUID(eip) {
		var err error
		eipID, err = getEIPIDByIP(cs, eip)
		if err != nil {
			return err
		}
	}

	addr := &egoscale.IPAddress{ID: eipID, IsElastic: true}
	if err := cs.Get(addr); err != nil {
		return err
	}

	vms, err := cs.List(&egoscale.VirtualMachine{ZoneID: addr.ZoneID})
	if err != nil {
		return err
	}

	vmAssociated := []egoscale.VirtualMachine{}

	for _, value := range vms {
		vm := value.(*egoscale.VirtualMachine)
		nic := vm.DefaultNic()
		if nic == nil {
			continue
		}
		for _, sIP := range nic.SecondaryIP {
			if sIP.IPAddress.Equal(addr.IPAddress) {
				vmAssociated = append(vmAssociated, *vm)
			}
		}
	}

	table := table.NewTable(os.Stdout)
	table.SetHeader([]string{"Zone", "IP", "Virtual Machine", "Virtual Machine ID"})

	zone := addr.ZoneName
	ipaddr := addr.IPAddress.String()
	if len(vmAssociated) > 0 {
		for _, vm := range vmAssociated {
			table.Append([]string{zone, ipaddr, vm.Name, vm.ID})
			zone = ""
			ipaddr = ""
		}
	} else {
		table.Append([]string{zone, ipaddr})
	}
	table.Render()

	return nil
}

func init() {
	eipCmd.AddCommand(eipShowCmd)
}
