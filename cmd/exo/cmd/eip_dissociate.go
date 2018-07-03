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
	"fmt"
	"net"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// dissociateCmd represents the dissociate command
var dissociateCmd = &cobra.Command{
	Use:   "dissociate <IP address> <instance name | instance id>",
	Short: "Dissociate an IP from an instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}
		return dissociateIP(args[0], args[1])
	},
}

func dissociateIP(ipAddr, instance string) error {
	ip := net.ParseIP(ipAddr)
	if ip == nil {
		return fmt.Errorf("Invalide IP address")
	}

	vm, err := getVMWithNameOrID(cs, instance)
	if err != nil {
		return err
	}

	defaultNic := vm.DefaultNic()

	if defaultNic == nil {
		return fmt.Errorf("No default NIC on this instance")
	}

	eipID, err := getSecondaryIP(defaultNic, ip)
	if err != nil {
		return err
	}

	req := &egoscale.RemoveIPFromNic{ID: eipID}

	if err := cs.BooleanRequest(req); err != nil {
		return err
	}
	println(req.ID)
	return nil
}

func getSecondaryIP(nic *egoscale.Nic, ip net.IP) (string, error) {
	for _, sIP := range nic.SecondaryIP {
		if sIP.IPAddress.Equal(ip) {
			return sIP.ID, nil
		}
	}
	return "", fmt.Errorf("eip: %q not found", ip.String())
}

func init() {
	eipCmd.AddCommand(dissociateCmd)
}
