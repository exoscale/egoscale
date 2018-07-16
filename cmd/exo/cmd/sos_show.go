package cmd

import (
	"fmt"
	"os"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// showCmd represents the show command
var sosShowCmd = &cobra.Command{
	Use:     "show <bucket name>",
	Short:   "Show bucket object(s)",
	Aliases: gShowAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			return err
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Name", "Size", "Last modified"})

		doneCh := make(chan struct{})
		defer close(doneCh)
		recursive := false

		for message := range minioClient.ListObjectsV2(args[0], "", recursive, doneCh) {
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
