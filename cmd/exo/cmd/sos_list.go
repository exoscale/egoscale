package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sosListCmd represents the list command
var sosListCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
	},
}

func init() {
	sosCmd.AddCommand(sosListCmd)
}
