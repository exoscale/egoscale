package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var sosShowCmd = &cobra.Command{
	Use:     "show <bucket name> <keyword>",
	Short:   "Show bucket object(s)",
	Aliases: gShowAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		var prefix string

		if len(args) >= 2 {
			prefix = strings.Join(args[1:], "/")
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

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Name", "Size", "Last modified"})

		doneCh := make(chan struct{})
		defer close(doneCh)
		recursive := true

		for message := range minioClient.ListObjectsV2(args[0], prefix, recursive, doneCh) {
			size := fmt.Sprintf("%d", message.Size)
			lastModified := fmt.Sprintf("%v", message.LastModified)
			table.Append([]string{message.Key, size, lastModified})
		}

		table.Render()

		return nil
	},
}

func init() {
	sosCmd.AddCommand(sosShowCmd)
}
