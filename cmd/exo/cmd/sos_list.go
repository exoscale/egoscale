package cmd

import (
	"log"
	"os"

	"github.com/minio/minio-go"

	"github.com/exoscale/egoscale/cmd/exo/table"
	"github.com/spf13/cobra"
)

// sosListCmd represents the list command
var sosListCmd = &cobra.Command{
	Use:     "list",
	Short:   "List all buckets",
	Aliases: gListAlias,
	Run: func(cmd *cobra.Command, args []string) {

		table := table.NewTable(os.Stdout)
		table.SetHeader([]string{"Zone", "Name", "Creation Date"})

		///XXX you must use a default zone support SOS
		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			println("XXX wainting for all zone supporting SOS: use a supported defaultZone")
			log.Fatal(err)
		}

		allBuckets, err := listBucket(minioClient)
		if err != nil {
			log.Fatal(err)
		}

		for zoneName, buckets := range allBuckets {
			zName := zoneName
			for _, bucket := range buckets {
				table.Append([]string{zName, bucket.Name, bucket.CreationDate.String()})
				zName = ""
			}
		}

		table.Render()

	},
}

func listBucket(minioClient *minio.Client) (map[string][]minio.BucketInfo, error) {
	bucketInfos, err := minioClient.ListBuckets()
	if err != nil {
		return nil, err
	}

	res := map[string][]minio.BucketInfo{}

	for _, bucketInfo := range bucketInfos {

		bucketLocation, err := minioClient.GetBucketLocation(bucketInfo.Name)
		if err != nil {
			return nil, err
		}
		if _, ok := res[bucketLocation]; !ok {
			res[bucketLocation] = []minio.BucketInfo{bucketInfo}
			continue
		}

		res[bucketLocation] = append(res[bucketLocation], bucketInfo)

	}
	return res, nil
}

func init() {
	sosCmd.AddCommand(sosListCmd)
}
