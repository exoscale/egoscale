package cmd

import (
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/spf13/cobra"
)

// kubectlCmd represents the kubectl command
var kubectlCmd = &cobra.Command{
	Use:   "kubectl <cluster name> [all kubectl cmd you want]",
	Short: "Run kubectl command for a kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		configFilePath := path.Join(configFolder, "k8s", "clusters", args[0])

		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			log.Fatalf("kubernetes cluster %q not found", args[0])
		}

		configFilePath = path.Join(configFilePath, "kube_config_cluster.yml")

		if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
			log.Fatalf("kubectl cluster config file %q not found", configFilePath)
		}

		cmdArgs := []string{
			"--kubeconfig",
			configFilePath,
		}

		for _, param := range args[1:] {
			cmdArgs = append(cmdArgs, param)
		}

		c := exec.Command("kubectl", cmdArgs...)

		c.Stdin = os.Stdin
		c.Stderr = os.Stderr
		c.Stdout = os.Stdout

		if err := c.Run(); err != nil {
			log.Fatal(err)
		}

	},
}

func init() {
	RootCmd.AddCommand(kubectlCmd)
}
