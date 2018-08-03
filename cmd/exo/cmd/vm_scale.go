package cmd

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// scaleCmd represents the scale command
var vmScaleCmd = &cobra.Command{
	Use:   "scale <vm name> [vm name] ...",
	Short: "Scale virtual machine",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		so, err := cmd.Flags().GetString("service-offering")
		if err != nil {
			return err
		}

		serviceoffering, err := getServiceOfferingByName(cs, so)
		if err != nil {
			return err
		}

		for _, v := range args {
			if err := asyncScaleWithCtx(gContext, v, serviceoffering.ID); err != nil {
				fmt.Fprintf(os.Stderr, "\n%v\n", err) // nolint: errcheck
			} else {
				fmt.Println("\nScaled !")
			}
		}
		return nil
	},
}

// AsyncscaleWithCtx scale a virtual machine instance Async with context
func asyncScaleWithCtx(ctx context.Context, vmName, serviceofferingID string) error {
	vm, err := getVMWithNameOrID(cs, vmName)
	if err != nil {
		return err
	}

	if egoscale.VirtualMachineState(vm.State) != egoscale.VirtualMachineStopped {
		return fmt.Errorf("this operation is not permitted if your VM is running")
	}

	fmt.Printf("Scaling %q ", vm.Name)
	var errorReq error
	cs.AsyncRequest(&egoscale.ScaleVirtualMachine{ID: vm.ID, ServiceOfferingID: serviceofferingID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {

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
	vmCmd.AddCommand(vmScaleCmd)
	vmScaleCmd.Flags().StringP("service-offering", "o", "", "<name | id> (micro|tiny|small|medium|large|extra-large|huge|mega|titan")
	if err := vmScaleCmd.MarkFlagRequired("service-offering"); err != nil {
		log.Fatal(err)
	}
}
