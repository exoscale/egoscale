package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var vmStopCmd = &cobra.Command{
	Use:   "stop <vm name> [vm name] ...",
	Short: "Stop a virtual machine instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		for _, v := range args {

			vm, err := getVMWithNameOrID(cs, v)
			if err != nil {
				return err
			}

			fmt.Printf("Stoping %q ", vm.Name)
			if err := vm.AsyncStopWithCtx(gContext, *cs, true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("\nstopped !")
			}
		}
		return nil
	},
}

func init() {
	vmCmd.AddCommand(vmStopCmd)
}
