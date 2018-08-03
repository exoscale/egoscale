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

			vm, err := getVMWithNameOrID(cs, v)
			if err != nil {
				return err
			}

			fmt.Printf("Stoping %q ", vm.Name)
			if err := vm.asyncStopWithCtx(gContext, *cs, true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("\nStopped !")
			}
		}
		return nil
	},
}

// Stop a virtual machine instance
func (vm virtualMachine) stop(cs egoscale.Client) error {
	return vm.stopWithCtx(context.Background(), cs)
}

// StopWithCtx a virtual machine instance
func (vm virtualMachine) stopWithCtx(ctx context.Context, cs egoscale.Client) error {
	_, err := cs.RequestWithContext(ctx, &egoscale.StopVirtualMachine{ID: vm.ID})
	if err != nil {
		return err
	}
	return nil
}

// AsyncStop stop a virtual machine instance Async
func (vm virtualMachine) asyncStop(cs egoscale.Client, displayLoad bool) error {
	return vm.asyncStopWithCtx(context.Background(), cs, displayLoad)
}

// AsyncStopWithCtx stop a virtual machine instance Async
func (vm virtualMachine) asyncStopWithCtx(ctx context.Context, cs egoscale.Client, displayLoad bool) error {
	var errorReq error
	cs.AsyncRequest(&egoscale.StopVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {
		if displayLoad {
			fmt.Printf(".")
		}

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
