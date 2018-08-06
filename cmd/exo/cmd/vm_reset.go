package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

var resetCmdDiskSize *int64

// vmResetCmd represents the stop command
var vmResetCmd = &cobra.Command{
	Use:   "reset <vm name> [vm name] ...",
	Short: "Reset virtual machine instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		diskValue, err := getInt64CustomFlag(cmd, "disk")
		if err != nil {
			return err
		}

		errs := []error{}
		for _, v := range args {
			if err := resetVirtualMachine(v, diskValue); err != nil {
				errs = append(errs, fmt.Errorf("could not reset %q: %s", v, err))
			}
		}

		if len(errs) == 1 {
			return errs[0]
		}
		if len(errs) > 1 {
			var b strings.Builder
			for _, err := range errs {
				if _, e := fmt.Fprintln(&b, err); e != nil {
					return e
				}
			}
			return errors.New(b.String())
		}

		return nil
	},
}

// resetVirtualMachine stop a virtual machine instance
func resetVirtualMachine(vmName string, diskValue int64PtrValue) error {
	vm, err := getVMWithNameOrID(vmName)
	if err != nil {
		return err
	}

	volume := &egoscale.Volume{
		VirtualMachineID: vm.ID,
		Type:             "ROOT",
	}

	if err := cs.Get(volume); err != nil {
		return err
	}

	rootDiskSize := int64(volume.Size >> 30)

	if diskValue.int64 != nil {
		if *diskValue.int64 < 10 {
			return fmt.Errorf("root disk size must be greater or equal than 10GB")
		}
		rootDiskSize = *diskValue.int64
	}

	fmt.Printf("Resetting %q ", vm.Name)
	var errorReq error
	cs.AsyncRequestWithContext(gContext, &egoscale.RestoreVirtualMachine{VirtualMachineID: vm.ID, RootDiskSize: rootDiskSize}, func(jobResult *egoscale.AsyncJobResult, err error) bool {

		fmt.Print(".")

		if err != nil {
			errorReq = err
			return false
		}

		if jobResult.JobStatus == egoscale.Success {
			fmt.Println(" success.")
			return false
		}

		return true
	})

	if errorReq != nil {
		fmt.Println(" failure!")
	}

	return errorReq
}

func init() {
	vmCmd.AddCommand(vmResetCmd)
	diskSizeVarP := new(int64PtrValue)
	vmResetCmd.Flags().VarP(diskSizeVarP, "disk", "d", "New disk size after reset in GB")
}
