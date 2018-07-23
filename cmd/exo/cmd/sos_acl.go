package cmd

import (
	"fmt"
	"strings"

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
	manualRead        string = "X-Amz-Grant-Read"
	manualWrite       string = "X-Amz-Grant-Write"
	manualReadACP     string = "X-Amz-Grant-Read-Acp"
	manualWriteACP    string = "X-Amz-Grant-Write-Acp"
	manualFullControl string = "X-Amz-Grant-Full-Control"
)

// aclCmd represents the acl command
var sosACLCmd = &cobra.Command{
	Use:   "acl <bucket name> <object name> [object name] ...",
	Short: "Object(s) ACLs managment",
}

func init() {
	sosCmd.AddCommand(sosACLCmd)
}

// aclCmd represents the acl command
var sosAddACLCmd = &cobra.Command{
	Use:   "add <bucket name> <object name> [object name] ...",
	Short: "Add ACL(s) to object(s)",
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

		objInfo, err := minioClient.GetObjectACL(args[0], args[1])
		if err != nil {
			return err
		}

		src := minio.NewSourceInfo(args[0], args[1], nil)

		_, okMeta := meta["X-Amz-Acl"]
		_, okHeader := objInfo.Metadata["X-Amz-Acl"]

		if okHeader && !okMeta {
			objInfo.Metadata.Del("X-Amz-Acl")
			objInfo.Metadata.Add(manualFullControl, "id="+gCurrentAccount.Account)
		}

		src.Headers = objInfo.Metadata

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
		meta["X-Amz-Acl"] = defACL
		return meta, nil
	}

	manualACLs, err := getManualACL(cmd)
	if err != nil {
		return nil, err
	}

	if manualACLs == nil {
		return nil, nil
	}

	for k, v := range manualACLs {

		for i := range v {
			v[i] = fmt.Sprintf("id=%s", v[i])
		}

		meta[k] = strings.Join(v, ", ")
	}

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

func getManualACL(cmd *cobra.Command) (map[string][]string, error) {

	res := map[string][]string{}

	acl, err := cmd.Flags().GetString("read")
	if err != nil {
		return nil, err
	}
	if acl != "" {
		res[manualRead] = getCommaflag(acl)
	}

	acl, err = cmd.Flags().GetString("write")
	if err != nil {
		return nil, err
	}
	if acl != "" {
		res[manualWrite] = getCommaflag(acl)
	}

	acl, err = cmd.Flags().GetString("read-acp")
	if err != nil {
		return nil, err
	}
	if acl != "" {
		res[manualReadACP] = getCommaflag(acl)
	}

	acl, err = cmd.Flags().GetString("write-acp")
	if err != nil {
		return nil, err
	}
	if acl != "" {
		res[manualWriteACP] = getCommaflag(acl)
	}

	acl, err = cmd.Flags().GetString("full-control")
	if err != nil {
		return nil, err
	}
	if acl != "" {
		res[manualFullControl] = getCommaflag(acl)
	}

	if len(res) == 0 {
		return nil, nil
	}

	return res, err
}

func init() {
	sosACLCmd.AddCommand(sosAddACLCmd)
	sosAddACLCmd.Flags().SortFlags = false

	//Canned ACLs
	sosAddACLCmd.Flags().BoolP(private, "p", false, "Canned ACL private")
	sosAddACLCmd.Flags().BoolP(publicRead, "r", false, "Canned ACL public read")
	sosAddACLCmd.Flags().BoolP(publicReadWrite, "w", false, "Canned ACL public read and write")
	sosAddACLCmd.Flags().BoolP(authenticatedRead, "", false, "Canned ACL authenticated read")
	sosAddACLCmd.Flags().BoolP(bucketOwnerRead, "", false, "Canned ACL bucket owner read")
	sosAddACLCmd.Flags().BoolP(bucketOwnerFullControl, "f", false, "Canned ACL bucket owner full control")

	//Manual ACLs
	sosAddACLCmd.Flags().StringP("read", "", "", "Manual acl edit grant read e.g(value, value, ...)")
	sosAddACLCmd.Flags().StringP("write", "", "", "Manual acl edit grant write e.g(value, value, ...)")
	sosAddACLCmd.Flags().StringP("read-acp", "", "", "Manual acl edit grant acp read e.g(value, value, ...)")
	sosAddACLCmd.Flags().StringP("write-acp", "", "", "Manual acl edit grant acp write e.g(value, value, ...)")
	sosAddACLCmd.Flags().StringP("full-control", "", "", "Manual acl edit grant full control e.g(value, value, ...)")
}

// aclCmd represents the acl command
var sosRemoveACLCmd = &cobra.Command{
	Use:     "remove <bucket name> <object name> [object name] ...",
	Short:   "Remove ACL(s) from object(s)",
	Aliases: gRemoveAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return cmd.Usage()
		}

		meta, err := getManualACLBool(cmd)
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

		objInfo, err := minioClient.GetObjectACL(args[0], args[1])
		if err != nil {
			return err
		}

		src := minio.NewSourceInfo(args[0], args[1], nil)

		_, okHeader := objInfo.Metadata["X-Amz-Acl"]

		if okHeader {
			return fmt.Errorf("Error: No Manual ACL are set")
		}

		for _, k := range meta {
			objInfo.Metadata[k] = []string{""}
		}

		src.Headers = objInfo.Metadata

		// Destination object
		dst, err := minio.NewDestinationInfo(args[0], args[1], nil, nil)
		if err != nil {
			return err
		}

		// Copy object call
		return minioClient.CopyObject(dst, src)
	},
}

func init() {
	sosACLCmd.AddCommand(sosRemoveACLCmd)
	sosRemoveACLCmd.Flags().SortFlags = false
	sosRemoveACLCmd.Flags().BoolP("read", "r", false, "Remove grant read ACL")
	sosRemoveACLCmd.Flags().BoolP("write", "w", false, "Remove grant write ACL")
	sosRemoveACLCmd.Flags().BoolP("read-acp", "", false, "Remove grant acp read ACL")
	sosRemoveACLCmd.Flags().BoolP("write-acp", "", false, "Remove grant acp write ACL")
	sosRemoveACLCmd.Flags().BoolP("full-control", "f", false, "Remove grant full control ACL")
}

func getManualACLBool(cmd *cobra.Command) ([]string, error) {

	var res []string

	acl, err := cmd.Flags().GetBool("read")
	if err != nil {
		return nil, err
	}
	if acl {
		res = append(res, manualRead)
	}

	acl, err = cmd.Flags().GetBool("write")
	if err != nil {
		return nil, err
	}
	if acl {
		res = append(res, manualWrite)
	}

	acl, err = cmd.Flags().GetBool("read-acp")
	if err != nil {
		return nil, err
	}
	if acl {
		res = append(res, manualReadACP)
	}

	acl, err = cmd.Flags().GetBool("write-acp")
	if err != nil {
		return nil, err
	}
	if acl {
		res = append(res, manualWriteACP)
	}

	acl, err = cmd.Flags().GetBool("full-control")
	if err != nil {
		return nil, err
	}
	if acl {
		res = append(res, manualFullControl)
	}

	return res, nil
}

// aclCmd represents the acl command
var sosShowACLCmd = &cobra.Command{
	Use:     "show <bucket name> <object name> [object name] ...",
	Short:   "show Object ACLs",
	Aliases: gShowAlias,
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	sosACLCmd.AddCommand(sosShowACLCmd)
}
