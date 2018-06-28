package cmd

import (
	"log"

	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var sosCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "create bucket",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		zone, err := cmd.Flags().GetString("zone")
		if err != nil {
			log.Fatal(err)
		}

		minioClient, err := newMinioClient(zone)
		if err != nil {
			log.Fatal(err)
		}

		println("test")

		if err := createBucket(minioClient, args[0], zone); err != nil {
			log.Fatal(err)
		}
	},
}

func createBucket(minioClient *minio.Client, bucketName, zone string) error {
	return minioClient.MakeBucket(bucketName, zone)
}

func init() {
	sosCmd.AddCommand(sosCreateCmd)
	sosCreateCmd.Flags().StringP("zone", "z", "ch-dk-2", "Object storage zone")
}
