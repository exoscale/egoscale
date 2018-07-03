package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// privnetCmd represents the pn command
var privnetCmd = &cobra.Command{
	Use:   "privnet",
	Short: "Private networks management",
}

func getNetworkIDByName(cs *egoscale.Client, name, zone string) (*egoscale.Network, error) {
	nets, err := cs.List(&egoscale.Network{Type: "Isolated", CanUseForDeploy: true, ZoneID: zone})
	if err != nil {
		log.Fatal(err)
	}

	var res *egoscale.Network
	match := 0
	for _, net := range nets {
		n := net.(*egoscale.Network)
		if strings.Compare(name, n.Name) == 0 || strings.Compare(name, n.ID) == 0 {
			res = n
			match++
		}
	}
	switch match {
	case 0:
		return nil, fmt.Errorf("Unable to find this private network")
	case 1:
		return res, nil
	default:
		return nil, fmt.Errorf("Multiple private network found")

	}
}

func init() {
	RootCmd.AddCommand(privnetCmd)
}
