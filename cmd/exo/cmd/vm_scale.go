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
			vm, err := getVMWithNameOrID(cs, v)
			if err != nil {
				return err
			}

			fmt.Printf("Scaling %q ", vm.Name)
			if err := vm.asyncScaleWithCtx(gContext, *cs, serviceoffering, true); err != nil {
				fmt.Fprintln(os.Stderr, err)
			} else {
				fmt.Println("\nScaled !")
			}
		}
		return nil
	},
}

// Scale a virtual machine instance
func (vm virtualMachine) scale(cs egoscale.Client, so *egoscale.ServiceOffering) error {
	return vm.scaleWithCtx(context.Background(), cs, so)
}

// ScaleWithCtx a virtual machine instance
func (vm virtualMachine) scaleWithCtx(ctx context.Context, cs egoscale.Client, so *egoscale.ServiceOffering) error {

	if err := vm.stopWithCtx(ctx, cs); err != nil {
		return err
	}

	_, err := cs.RequestWithContext(ctx, &egoscale.ScaleVirtualMachine{ID: vm.ID, ServiceOfferingID: so.ID})
	if err != nil {
		return err
	}
	return vm.startWithCtx(ctx, cs)
}

// AsyncScale scale a virtual machine instance Async
func (vm virtualMachine) asyncScale(cs egoscale.Client, so *egoscale.ServiceOffering, displayLoad bool) error {
	return vm.asyncScaleWithCtx(context.Background(), cs, so, displayLoad)
}

// AsyncscaleWithCtx scale a virtual machine instance Async with context
func (vm virtualMachine) asyncScaleWithCtx(ctx context.Context, cs egoscale.Client, so *egoscale.ServiceOffering, displayLoad bool) error {

	if err := vm.asyncStopWithCtx(ctx, cs, displayLoad); err != nil {
		return err
	}

	var errorReq error
	cs.AsyncRequest(&egoscale.ScaleVirtualMachine{ID: vm.ID, ServiceOfferingID: so.ID}, func(jobResult *egoscale.AsyncJobResult, err error) bool {
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
	if errorReq != nil {
		return errorReq
	}

	return vm.asyncStartWithCtx(ctx, cs, displayLoad)
}

func init() {
	vmCmd.AddCommand(vmScaleCmd)
	vmScaleCmd.Flags().StringP("service-offering", "o", "", "<name | id> (micro|tiny|small|medium|large|extra-large|huge|mega|titan")
	if err := vmScaleCmd.MarkFlagRequired("service-offering"); err != nil {
		log.Fatal(err)
	}
}
