package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

// sosDeleteCmd represents the delete command
var sosDeleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Short:   "Delete a bucket",
	Aliases: gDeleteAlias,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		///XXX you must use a default zone support SOS
		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			println("XXX wainting for all zone supporting SOS: use a supported defaultZone")
			log.Fatal(err)
		}

		zone, err := minioClient.GetBucketLocation(args[0])
		if err != nil {
			log.Fatal(err)
		}

		minioClient, err = newMinioClient(zone)
		if err != nil {
			log.Fatal(err)
		}

		if err := minioClient.RemoveBucket(args[0]); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	sosCmd.AddCommand(sosDeleteCmd)
}
