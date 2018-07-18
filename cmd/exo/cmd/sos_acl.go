package cmd

import (
	"github.com/spf13/cobra"
)

// aclCmd represents the acl command
var aclCmd = &cobra.Command{
	Use:   "acl <bucket name> <object name> [object name] ...",
	Short: "Set acl on object(s) in bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			return err
		}

		location, err := minioClient.GetBucketLocation(args[0])
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(location)
		if err != nil {
			return err
		}

		return nil
	},
}

func init() {
	sosCmd.AddCommand(aclCmd)
}
