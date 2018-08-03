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
			vm, err := getVMWithNameOrID(cs, v)
			if err != nil {
				return err
			}

			fmt.Printf("Starting %q ", vm.Name)
			if err := vm.asyncStartWithCtx(gContext, *cs, true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("\nStarted !")
			}
		}
		return nil
	},
}

// Start a virtual machine instance
func (vm virtualMachine) start(cs egoscale.Client) error {
	return vm.startWithCtx(context.Background(), cs)
}

// StartWithCtx a virtual machine instance
func (vm virtualMachine) startWithCtx(ctx context.Context, cs egoscale.Client) error {
	_, err := cs.RequestWithContext(ctx, &egoscale.StartVirtualMachine{ID: vm.ID})
	if err != nil {
		return err
	}
	return nil
}

// AsyncStart start a virtual machine instance Async
func (vm virtualMachine) asyncStart(cs egoscale.Client, displayLoad bool) error {
	return vm.asyncStartWithCtx(context.Background(), cs, displayLoad)
}

// AsyncStartWithCtx start a virtual machine instance Async
func (vm virtualMachine) asyncStartWithCtx(ctx context.Context, cs egoscale.Client, displayLoad bool) error {
	var errorReq error
	cs.AsyncRequest(&egoscale.StartVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {
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
	vmCmd.AddCommand(vmStartCmd)
}
