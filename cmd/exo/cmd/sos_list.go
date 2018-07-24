package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/exoscale/egoscale/cmd/exo/table"

	"github.com/minio/minio-go"

	"github.com/spf13/cobra"
)

// sosListCmd represents the list command
var sosListCmd = &cobra.Command{
	Use:     "list [<bucket name>/path]",
	Short:   "List file and folder",
	Aliases: gListAlias,
	RunE: func(cmd *cobra.Command, args []string) error {

		///XXX you must use a default zone support SOS
		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			println("XXX wainting for all zone supporting SOS: use a supported defaultZone")
			log.Fatal(err)
		}

		isRec, err := cmd.Flags().GetBool("recursive")
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return displayBucket(minioClient, isRec)
		}

		path := filepath.ToSlash(args[0])
		path = strings.Trim(path, "/")
		p := splitPath(args[0])

		if len(p) == 0 {
			return displayBucket(minioClient, isRec)
		}
		if p[0] == "" {
			return displayBucket(minioClient, isRec)
		}

		var prefix string
		if len(p) > 1 {
			prefix = path[len(p[0]):]
			prefix = strings.Trim(prefix, "/")
		}

		bucketName := p[0]

		///XXX waiting for pithos 301 redirect
		location, err := minioClient.GetBucketLocation(bucketName)
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(location)
		if err != nil {
			return err
		}
		///

		table := table.NewTable(os.Stdout)
		table.RemoveFrame()

		if isRec {
			listRecursively(minioClient, bucketName, prefix, "", false, table)
			table.Render()
			return nil
		}

		doneCh := make(chan struct{})
		defer close(doneCh)
		recursive := true

		for message := range minioClient.ListObjectsV2(bucketName, prefix, recursive, doneCh) {
			sPrefix := splitPath(prefix)
			sKey := splitPath(message.Key)
			if sPrefix[0] == "" && len(sKey) > 1 {
				continue
			}
			if len(sKey) > len(sPrefix)+1 {
				continue
			}
			if isPrefix(prefix, message.Key) {
				continue
			}
			size := fmt.Sprintf("%dB", message.Size)
			lastModified := fmt.Sprintf("%v", message.LastModified)
			key := filepath.ToSlash(message.Key)
			key = strings.TrimLeft(message.Key[len(prefix):], "/")

			table.AppendArgs(fmt.Sprintf("[%s]    ", lastModified), size, key)
		}

		table.Render()

		return nil

	},
}

func listRecursively(c *minio.Client, bucketName, prefix, zone string, displayBucket bool, table *table.Table) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	table.RemoveFrame()

	for message := range c.ListObjectsV2(bucketName, prefix, true, doneCh) {

		size := fmt.Sprintf("%dB", message.Size)
		lastModified := fmt.Sprintf("%v", message.LastModified)
		if displayBucket {
			table.AppendArgs(fmt.Sprintf("[%s]", lastModified),
				fmt.Sprintf("[%s]    ", zone),
				size,
				fmt.Sprintf("%s/%s", bucketName, message.Key))
		} else {
			table.AppendArgs(fmt.Sprintf("[%s]    ", lastModified),
				size,
				fmt.Sprintf("%s/%s", bucketName, message.Key))
		}
	}
}

func displayBucket(minioClient *minio.Client, isRecursive bool) error {
	allBuckets, err := listBucket(minioClient)
	if err != nil {
		return err
	}

	table := table.NewTable(os.Stdout)
	table.RemoveFrame()

	for zoneName, buckets := range allBuckets {
		for _, bucket := range buckets {
			if isRecursive {
				///XXX Waiting for pithos 301 redirect
				minioClient, err = newMinioClient(zoneName)
				if err != nil {
					return err
				}
				///
				listRecursively(minioClient, bucket.Name, "", zoneName, true, table)
				continue
			}
			table.AppendArgs(fmt.Sprintf("[%s]",
				bucket.CreationDate.String()),
				fmt.Sprintf("[%s]    ", zoneName),
				"0B",
				bucket.Name)
		}
	}
	table.Render()
	return nil
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

func isPrefix(prefix, file string) bool {
	prefix = strings.Trim(prefix, "/")
	file = strings.Trim(file, "/")
	return prefix == file
}

func splitPath(s string) []string {
	path := filepath.ToSlash(s)
	path = strings.Trim(path, "/")
	return strings.Split(path, "/")
}

func init() {
	sosCmd.AddCommand(sosListCmd)
	sosListCmd.Flags().BoolP("recursive", "r", false, "List recursively")
}
