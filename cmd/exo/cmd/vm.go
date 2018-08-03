package cmd

import (
	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// vmCmd represents the vm command
var vmCmd = &cobra.Command{
	Use:   "vm",
	Short: "Virtual machines management",
}

const (
	starting = "Starting"
	running  = "Running"
	stoping  = "Stopping"
	stopped  = "Stopped"
)

// VirtualMachine represente egoscale virtual machine
type virtualMachine struct {
	*egoscale.VirtualMachine
}

func getVMWithNameOrID(cs *egoscale.Client, name string) (*virtualMachine, error) {
	vm := &egoscale.VirtualMachine{ID: name}
	if err := cs.Get(vm); err == nil {
		return &virtualMachine{vm}, err
	}

	vm.Name = name
	vm.ID = ""

	if err := cs.Get(vm); err != nil {
		return nil, err
	}
	return &virtualMachine{vm}, nil
}

func getSecurityGroup(vm *virtualMachine) []string {
	sgs := []string{}
	for _, sgN := range vm.SecurityGroup {
		sgs = append(sgs, sgN.Name)
	}
	return sgs
}

func init() {
	RootCmd.AddCommand(vmCmd)
}
