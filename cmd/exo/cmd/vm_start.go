package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var vmStartCmd = &cobra.Command{
	Use:   "start <vm name> [vm name] ...",
	Short: "Start virtual machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		for _, v := range args {
			if err := asyncStartWithCtx(gContext, v); err != nil {
				fmt.Fprintln(os.Stderr, err) // nolint: errcheck
			} else {
				fmt.Println("\nStarted !")
			}
		}
		return nil
	},
}

// AsyncStartWithCtx start a virtual machine instance Async
func asyncStartWithCtx(ctx context.Context, vmName string) error {
	vm, err := getVMWithNameOrID(cs, vmName)
	if err != nil {
		return err
	}

	fmt.Printf("Starting %q ", vm.Name)
	var errorReq error
	cs.AsyncRequest(&egoscale.StartVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {

		fmt.Printf(".")

		if err != nil {
			errorReq = err
			return false
		}

		if jobResult.JobStatus == egoscale.Success {
			return false
		}
		return true
	})
	return errorReq
}

func init() {
	vmCmd.AddCommand(vmStartCmd)
}
