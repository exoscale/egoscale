package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// dissociateCmd represents the dissociate command
var dissociateCmd = &cobra.Command{
	Use:   "dissociate",
	Short: "Dissociate a private network from an instance",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("dissociate called")
	},
}

func init() {
	privnetCmd.AddCommand(dissociateCmd)
}
