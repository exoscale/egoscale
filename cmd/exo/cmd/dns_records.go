package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// dnsACmd represents the A command
var dnsACmd = &cobra.Command{
	Use:   "A <domain name>",
	Short: "Add A record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		addr, err := cmd.Flags().GetString("address")
		if err != nil {
			log.Fatal(err)
		}
		ttl, err := cmd.Flags().GetInt("ttl")
		if err != nil {
			log.Fatal(err)
		}

		domain, err := csDNS.GetDomain(args[0])
		if err != nil {
			log.Fatal(err)
		}

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "A",
			Name:       name,
			Content:    addr,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsACmd)
	dnsACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsACmd.Flags().StringP("address", "a", "", "Example: 127.0.0.1")
	dnsACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsACmd.MarkFlagRequired("address")
}

// AAAACmd represents the AAAA command
var dnsAAAACmd = &cobra.Command{
	Use:   "AAAA <domain name>",
	Short: "Add AAAA record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		addr, err := cmd.Flags().GetString("address")
		if err != nil {
			log.Fatal(err)
		}
		ttl, err := cmd.Flags().GetInt("ttl")
		if err != nil {
			log.Fatal(err)
		}

		domain, err := csDNS.GetDomain(args[0])
		if err != nil {
			log.Fatal(err)
		}

		resp, err := csDNS.CreateRecord(args[0], egoscale.DNSRecord{
			DomainID:   domain.ID,
			TTL:        ttl,
			RecordType: "AAAA",
			Name:       name,
			Content:    addr,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsAAAACmd)
	dnsAAAACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsAAAACmd.Flags().StringP("address", "a", "", "Example: 2001:0db8:85a3:0000:0000:EA75:1337:BEEF")
	dnsAAAACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsAAAACmd.MarkFlagRequired("address")
}

// ALIASCmd represents the ALIAS command
var dnsALIASCmd = &cobra.Command{
	Use:   "ALIAS <domain name>",
	Short: "Add ALIAS record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("ALIAS called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsALIASCmd)
}

// CNAMECmd represents the CNAME command
var dnsCNAMECmd = &cobra.Command{
	Use:   "CNAME <domain name>",
	Short: "Add CNAME record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("CNAME called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsCNAMECmd)

}

// HINFOCmd represents the HINFO command
var dnsHINFOCmd = &cobra.Command{
	Use:   "HINFO <domain name>",
	Short: "Add HINFO record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("HINFO called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsHINFOCmd)
}

// MXCmd represents the MX command
var dnsMXCmd = &cobra.Command{
	Use:   "MX <domain name>",
	Short: "Add MX record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("MX called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsMXCmd)
}

// NAPTRCmd represents the NAPTR command
var dnsNAPTRCmd = &cobra.Command{
	Use:   "NAPTR <domain name>",
	Short: "Add NAPTR record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("NAPTR called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsNAPTRCmd)
}

// NSCmd represents the NS command
var dnsNSCmd = &cobra.Command{
	Use:   "NS <domain name>",
	Short: "Add NS record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("NS called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsNSCmd)
}

// POOLCmd represents the POOL command
var dnsPOOLCmd = &cobra.Command{
	Use:   "POOL <domain name>",
	Short: "Add POOL record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("POOL called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsPOOLCmd)
}

// SPFCmd represents the SPF command
var dnsSPFCmd = &cobra.Command{
	Use:   "SPF <domain name>",
	Short: "Add SPF record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("SPF called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSPFCmd)
}

// SRVCmd represents the SRV command
var dnsSRVCmd = &cobra.Command{
	Use:   "SRV <domain name>",
	Short: "Add SRV record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("SRV called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSRVCmd)
}

// SSHFPCmd represents the SSHFP command
var dnsSSHFPCmd = &cobra.Command{
	Use:   "SSHFP <domain name>",
	Short: "Add SSHFP record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("SSHFP called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsSSHFPCmd)
}

// TXTCmd represents the TXT command
var dnsTXTCmd = &cobra.Command{
	Use:   "TXT <domain name>",
	Short: "Add TXT record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("TXT called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsTXTCmd)
}

// URLCmd represents the URL command
var dnsURLCmd = &cobra.Command{
	Use:   "URL <domain name>",
	Short: "Add URL record type to a domain",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			cmd.Usage()
			return
		}
		fmt.Println("URL called")
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsURLCmd)
}
