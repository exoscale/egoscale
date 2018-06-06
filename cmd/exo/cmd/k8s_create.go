package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	ttime "time"

	"github.com/exoscale/egoscale"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var k8sCreateCmd = &cobra.Command{
	Use:   "create <name>",
	Short: "Create k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		node, err := cmd.Flags().GetString("node")
		if err != nil {
			log.Fatal(err)
		}

		nodes := getCommaflag(node)

		if node == "" {

			nodeCap, err := cmd.Flags().GetString("node-capacity")
			if err != nil {
				log.Fatal(err)
			}

			nodeNumber, err := cmd.Flags().GetInt("node-number")
			if err != nil {
				log.Fatal(err)
			}

			nodes, err = deployNodes(nodeNumber, nodeCap)
			if err != nil {
				log.Fatal(err)
			}

			///////MMMDDRRR J'avais pas pensé au CLOUDINIT qui job !!!!!
			////WIP crado
			ttime.Sleep(ttime.Second * 60)
		}

		clusterFile, err := createK8sClusterFile(nodes)
		if err != nil {
			log.Fatal(err)
		}

		if err := startWithRKE(clusterFile, args[0]); err != nil {
			log.Fatal(err)
		}
	},
}

func deployNodes(nodeNumber int, nodeCapacity string) ([]string, error) {

	nodes := make([]string, nodeNumber)

	for i := 0; i < nodeNumber; i++ {

		userData, err := getUserData(path.Join(configFolder, "cloudinit.yml"))
		if err != nil {
			log.Fatal(err)
		}

		zone, err := getZoneIDByName(cs, "ch-dk-2")
		if err != nil {
			log.Fatal(err)
		}

		template, err := getTemplateIDByName(cs, "Linux Ubuntu 18.04", zone)
		if err != nil {
			log.Fatal(err)
		}

		sgs, err := getSecuGrpList(cs, "default")
		if err != nil {
			log.Fatal(err)
		}

		servOffering, err := getServiceOfferingIDByName(cs, "medium")
		if err != nil {
			log.Fatal(err)
		}

		req := &egoscale.DeployVirtualMachine{
			Name:              fmt.Sprintf("test-node-%d", i),
			UserData:          userData,
			ZoneID:            zone,
			TemplateID:        template,
			RootDiskSize:      50,
			SecurityGroupIDs:  sgs,
			ServiceOfferingID: servOffering,
		}

		vm, err := createVM(req)
		if err != nil {
			return nil, err
		}

		nodes[i] = vm.ID
	}

	return nodes, nil
}

func startWithRKE(clusterFile, clusterName string) error {

	filePath := path.Join(configFolder, "k8s", "clusters", clusterName)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return err
		}
	}

	filePath = path.Join(filePath, "cluster.yml")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filePath, []byte(clusterFile), 0600); err != nil {
			log.Fatalf("cluster.yml could not be written. %s", err)
		}
	}

	args := []string{
		"up",
		"--config",
		filePath,
	}

	cmd := exec.Command(path.Join(configFolder, "rke"), args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func createK8sClusterFile(nodes []string) (string, error) {

	clusterYML := k8sYMLHeader
	for _, node := range nodes {
		vm, err := getVMWithNameOrID(cs, node)
		if err != nil {
			return "", err
		}

		temp := &egoscale.Template{IsFeatured: true, ID: vm.TemplateID, ZoneID: "1"}

		if err := cs.Get(temp); err != nil {
			return "", err
		}

		sshInfos, err := getSSHInfos(node)
		if err != nil {
			return "", err
		}

		clusterYML += fmt.Sprintf(k8sNodeYML, vm.IP().String(), sshInfos.userName, sshInfos.sshKeys)

	}
	clusterYML += k8sClusterYML

	return clusterYML, nil
}

func init() {
	k8sCmd.AddCommand(k8sCreateCmd)
	k8sCreateCmd.Flags().StringP("node", "n", "", "Node can provision existing instances [vm name | id, vm name | id,...]")
	k8sCreateCmd.Flags().IntP("node-number", "", 1, "Node number to create (if --node not set)")
	k8sCreateCmd.Flags().StringP("node-capacity", "", "medium", "Node(s) capacity (if --node not set) (micro|tiny|small|medium|large|extra-large|huge|mega|titan)")
}
