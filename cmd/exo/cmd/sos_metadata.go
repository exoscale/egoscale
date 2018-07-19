package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// metadataCmd represents the metadata command
var sosMetadataCmd = &cobra.Command{
	Use:   "metadata",
	Short: "Object metadate management",
}

func init() {
	sosCmd.AddCommand(sosMetadataCmd)
}

// metadataCmd represents the metadata command
var sosAddMetadataCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a metadata to an object",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("metadata called")
	},
}

func init() {
	sosMetadataCmd.AddCommand(sosAddMetadataCmd)
}

// metadataCmd represents the metadata command
var sosRemoveMetadataCmd = &cobra.Command{
	Use:     "remove",
	Aliases: gRemoveAlias,
	Short:   "Remove a metadata from an object",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("metadata called")
	},
}

func init() {
	sosMetadataCmd.AddCommand(sosRemoveMetadataCmd)
}

// metadataCmd represents the metadata command
var sosShowMetadataCmd = &cobra.Command{
	Use:     "show",
	Short:   "Show object metadatas",
	Aliases: gShowAlias,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Show object metadatas")
	},
}

func init() {
	sosMetadataCmd.AddCommand(sosShowMetadataCmd)
}
