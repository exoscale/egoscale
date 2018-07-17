package cmd

import (
	"github.com/spf13/cobra"
)

// sosDeleteCmd represents the delete command
var sosDeleteCmd = &cobra.Command{
	Use:     "delete <name>",
	Short:   "Delete a bucket",
	Aliases: gDeleteAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		///XXX you must use a default zone support SOS
		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			println("XXX wainting for all zone supporting SOS: use a supported defaultZone")
			return err
		}

		zone, err := minioClient.GetBucketLocation(args[0])
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(zone)
		if err != nil {
			return err
		}

		return minioClient.RemoveBucket(args[0])
	},
}

func init() {
	sosCmd.AddCommand(sosDeleteCmd)
}
