package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// dnsDeleteCmd represents the delete command
var dnsDeleteCmd = &cobra.Command{
	Use:   "delete <domain name>",
	Short: "Delete a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		if err := deleteDomain(args[0]); err != nil {
			log.Fatal(err)
		}

		println(args[0])
	},
}

func deleteDomain(domainName string) error {
	if err := csDNS.DeleteDomain(domainName); err != nil {
		return err
	}
	return nil
}

func init() {
	dnsCmd.AddCommand(dnsDeleteCmd)
}
