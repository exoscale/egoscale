package cmd

import (
	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

const (
	private                string = "private"
	publicRead             string = "public-read"
	publicReadWrite        string = "public-read-write"
	authenticatedRead      string = "authenticated-read"
	bucketOwnerRead        string = "bucket-owner-read"
	bucketOwnerFullControl string = "bucket-owner-full-control"
)

// aclCmd represents the acl command
var sosACLCmd = &cobra.Command{
	Use:   "acl <bucket name> <object name> [object name] ...",
	Short: "Set acl on object(s) in bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		meta := map[string]string{
			//"x-amz-grant-read": "pej@exoscale.ch",
			//"x-amz-grant-read": "bre@exoscale.ch,pej@exoscale.ch",
			// "x-amz-grant-read-acp":     "",
			// "x-amz-grant-write-acp":    "",
			//"x-amz-grant-full-control": "exoscale-1",
			//"x-amz-acl": "public-read",
		}

		defACL, err := getDefaultACL(cmd)
		if err != nil {
			return err
		}

		//WIP
		if defACL == "" {
			println("Error: Choose one default Quick ACL flag")
			return cmd.Usage()
		}

		meta["x-amz-acl"] = defACL

		minioClient, err := newMinioClient(gCurrentAccount.DefaultZone)
		if err != nil {
			return err
		}

		location, err := minioClient.GetBucketLocation(args[0])
		if err != nil {
			return err
		}

		minioClient, err = newMinioClient(location)
		if err != nil {
			return err
		}

		src := minio.NewSourceInfo(args[0], args[1], nil)

		// Destination object
		dst, err := minio.NewDestinationInfo(args[0], args[1], nil, meta)
		if err != nil {
			return err
		}

		// Copy object call
		err = minioClient.CopyObject(dst, src)
		if err != nil {
			return err
		}

		return nil
	},
}

func getDefaultACL(cmd *cobra.Command) (string, error) {

	acl, err := cmd.Flags().GetBool(private)
	if err != nil {
		return "", err
	}
	if acl {
		return private, nil
	}

	acl, err = cmd.Flags().GetBool(publicRead)
	if err != nil {
		return "", err
	}
	if acl {
		return publicRead, nil
	}

	acl, err = cmd.Flags().GetBool(publicReadWrite)
	if err != nil {
		return "", err
	}
	if acl {
		return publicReadWrite, nil
	}

	acl, err = cmd.Flags().GetBool(authenticatedRead)
	if err != nil {
		return "", err
	}
	if acl {
		return authenticatedRead, nil
	}

	acl, err = cmd.Flags().GetBool(bucketOwnerRead)
	if err != nil {
		return "", err
	}
	if acl {
		return bucketOwnerRead, nil
	}

	acl, err = cmd.Flags().GetBool(bucketOwnerFullControl)
	if err != nil {
		return "", err
	}
	if acl {
		return bucketOwnerFullControl, nil
	}

	return "", nil
}

func init() {
	sosCmd.AddCommand(sosACLCmd)
	sosACLCmd.Flags().BoolP(private, "p", false, "Quick ACL private")
	sosACLCmd.Flags().BoolP(publicRead, "r", false, "Quick ACL public read")
	sosACLCmd.Flags().BoolP(publicReadWrite, "w", false, "Quick ACL public read and write")
	sosACLCmd.Flags().BoolP(authenticatedRead, "", false, "Quick ACL authenticated read")
	sosACLCmd.Flags().BoolP(bucketOwnerRead, "", false, "Quick ACL bucket owner read")
	sosACLCmd.Flags().BoolP(bucketOwnerFullControl, "f", false, "Quick ACL bucket owner full control")
}
