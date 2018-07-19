package cmd

import (
	"fmt"

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
	Use:   "show",
	Short: "Show object headers",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("headers called")
	},
}

func init() {
	sosHeadersCmd.AddCommand(sosShowHeadersCmd)
}
