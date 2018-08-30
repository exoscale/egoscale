package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:     "remove <bucket name> <object name> [object name] ...",
	Short:   "Remove object(s) from a bucket",
	Aliases: gRemoveAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		minioClient, err := newMinioClient(sosZone)
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

		chName := make(chan string, len(args[1:]))

		go func() {
			defer close(chName)
			for _, obj := range args[1:] {
				chName <- obj
			}
		}()

		for objectErr := range minioClient.RemoveObjectsWithContext(gContext, args[0], chName) {
			return fmt.Errorf("error detected during deletion: %v", objectErr)
		}

		log.Printf("Object(s):\n")

		for _, obj := range args[1:] {
			log.Println("-", obj)
		}

		log.Printf("Successfully removed\n")

		return nil
	},
}

func init() {
	sosCmd.AddCommand(removeCmd)
}
