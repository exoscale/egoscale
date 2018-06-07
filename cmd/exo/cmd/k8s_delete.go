package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"
)

// k8sDeleteCmd represents the delete command
var k8sDeleteCmd = &cobra.Command{
	Use:   "delete <cluster name>",
	Short: "Remove kubernetes cluster",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		if err := deleteCluster(args[0]); err != nil {
			log.Fatal(err)
		}
		println(args[0])
	},
}

func deleteCluster(clusterName string) error {

	filePath := path.Join(configFolder, "k8s", "clusters", clusterName, "nodes")

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		log.Fatalf("kubernetes cluster %q not found", clusterName)
	}

	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	nodes := strings.Split(string(dat), "\n")

	if err := destroyAllNodes(nodes); err != nil {
		return err
	}

	folder := path.Join(configFolder, "k8s", "clusters", clusterName)

	return os.RemoveAll(folder)

}

func init() {
	k8sCmd.AddCommand(k8sDeleteCmd)
}
