package cmd

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh <vm name | id>",
	Short: "SSH into a virtual machine instance",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return cmd.Usage()
		}

		ipv6, err := cmd.Flags().GetBool("ipv6")
		if err != nil {
			return err
		}

		isInfo, err := cmd.Flags().GetBool("info")
		if err != nil {
			return err
		}

		isConnectionSTR, err := cmd.Flags().GetBool("print")
		if err != nil {
			return err
		}

		info, err := getSSHInfo(args[0], ipv6)
		if err != nil {
			return err
		}

		if isConnectionSTR {
			printSSHConnectSTR(info)
			return nil
		}

		if isInfo {
			printSSHInfo(info)
			return nil
		}
		return connectSSH(info)
	},
}

type sshInfo struct {
	sshKeys  string
	userName string
	ip       *net.IP
	vmName   string
	vmID     string
}

func getSSHInfo(name string, isIpv6 bool) (*sshInfo, error) {
	vm, err := getVMWithNameOrID(cs, name)
	if err != nil {
		return nil, err
	}

	sshKeyPath := path.Join(gConfigFolder, "instances", vm.ID, "id_rsa")

	if _, err := os.Stat(sshKeyPath); os.IsNotExist(err) {
		sshKeyPath = "Default ssh keypair not found"
	}

	nic := vm.DefaultNic()
	if nic == nil {
		return nil, fmt.Errorf("this instance %q has no default NIC", vm.ID)
	}

	vmIP := vm.IP()

	if isIpv6 {
		if nic.IP6Address == nil {
			return nil, fmt.Errorf("missing IPv6 address on the instance %q", vm.ID)
		}
		vmIP = &nic.IP6Address
	}

	template := &egoscale.Template{ID: vm.TemplateID, IsFeatured: true, ZoneID: "1"}

	if err := cs.Get(template); err != nil {
		return nil, err
	}

	tempUser, ok := template.Details["username"]
	if !ok {
		return nil, fmt.Errorf("missing username information in Template %q", template.ID)
	}

	return &sshInfo{
		sshKeys:  sshKeyPath,
		userName: tempUser,
		ip:       vmIP,
		vmName:   vm.Name,
		vmID:     vm.ID,
	}, nil

}

func printSSHConnectSTR(info *sshInfo) {

	if _, err := os.Stat(info.sshKeys); os.IsNotExist(err) {
		log.Println("Warning: Identity file Default ssh keypair not found not accessible: No such file or directory.")
	}

	fmt.Printf("ssh -i %s %s@%s\n", info.sshKeys, info.userName, info.ip)
}

func printSSHInfo(info *sshInfo) {
	if _, err := os.Stat(info.sshKeys); os.IsNotExist(err) {
		log.Println("Warning: Identity file Default ssh keypair not found not accessible: No such file or directory.")
	}
	println("Host", info.vmName)
	println("\tHostName", info.ip.String())
	println("\tUser", info.userName)
	println("\tIdentityFile", info.sshKeys)
}

func connectSSH(cred *sshInfo) error {

	args := []string{
		"-i",
		cred.sshKeys,
		cred.userName + "@" + cred.ip.String(),
	}

	cmd := exec.Command("ssh", args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	return cmd.Run()
}

func init() {
	sshCmd.Flags().BoolP("info", "i", false, "Print SSH config information")
	sshCmd.Flags().BoolP("print", "p", false, "Print SSH command")
	sshCmd.Flags().BoolP("ipv6", "6", false, "Use IPv6 for SSH")
	RootCmd.AddCommand(sshCmd)
}
