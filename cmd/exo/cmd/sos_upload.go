package cmd

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

// uploadCmd represents the upload command
var sosUploadCmd = &cobra.Command{
	Use:     "upload <bucket name> <file path>",
	Short:   "Upload an object into a bucket",
	Aliases: gUploadAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			return err
		}

		// Upload the zip file
		bucketName := args[0]
		objectName := filepath.Base(args[1])
		filePath := args[1]

		file, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer file.Close() // nolint: errcheck

		// Only the first 512 bytes are used to sniff the content type.
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			return err
		}

		contentType := http.DetectContentType(buffer)

		// Upload object with FPutObject
		n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
		if err != nil {
			return err
		}

		log.Printf("Successfully uploaded %s of size %d\n", objectName, n)

		return nil
	},
}

func init() {
	sosCmd.AddCommand(sosUploadCmd)
}
