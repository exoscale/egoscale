package cmd

import (
	"errors"
	"os"
	"strconv"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

func init() {
	dnsCmd.AddCommand(dnsShowCmd)
	dnsShowCmd.Flags().StringP("name", "n", "", "List records by name")
	dnsShowCmd.Flags().StringP("content", "c", "", "List records by content keyword")
}

// dnsShowCmd represents the show command
var dnsShowCmd = &cobra.Command{
	Use:   "show <domain name | id> [type]",
	Short: "Show the domain records",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("show expects one DNS domain by name or id")
		}

		types := []string{}
		if len(args) > 1 {
			types = args[1:]
		} else {
			types = []string{""}
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}
		content, err := cmd.Flags().GetString("content")
		if err != nil {
			return err
		}

		t := table.NewTable(os.Stdout)
		err = domainListRecords(t, args[0], name, content, types)
		if err == nil {
			t.Render()
		}
		return err
	},
}

func domainListRecords(t *table.Table, domain, name, content string, types []string) error {

	t.SetHeader([]string{"Type", "Name", "Content", "TTL", "Prio", "ID"})

	for _, recordType := range types {
		records, err := csDNS.GetRecordsWithFilters(domain, name, content, recordType)
		if err != nil {
			return err
		}

		for _, record := range records {
			t.Append([]string{
				record.RecordType,
				record.Name,
				record.Content,
				strconv.Itoa(record.TTL),
				strconv.Itoa(record.Prio),
				strconv.FormatInt(record.ID, 10),
			})
		}
	}

	return nil
}
