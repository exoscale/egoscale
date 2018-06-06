package cmd

import (
	"github.com/spf13/cobra"
)

// k8sCmd represents the k8s command
var k8sCmd = &cobra.Command{
	Use:   "k8s",
	Short: "Manage kubernetes clusters",
}

func init() {
	rootCmd.AddCommand(k8sCmd)
}

const (
	k8sYMLHeader = `---
nodes:
`

	k8sNodeYML = `  - address: %s
    user: %s
    role: [controlplane,etcd,worker]
    ssh_key_path: %s
`

	k8sClusterYML = `

services:
  etcd:
    image: quay.io/coreos/etcd:latest
  kubelet:
    image: rancher/k8s:v1.10.0-rancher1-2
  kube-api:
    image: rancher/k8s:v1.10.0-rancher1-2
  kube-controller:
    image: rancher/k8s:v1.10.0-rancher1-2
  scheduler:
    image: rancher/k8s:v1.10.0-rancher1-2
  kubeproxy:
    image: rancher/k8s:v1.10.0-rancher1-2
`
)
