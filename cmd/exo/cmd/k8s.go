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
	RootCmd.AddCommand(k8sCmd)
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

	cloudINIT = `#cloud-config

manage_etc_hosts: true

apt_sources:
- source: "deb [arch=amd64] https://download.docker.com/linux/ubuntu xenial stable"
  keyid: 9DC858229FC7DD38854AE2D88D81803C0EBFCD88

package_update: true
package_upgrade: true

packages:
- [ docker-ce, "17.03.2~ce-0~ubuntu-xenial" ]

runcmd:
  - [ usermod, -aG, docker, ubuntu ]
`

	clusterFileName = "cluster.yml"
)
