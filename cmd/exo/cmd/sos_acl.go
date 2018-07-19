package cmd

import (
	"fmt"

	minio "github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

const (
	//Canned ACLs
	private                string = "private"
	publicRead             string = "public-read"
	publicReadWrite        string = "public-read-write"
	authenticatedRead      string = "authenticated-read"
	bucketOwnerRead        string = "bucket-owner-read"
	bucketOwnerFullControl string = "bucket-owner-full-control"

	//Manual edit ACLs
	manualRead        string = "x-amz-grant-read"
	manualWrite       string = "x-amz-grant-write"
	manualReadACP     string = "x-amz-grant-read-acp"
	manualWriteACP    string = "x-amz-grant-write-acp"
	manualFullControl string = "x-amz-grant-full-control"
)

// aclCmd represents the acl command
var sosACLCmd = &cobra.Command{
	Use:   "acl <bucket name> <object name> [object name] ...",
	Short: "Set acl on object(s) in bucket",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		meta, err := getACL(cmd)
		if err != nil {
			return err
		}

		if meta == nil {
			println("error: You have to choose one flag")
			if err = cmd.Usage(); err != nil {
				return err
			}
			return fmt.Errorf("error: You have to choose one flag")
		}

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
		return minioClient.CopyObject(dst, src)
	},
}

func getACL(cmd *cobra.Command) (map[string]string, error) {

	meta := map[string]string{}

	defACL, err := getDefaultCannedACL(cmd)
	if err != nil {
		return nil, err
	}

	if defACL != "" {
		meta["x-amz-acl"] = defACL
		return meta, nil
	}

	ManualACL, acl, err := getManualACL(cmd)
	if err != nil {
		return nil, err
	}

	meta[manualFullControl] = "id=" + gCurrentAccount.Account

	if ManualACL == "" {
		return nil, nil
	}

	meta[ManualACL] = "id=" + acl

	return meta, nil
}

func getDefaultCannedACL(cmd *cobra.Command) (string, error) {

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

func getManualACL(cmd *cobra.Command) (string, string, error) {

	acl, err := cmd.Flags().GetString("read")
	if err != nil {
		return "", "", err
	}
	if acl != "" {
		return manualRead, acl, nil
	}

	acl, err = cmd.Flags().GetString("write")
	if err != nil {
		return "", "", err
	}
	if acl != "" {
		return manualWrite, acl, nil
	}

	acl, err = cmd.Flags().GetString("read-acp")
	if err != nil {
		return "", "", err
	}
	if acl != "" {
		return manualReadACP, acl, nil
	}

	acl, err = cmd.Flags().GetString("write-acp")
	if err != nil {
		return "", "", err
	}
	if acl != "" {
		return manualWriteACP, acl, nil
	}

	acl, err = cmd.Flags().GetString("full-control")
	if err != nil {
		return "", "", err
	}
	if acl != "" {
		return manualFullControl, acl, nil
	}

	return "", "", nil
}

func init() {
	sosCmd.AddCommand(sosACLCmd)
	sosACLCmd.Flags().SortFlags = false

	//Canned ACLs
	sosACLCmd.Flags().BoolP(private, "p", false, "Canned ACL private")
	sosACLCmd.Flags().BoolP(publicRead, "r", false, "Canned ACL public read")
	sosACLCmd.Flags().BoolP(publicReadWrite, "w", false, "Canned ACL public read and write")
	sosACLCmd.Flags().BoolP(authenticatedRead, "", false, "Canned ACL authenticated read")
	sosACLCmd.Flags().BoolP(bucketOwnerRead, "", false, "Canned ACL bucket owner read")
	sosACLCmd.Flags().BoolP(bucketOwnerFullControl, "f", false, "Canned ACL bucket owner full control")

	//Manual ACLs
	sosACLCmd.Flags().StringP("read", "", "", "Manual acl edit grant read")
	sosACLCmd.Flags().StringP("write", "", "", "Manual acl edit grant write")
	sosACLCmd.Flags().StringP("read-acp", "", "", "Manual acl edit grant acp read")
	sosACLCmd.Flags().StringP("write-acp", "", "", "Manual acl edit grant acp write")
	sosACLCmd.Flags().StringP("full-control", "", "", "Manual acl edit grant full control")
}
