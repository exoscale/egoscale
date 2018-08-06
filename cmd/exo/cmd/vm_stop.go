package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var vmStopCmd = &cobra.Command{
	Use:   "stop <vm name> [vm name] ...",
	Short: "Stop virtual machine instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		for _, v := range args {
			if err := asyncStopWithCtx(gContext, v); err != nil {
				fmt.Fprintln(os.Stderr, err) // nolint: errcheck
			} else {
				fmt.Println("\nStopped !")
			}
		}
		return nil
	},
}

// AsyncStopWithCtx stop a virtual machine instance Async
func asyncStopWithCtx(ctx context.Context, vmName string) error {
	vm, err := getVMWithNameOrID(cs, vmName)
	if err != nil {
		return err
	}

	if egoscale.VirtualMachineState(vm.State) == egoscale.VirtualMachineStopped {
		return fmt.Errorf("virtual machine already stopped")
	}

	fmt.Printf("Stoping %q ", vm.Name)
	var errorReq error
	cs.AsyncRequest(&egoscale.StopVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {

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
	vmCmd.AddCommand(vmStopCmd)
}
