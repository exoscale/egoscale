package cmd

import (
	"log"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// dnsCreateCmd represents the create command
var dnsCreateCmd = &cobra.Command{
	Use:     "create <domain name>",
	Short:   "Create a domain",
	Aliases: gCreateAlias,
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

	return csDNS.CreateDomain(domainName)
}

func init() {
	dnsCmd.AddCommand(dnsCreateCmd)
}
