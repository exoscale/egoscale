package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/exoscale/egoscale"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var dnsListCmd = &cobra.Command{
	Use:   "list [domain name]",
	Short: "List all domain or domain records",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			domains, err := listDomains()
			if err != nil {
				log.Fatal(err)
			}

			table := table.NewTable(os.Stdout)
			table.SetHeader([]string{"Name", "ID"})

			for _, domain := range domains {
				table.Append([]string{domain.Name, fmt.Sprintf("%d", domain.ID)})
			}
			table.Render()
			return
		}
		records, err := csDNS.GetRecords(args[0])
		if err != nil {
			log.Fatal(err)
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Type", "Name", "Content", "TTL", "Priority", "id"})

		for _, record := range records {
			s := []string{
				record.RecordType,
				record.Name,
				record.Content,
				fmt.Sprintf("%d", record.TTL),
				fmt.Sprintf("%d", record.Prio),
				fmt.Sprintf("%v", record.ID),
			}
			table.Append(s)
		}

		table.Render()
	},
}

func listDomains() ([]egoscale.DNSDomain, error) {
	domains, err := csDNS.GetDomains()
	if err != nil {
		return nil, err
	}
	return domains, nil
}

func init() {
	dnsCmd.AddCommand(dnsListCmd)
}
