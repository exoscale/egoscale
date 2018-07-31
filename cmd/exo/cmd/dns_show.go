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
}

// dnsShowCmd represents the show command
var dnsShowCmd = &cobra.Command{
	Use:   "show <domain name | id>",
	Short: "Show the domain records",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("show expects one DNS domain by name or id")
		}

		t := table.NewTable(os.Stdout)
		err := domainListRecords(t, args[0])
		if err == nil {
			t.Render()
		}
		return err
	},
}

func domainListRecords(t *table.Table, name string) error {
	records, err := csDNS.GetRecords(name)
	if err != nil {
		return err
	}

	t.SetHeader([]string{"Type", "Name", "Content", "TTL", "Prio", "ID"})

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

	return nil
}
