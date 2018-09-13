package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var privnetUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update private network",
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		desc, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		startip, err := cmd.Flags().GetString("startip")
		if err != nil {
			return err
		}
		endip, err := cmd.Flags().GetString("endip")
		if err != nil {
			return err
		}
		netmask, err := cmd.Flags().GetString("netmask")
		if err != nil {
			return err
		}
		return privnetUpdate(id, name, desc, startip, endip, netmask)
	},
}

func privnetUpdate(id, name, desc, startIPAddr, endIPAddr, netmaskAddr string) error {
	uuid, err := egoscale.ParseUUID(id)
	if err != nil {
		return fmt.Errorf("update the network by ID, got %q", id)
	}
	var startip, endip, netmask net.IP
	if startIPAddr != "" {
		startip = net.ParseIP(startIPAddr)
	}
	if endIPAddr != "" {
		endip = net.ParseIP(endIPAddr)
	}
	if netmaskAddr != "" {
		netmask = net.ParseIP(netmaskAddr)
	}

	req := &egoscale.UpdateNetwork{
		ID:          uuid,
		DisplayText: desc,
		Name:        name,
		StartIP:     startip,
		EndIP:       endip,
		Netmask:     netmask,
	}

	resp, err := asyncRequest(req, fmt.Sprintf("Updating the network %q ", id))
	if err != nil {
		return err
	}

	newNet := resp.(*egoscale.Network)

	table := table.NewTable(os.Stdout)
	table.SetHeader([]string{"Name", "Description", "ID", "StartIP", "EndIP", "Netmask"})
	table.Append([]string{newNet.Name, newNet.DisplayText, newNet.ID.String(), newNet.StartIP.String(), newNet.EndIP.String(), newNet.Netmask.String()})
	table.Render()
	return nil
}

func init() {
	privnetUpdateCmd.Flags().StringP("id", "i", "", "Private network id")
	privnetUpdateCmd.Flags().StringP("name", "n", "", "Private network name")
	privnetUpdateCmd.Flags().StringP("description", "d", "", "Private network description")
	privnetUpdateCmd.Flags().StringP("startip", "s", "", "the beginning IP address in the network IP range. Required for managed networks.")
	privnetUpdateCmd.Flags().StringP("endip", "e", "", "the ending IP address in the network IP range. Required for managed networks.")
	privnetUpdateCmd.Flags().StringP("netmask", "m", "", "the netmask of the network.  Required for managed networks")
	privnetCmd.AddCommand(privnetUpdateCmd)
}
