package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

func init() {

	for i := egoscale.A; i <= egoscale.URL; i++ {

		var cmdUpdateRecord = &cobra.Command{
			Use:   fmt.Sprintf("%s <domain name> <record name | id>", egoscale.Record.String(i)),
			Short: fmt.Sprintf("Update %s record type to a domain", egoscale.Record.String(i)),
			Run: func(cmd *cobra.Command, args []string) {
				if len(args) < 2 {
					cmd.Usage()
					return
				}

				recordID, err := csDNS.GetRecordIDByName(args[0], args[1])
				if err != nil {
					log.Fatal(err)
				}

				name, err := cmd.Flags().GetString("name")
				if err != nil {
					log.Fatal(err)
				}
				addr, err := cmd.Flags().GetString("content")
				if err != nil {
					log.Fatal(err)
				}
				ttl, err := cmd.Flags().GetInt("ttl")
				if err != nil {
					log.Fatal(err)
				}

				domain, err := csDNS.GetDomain(args[0])
				if err != nil {
					log.Fatal(err)
				}

				resp, err := csDNS.UpdateRecord(args[0], egoscale.UpdateDNSRecord{
					ID:         recordID,
					DomainID:   domain.ID,
					TTL:        ttl,
					RecordType: egoscale.Record.String(i),
					Name:       name,
					Content:    addr,
				})
				if err != nil {
					log.Fatal(err)
				}
				println(resp.ID)
			},
		}
		cmdUpdateRecord.Flags().StringP("name", "n", "", "Update name")
		cmdUpdateRecord.Flags().StringP("content", "c", "", "Update Content")
		cmdUpdateRecord.Flags().IntP("ttl", "t", 0, "Update ttl")
		cmdUpdateRecord.Flags().IntP("priority", "p", 0, "Update priority")
		dnsUpdateCmd.AddCommand(cmdUpdateRecord)
	}
}
