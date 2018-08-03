package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// rebootCmd represents the reboot command
var vmRebootCmd = &cobra.Command{
	Use:   "reboot <vm name> [vm name] ...",
	Short: "Reboot virtual machine instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		for _, v := range args {
			vm, err := getVMWithNameOrID(cs, v)
			if err != nil {
				return err
			}

			fmt.Printf("Rebooting %q ", vm.Name)
			if err := vm.asyncStopWithCtx(gContext, *cs, true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("\nRebooted !")
			}
		}
		return nil
	},
}

// Stop a virtual machine instance
func (vm virtualMachine) reboot(cs egoscale.Client) error {
	return vm.stopWithCtx(context.Background(), cs)
}

// StopWithCtx a virtual machine instance
func (vm virtualMachine) rebootWithCtx(ctx context.Context, cs egoscale.Client) error {
	_, err := cs.RequestWithContext(ctx, &egoscale.RebootVirtualMachine{ID: vm.ID})
	if err != nil {
		return err
	}
	return nil
}

// AsyncStop stop a virtual machine instance Async
func (vm virtualMachine) asyncReboot(cs egoscale.Client, displayLoad bool) error {
	return vm.asyncStopWithCtx(context.Background(), cs, displayLoad)
}

// AsyncStopWithCtx stop a virtual machine instance Async
func (vm virtualMachine) asyncRebootWithCtx(ctx context.Context, cs egoscale.Client, displayLoad bool) error {
	var errorReq error
	cs.AsyncRequest(&egoscale.RebootVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {
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
	vmCmd.AddCommand(vmRebootCmd)
}
