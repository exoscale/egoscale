package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/exoscale/egoscale/cmd/exo/table"

	"github.com/spf13/cobra"
)

// zoneCmd represents the zone command
var zoneCmd = &cobra.Command{
	Use:   "zone",
	Short: "List all available zones",
	RunE: func(cmd *cobra.Command, args []string) error {
		return listZones()
	},
}

func listZones() error {
	zones, err := cs.ListWithContext(gContext, &egoscale.Zone{})
	if err != nil {
		return err
	}

	table := table.NewTable(os.Stdout)
	table.SetHeader([]string{"Name", "ID"})

	for _, zone := range zones {
		z := zone.(*egoscale.Zone)
		table.Append([]string{z.Name, z.ID.String()})
	}
	table.Render()
	return nil
}

// func getZoneIDByName(name string) (*egoscale.UUID, error) {
func getZoneIDByName(name string) (string, error) {
	zoneReq := egoscale.Zone{}

	zones, err := cs.ListWithContext(gContext, &zoneReq)
	if err != nil {
		return "", err
	}

	var zoneID *egoscale.UUID

	for _, zone := range zones {
		z := zone.(*egoscale.Zone)
		if name == z.ID.String() {
			zoneID = z.ID
			break
		}

		if strings.Contains(strings.ToLower(z.Name), strings.ToLower(name)) {
			if zoneID != nil {
				return "", fmt.Errorf("more than one zones were found")
			}
			zoneID = z.ID
		}
	}

	if zoneID == nil {
		return "", fmt.Errorf("zone %q was not found", name)
	}

	return zoneID.String(), nil
}

func init() {
	RootCmd.AddCommand(zoneCmd)
}
