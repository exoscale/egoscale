package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"strings"

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

		location, err := minioClient.GetBucketLocation(bucketName)
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(location)
		if err != nil {
			return err
		}

		if isRec {
			listRecursively(minioClient, bucketName, prefix, "", false)
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
			fmt.Printf("[%s]    %s %s\n", lastModified, size, key)
		}

		return nil

	},
}

func listRecursively(c *minio.Client, bucketName, prefix, zone string, displayBucket bool) {
	doneCh := make(chan struct{})
	defer close(doneCh)

	for message := range c.ListObjectsV2(bucketName, prefix, true, doneCh) {

		size := fmt.Sprintf("%dB", message.Size)
		lastModified := fmt.Sprintf("%v", message.LastModified)
		if displayBucket {
			fmt.Printf("[%s] [%s]%s    %s %s/%s\n", lastModified, zone, alignZone(zone), size, bucketName, message.Key)
		} else {
			fmt.Printf("[%s]    %s %s\n", lastModified, size, message.Key)
		}
	}
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

func displayBucket(minioClient *minio.Client, isRecursive bool) error {
	allBuckets, err := listBucket(minioClient)
	if err != nil {
		return err
	}

	for zoneName, buckets := range allBuckets {
		for _, bucket := range buckets {
			if isRecursive {
				///XXX Waiting for pithos 301 redirect
				minioClient, err = newMinioClient(zoneName)
				if err != nil {
					return err
				}
				///
				listRecursively(minioClient, bucket.Name, "", zoneName, true)
				continue
			}
			fmt.Println(
				fmt.Sprintf("[%s]", bucket.CreationDate.String()),
				fmt.Sprintf("[%s]%s     ", zoneName, alignZone(zoneName)),
				"0B",
				bucket.Name)
		}
	}
	return nil
}

func alignZone(z string) string {
	len := 8 - len(z)
	if len < 0 {
		len = 0
	}
	return strings.Repeat(" ", len)
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
	sosListCmd.Flags().BoolP("recursive", "r", false, "List recursively")
}
