package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var vmDeleteCmd = &cobra.Command{
	Use:     "delete <name | id>+",
	Short:   "Delete virtual machine instance(s)",
	Aliases: gDeleteAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		tasks := []task{}

		for _, arg := range args {
			tsk, err := prepareDeleteVM(&arg, force)
			if err != nil {
				return err
			}
			if tsk != nil {
				tasks = append(tasks, task{tsk, fmt.Sprintf("Destroying %q ", arg)})
			}
		}

		resps, errs := asyncTasks(tasks)
		if len(errs) > 0 {
			return errs[0]
		}

		for _, r := range resps {
			resp := r.(*egoscale.VirtualMachine)

			folder := path.Join(gConfigFolder, "instances", resp.ID)

			if _, err := os.Stat(folder); !os.IsNotExist(err) {
				if err := os.RemoveAll(folder); err != nil {
					return err
				}
			}
		}

		return nil
	},
}

func prepareDeleteVM(name *string, force bool) (*egoscale.DestroyVirtualMachine, error) {
	vm, err := getVMWithNameOrID(*name)
	if err != nil {
		return nil, err
	}

	if !force {
		if !askQuestion(fmt.Sprintf("sure you want to delete %q virtual machine", vm.Name)) {
			return nil, nil
		}
	}

	*name = vm.Name

	return &egoscale.DestroyVirtualMachine{ID: vm.ID}, nil
}

func init() {
	vmDeleteCmd.Flags().BoolP("force", "f", false, "Attempt to remove vitual machine without prompting for confirmation")
	vmCmd.AddCommand(vmDeleteCmd)
}
