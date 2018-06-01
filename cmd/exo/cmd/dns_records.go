package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ACmd represents the A command
var ACmd = &cobra.Command{
	Use:   "A",
	Short: "Add A record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("A called")
	},
}

func init() {
	dnsAddCmd.AddCommand(ACmd)
}

// AAAACmd represents the AAAA command
var AAAACmd = &cobra.Command{
	Use:   "AAAA",
	Short: "Add AAAA record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("AAAA called")
	},
}

func init() {
	dnsAddCmd.AddCommand(AAAACmd)
}

// ALIASCmd represents the ALIAS command
var ALIASCmd = &cobra.Command{
	Use:   "ALIAS",
	Short: "Add ALIAS record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ALIAS called")
	},
}

func init() {
	dnsAddCmd.AddCommand(ALIASCmd)
}

// CNAMECmd represents the CNAME command
var CNAMECmd = &cobra.Command{
	Use:   "CNAME",
	Short: "Add CNAME record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CNAME called")
	},
}

func init() {
	dnsAddCmd.AddCommand(CNAMECmd)

}

// HINFOCmd represents the HINFO command
var HINFOCmd = &cobra.Command{
	Use:   "HINFO",
	Short: "Add HINFO record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("HINFO called")
	},
}

func init() {
	dnsAddCmd.AddCommand(HINFOCmd)
}

// MXCmd represents the MX command
var MXCmd = &cobra.Command{
	Use:   "MX",
	Short: "Add MX record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("MX called")
	},
}

func init() {
	dnsAddCmd.AddCommand(MXCmd)
}

// NAPTRCmd represents the NAPTR command
var NAPTRCmd = &cobra.Command{
	Use:   "NAPTR",
	Short: "Add NAPTR record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NAPTR called")
	},
}

func init() {
	dnsAddCmd.AddCommand(NAPTRCmd)
}

// NSCmd represents the NS command
var NSCmd = &cobra.Command{
	Use:   "NS",
	Short: "Add NS record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("NS called")
	},
}

func init() {
	dnsAddCmd.AddCommand(NSCmd)
}

// POOLCmd represents the POOL command
var POOLCmd = &cobra.Command{
	Use:   "POOL",
	Short: "Add POOL record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("POOL called")
	},
}

func init() {
	dnsAddCmd.AddCommand(POOLCmd)
}

// SPFCmd represents the SPF command
var SPFCmd = &cobra.Command{
	Use:   "SPF",
	Short: "Add SPF record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SPF called")
	},
}

func init() {
	dnsAddCmd.AddCommand(SPFCmd)
}

// SRVCmd represents the SRV command
var SRVCmd = &cobra.Command{
	Use:   "SRV",
	Short: "Add SRV record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SRV called")
	},
}

func init() {
	dnsAddCmd.AddCommand(SRVCmd)
}

// SSHFPCmd represents the SSHFP command
var SSHFPCmd = &cobra.Command{
	Use:   "SSHFP",
	Short: "Add SSHFP record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("SSHFP called")
	},
}

func init() {
	dnsAddCmd.AddCommand(SSHFPCmd)
}

// TXTCmd represents the TXT command
var TXTCmd = &cobra.Command{
	Use:   "TXT",
	Short: "Add TXT record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("TXT called")
	},
}

func init() {
	dnsAddCmd.AddCommand(TXTCmd)
}

// URLCmd represents the URL command
var URLCmd = &cobra.Command{
	Use:   "URL",
	Short: "Add URL record type",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("URL called")
	},
}

func init() {
	dnsAddCmd.AddCommand(URLCmd)
}
