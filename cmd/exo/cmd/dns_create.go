package cmd

import (
	"log"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// dnsCreateCmd represents the create command
var dnsCreateCmd = &cobra.Command{
	Use:   "create <domain name>",
	Short: "Create a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		resp, err := createDomain(args[0])
		if err != nil {
			log.Fatal(err)
		}

		println(resp.ID)

	},
}

func createDomain(domainName string) (*egoscale.DNSDomain, error) {

	resp, err := csDNS.CreateDomain(domainName)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func init() {
	dnsCmd.AddCommand(dnsCreateCmd)
}
