package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// dnsRemoveCmd represents the remove command
var dnsRemoveCmd = &cobra.Command{
	Use:   "remove <domain name> <record name | id>",
	Short: "Remove a domain record",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}
		id, err := removeRecord(args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}
		println(id)
	},
}

func removeRecord(domainName, record string) (int64, error) {
	id, err := getRecordIDByName(csDNS, domainName, record)
	if err != nil {
		return 0, err
	}
	if err := csDNS.DeleteRecord(domainName, id); err != nil {
		return 0, err
	}

	return id, nil
}

func init() {
	dnsCmd.AddCommand(dnsRemoveCmd)
}
