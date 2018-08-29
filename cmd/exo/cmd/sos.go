package cmd

import (
	"strings"

	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// sosCmd represents the sos command
var sosCmd = &cobra.Command{
	Use:   "sos",
	Short: "Simple Object Storage management",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		zone, err := cmd.Flags().GetString("zone")
		if err != nil {
			return err
		}

		if zone != "" {
			gCurrentAccount.DefaultZone = zone
		}
		return nil
	},
}

func newMinioClient(zone string) (*minio.Client, error) {
	endpoint := strings.Replace(gCurrentAccount.SosEndpoint, "https://", "", -1)
	endpoint = strings.Replace(endpoint, "{zone}", zone, -1)
	return minio.NewV4(endpoint, gCurrentAccount.Key, gCurrentAccount.Secret, true)
}

func init() {
	RootCmd.AddCommand(sosCmd)
	sosCmd.PersistentFlags().StringP("zone", "z", "", "Simple object storage zone")
}
