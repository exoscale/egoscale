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
			if err := asyncRebootWithCtx(gContext, v); err != nil {
				fmt.Fprintln(os.Stderr, err) // nolint: errcheck
			} else {
				fmt.Println("\nRebooted !")
			}
		}
		return nil
	},
}

// AsyncRebootWithCtx reboot a virtual machine instance Async
func asyncRebootWithCtx(ctx context.Context, vmName string) error {
	vm, err := getVMWithNameOrID(cs, vmName)
	if err != nil {
		return err
	}

	fmt.Printf("Rebooting %q ", vm.Name)
	var errorReq error
	cs.AsyncRequest(&egoscale.RebootVirtualMachine{ID: vm.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {

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
	vmCmd.AddCommand(vmRebootCmd)
}
