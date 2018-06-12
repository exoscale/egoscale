package cmd

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
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

		noAuto, err := cmd.Flags().GetBool("no-auto")
		if err != nil {
			log.Fatal(err)
		}

		filePath := path.Join(configFolder, "k8s", "clusters", args[0])

		if _, err := os.Stat(filePath); !os.IsNotExist(err) {
			log.Fatalf("kubernetes cluster %q already exist", args[0])
		}

		node, err := cmd.Flags().GetString("node")
		if err != nil {
			log.Fatal(err)
		}

		nodes := getCommaflag(node)

		if node == "" {

			firewallRule, err := cmd.Flags().GetBool("firewall-rules-add")
			if err != nil {
				log.Fatal(err)
			}

			sg, err := cmd.Flags().GetString("security-group")
			if err != nil {
				log.Fatal(err)
			}

			securityGroup, err := getSecuGrpWithNameOrID(cs, sg)
			if err != nil {
				log.Fatal(err)
			}

			nodeCap, err := cmd.Flags().GetString("node-capacity")
			if err != nil {
				log.Fatal(err)
			}

			nodeNumber, err := cmd.Flags().GetInt("node-number")
			if err != nil {
				log.Fatal(err)
			}

			if firewallRule {
				if err := addK8sRules(securityGroup.ID); err != nil {
					log.Fatal(err)
				}
			}

			nodes, err = deployNodes(nodeNumber, nodeCap, securityGroup.ID)
			if err != nil {
				log.Fatal(err)
			}

			if err := checkingCloudInitJob(nodes); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
				if err := destroyAllNodes(nodes); err != nil {
					log.Fatal(err)
				}
				os.Exit(1)
			}
		}

		clusterFile, err := createK8sClusterFile(nodes, args[0])
		if err != nil {
			log.Fatal(err)
		}

		if err := storeConfig(args[0], clusterFile, nodes); err != nil {
			log.Fatal(err)
		}

		filePath = path.Join(configFolder, "k8s", "clusters", args[0], "cluster.yml")

		if noAuto {
			println(clusterFile)
			return
		}

		if err := startWithRKE(filePath); err != nil {
			log.Fatal(err)
		}

		displayHelper()
	},
}

func addK8sRules(securityGroupID string) error {

	rule, err := getDefaultRule(SSH.String(), false)
	if err != nil {
		return err
	}

	rule.SecurityGroupID = securityGroupID

	if err := addRule(rule, false); err != nil {
		return err
	}

	rule, err = getDefaultRule(ETCDCLIENT.String(), false)
	if err != nil {
		return err
	}

	rule.SecurityGroupID = securityGroupID

	if err := addRule(rule, false); err != nil {
		return err
	}

	rule, err = getDefaultRule(ETCDSERVER.String(), false)
	if err != nil {
		return err
	}

	rule.SecurityGroupID = securityGroupID

	if err := addRule(rule, false); err != nil {
		return err
	}

	rule, err = getDefaultRule(RKEHTTPS.String(), false)
	if err != nil {
		return err
	}

	rule.SecurityGroupID = securityGroupID

	if err := addRule(rule, false); err != nil {
		return err
	}

	return nil
}

func storeConfig(clusterName, clusterFile string, nodes []string) error {
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

	nodesFile := strings.Join(nodes, "\n")

	filePath = path.Join(configFolder, "k8s", "clusters", clusterName, "nodes")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filePath, []byte(nodesFile), 0600); err != nil {
			log.Fatalf("nodes could not be written. %s", err)
		}
	}
	return nil
}

func checkingCloudInitJob(vms []string) error {

	defer println("")

	print("Installing Docker on node(s)")
	for i := 0; i < 120; i++ {

		var errCMD error
		for _, vm := range vms {

			sshinfo, err := getSSHInfos(vm)
			if err != nil {
				return err
			}

			args := []string{
				fmt.Sprintf("%s@%s", sshinfo.userName, sshinfo.ip.String()),
				"-i",
				sshinfo.sshKeys,
				"-o",
				"StrictHostKeyChecking=no",
				"docker",
				"ps",
			}

			cmd := exec.Command("ssh", args...)

			errCMD = cmd.Run()
			if errCMD != nil {
				break
			}
		}
		if errCMD == nil {
			return nil
		}
		ttime.Sleep(ttime.Second)
		print(".")
		errCMD = nil
	}

	return fmt.Errorf("waiting to installing Docker Timeout")
}

func destroyAllNodes(nodes []string) error {
	for _, node := range nodes {
		if err := deleteVM(node); err != nil {
			return err
		}
	}
	return nil
}

func deployNodes(nodeNumber int, nodeCapacity, sg string) ([]string, error) {

	nodes := make([]string, nodeNumber)

	for i := 0; i < nodeNumber; i++ {

		zone, err := getZoneIDByName(cs, "ch-dk-2")
		if err != nil {
			log.Fatal(err)
		}

		template, err := getTemplateIDByName(cs, "Linux Ubuntu 18.04 LTS", zone)
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
			Name:              fmt.Sprintf("node-%d", i+1),
			UserData:          base64.StdEncoding.EncodeToString([]byte(cloudINIT)),
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

func startWithRKE(filePath string) error {

	args := []string{
		"up",
		"--config",
		filePath,
	}

	cmd := exec.Command("rke", args...)

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func createK8sClusterFile(nodes []string, clusterName string) (string, error) {

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

func displayHelper() {

}

func init() {
	k8sCmd.AddCommand(k8sCreateCmd)
	k8sCreateCmd.Flags().StringP("node", "n", "", "Node can provision existing instances [vm name | id, vm name | id,...]")
	k8sCreateCmd.Flags().IntP("node-number", "", 1, "Node number to create (if --node not set)")
	k8sCreateCmd.Flags().StringP("node-capacity", "", "medium", "Node(s) capacity (if --node not set) (micro|tiny|small|medium|large|extra-large|huge|mega|titan)")
	k8sCreateCmd.Flags().BoolP("firewall-rules-add", "f", false, "Add firewall rules for kubernetes (if --node not set)")
	k8sCreateCmd.Flags().StringP("security-group", "s", "default", "Create node(s) in a security group <security group name | id> (if --node not set)")
	k8sCreateCmd.Flags().BoolP("no-auto", "", false, "")
}
