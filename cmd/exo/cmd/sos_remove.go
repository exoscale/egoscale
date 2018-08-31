package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

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

		recursive, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		objectsCh := make(chan string)

		// Send object names that are needed to be removed to objectsCh
		go func() {
			defer close(objectsCh)
			// List all objects from a bucket-name with a matching prefix.

			for _, arg := range args[1:] {
				for object := range minioClient.ListObjects(args[0], arg, true, nil) {
					if object.Err != nil {
						log.Fatalln(object.Err)
					}

					obj := filepath.ToSlash(object.Key)
					arg = filepath.ToSlash(arg)
					arg = strings.Trim(arg, "/")

					if strings.HasPrefix(obj, fmt.Sprintf("%s/", arg)) && obj != arg {
						if !recursive {
							fmt.Fprintf(os.Stderr, "%s: is a directory\n", arg) // nolint: errcheck
							break
						}
						objectsCh <- object.Key
					} else if obj == arg {
						objectsCh <- object.Key
					}
				}
			}
		}()

		for objectErr := range minioClient.RemoveObjectsWithContext(gContext, args[0], objectsCh) {
			return fmt.Errorf("error detected during deletion: %v", objectErr)
		}

		return nil
	},
}

func init() {
	sosCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolP("recursive", "r", false, "Attempt to remove the file hierarchy rooted in each file argument")
}
