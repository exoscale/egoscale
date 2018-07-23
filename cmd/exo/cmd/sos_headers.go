package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// headersCmd represents the headers command
var sosHeadersCmd = &cobra.Command{
	Use:   "header",
	Short: "Object headers management",
}

func init() {
	sosCmd.AddCommand(sosHeadersCmd)
}

// headersCmd represents the headers command
var sosAddHeadersCmd = &cobra.Command{
	Use:   "add",
	Short: "Add an header key/value to an object",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("headers called")
	},
}

func init() {
	sosHeadersCmd.AddCommand(sosAddHeadersCmd)
}

// headersCmd represents the headers command
var sosRemoveHeadersCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an header key/value from an object",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("headers called")
	},
}

func init() {
	sosHeadersCmd.AddCommand(sosRemoveHeadersCmd)
}

// headersCmd represents the headers command
var sosShowHeadersCmd = &cobra.Command{
	Use:   "show <bucket name> <object name>",
	Short: "Show object headers",
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

		objInfo, err := minioClient.GetObjectACL(args[0], args[1])
		if err != nil {
			return err
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"File Name", "Key", "Value"})

		if objInfo.ContentType != "" {
			table.Append([]string{objInfo.Key, "content-type", objInfo.ContentType})
		}

		for k, v := range objInfo.Metadata {
			k = strings.ToLower(k)
			println(k)
			if isStandardHeader(k) && len(v) > 0 {
				table.Append([]string{objInfo.Key, k, v[0]})
			}
		}

		table.Render()

		return nil
	},
}

var supportedHeaders = []string{
	"content-type",
	"cache-control",
	"content-encoding",
	"content-disposition",
	"content-language",
	"x-amz-website-redirect-location",
	"expires",
}

func isStandardHeader(headerKey string) bool {
	key := strings.ToLower(headerKey)
	for _, header := range supportedHeaders {
		if strings.ToLower(header) == key {
			return true
		}
	}
	return false
}

func init() {
	sosHeadersCmd.AddCommand(sosShowHeadersCmd)
}
