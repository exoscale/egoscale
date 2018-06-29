package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// k8sListCmd represents the list command
var k8sListCmd = &cobra.Command{
	Use:   "list",
	Short: "List kubernetes cluster(s)",
	Run: func(cmd *cobra.Command, args []string) {
		folderPath := path.Join(configFolder, "k8s", "clusters")

		files, err := ioutil.ReadDir(folderPath)
		if err != nil {
			log.Fatal(err)
		}

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Kubernetes Cluster"})

		for _, f := range files {
			table.Append([]string{f.Name()})
		}

		table.Render()
	},
}

func init() {
	k8sCmd.AddCommand(k8sListCmd)
}
