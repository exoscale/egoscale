package cmd

import (
	"fmt"
	"log"
	"os"

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

		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			fmt.Fprintf(os.Stderr, "XXX waiting for all zone support SOS: current zone is %q. use a supported defaultZone", gCurrentAccount.DefaultZone) // nolint: errcheck
			log.Fatal(err)
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
