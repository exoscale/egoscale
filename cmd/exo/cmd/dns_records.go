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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		alias, err := cmd.Flags().GetString("alias")
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
			RecordType: "ALIAS",
			Name:       name,
			Content:    alias,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsALIASCmd)
	dnsALIASCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsALIASCmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsALIASCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsALIASCmd.MarkFlagRequired("alias")
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		alias, err := cmd.Flags().GetString("alias")
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
			RecordType: "CNAME",
			Name:       name,
			Content:    alias,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsCNAMECmd)
	dnsCNAMECmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsCNAMECmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsCNAMECmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsCNAMECmd.MarkFlagRequired("alias")
	dnsCNAMECmd.MarkFlagRequired("name")
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		cpu, err := cmd.Flags().GetString("cpu")
		if err != nil {
			log.Fatal(err)
		}
		os, err := cmd.Flags().GetString("os")
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
			RecordType: "HINFO",
			Name:       name,
			Content:    fmt.Sprintf("%s %s", cpu, os),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsHINFOCmd)
	dnsHINFOCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsHINFOCmd.Flags().StringP("cpu", "c", "", "Example: IBM-PC/AT")
	dnsHINFOCmd.Flags().StringP("os", "o", "", "The operating system of the machine, example: Linux")
	dnsHINFOCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsHINFOCmd.MarkFlagRequired("cpu")
	dnsHINFOCmd.MarkFlagRequired("os")
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		mailSrv, err := cmd.Flags().GetString("mail-server-host")
		if err != nil {
			log.Fatal(err)
		}
		priority, err := cmd.Flags().GetInt("priority")
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
			RecordType: "MX",
			Name:       name,
			Content:    mailSrv,
			Prio:       priority,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsMXCmd)
	dnsMXCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>")
	dnsMXCmd.Flags().StringP("mail-server-host", "m", "", "Example: mail-server.example.com")
	dnsMXCmd.Flags().IntP("priority", "p", 0, "Common values are for example 1, 5 or 10")
	dnsMXCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsMXCmd.MarkFlagRequired("mail-server-host")
	dnsMXCmd.MarkFlagRequired("priority")
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
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		order, err := cmd.Flags().GetInt("order")
		if err != nil {
			log.Fatal(err)
		}
		preference, err := cmd.Flags().GetInt("preference")
		if err != nil {
			log.Fatal(err)
		}

		flags := ""
		//flags
		s, err := cmd.Flags().GetBool("s")
		if err != nil {
			log.Fatal(err)
		}
		a, err := cmd.Flags().GetBool("a")
		if err != nil {
			log.Fatal(err)
		}
		u, err := cmd.Flags().GetBool("u")
		if err != nil {
			log.Fatal(err)
		}
		p, err := cmd.Flags().GetBool("p")
		if err != nil {
			log.Fatal(err)
		}

		if s {
			flags += "s"
		}
		if a {
			flags += "a"
		}
		if u {
			flags += "u"
		}
		if p {
			flags += "p"
		}

		service, err := cmd.Flags().GetString("service")
		if err != nil {
			log.Fatal(err)
		}
		regex, err := cmd.Flags().GetString("regex")
		if err != nil {
			log.Fatal(err)
		}
		replacement, err := cmd.Flags().GetString("replacement")
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
			RecordType: "NAPTR",
			Name:       name,
			Content:    fmt.Sprintf("%d %d %q %q %q %q", order, preference, flags, service, regex, replacement),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsAddCmd.AddCommand(dnsNAPTRCmd)
	dnsNAPTRCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsNAPTRCmd.Flags().IntP("order", "o", 0, "Used to determine the processing order, lowest first.")
	dnsNAPTRCmd.Flags().IntP("preference", "", 0, "Used to give weight to records with the same value in the 'order' field, low to high.")
	dnsNAPTRCmd.Flags().StringP("service", "", "", "Service")
	dnsNAPTRCmd.Flags().StringP("regex", "", "", "The substituion expression.")
	dnsNAPTRCmd.Flags().StringP("replacement", "", "", "The next record to look up, which must be a fully-qualified domain name.")

	//flags
	dnsNAPTRCmd.Flags().BoolP("s", "", false, "Flag indicates the next lookup is for an SRV.")
	dnsNAPTRCmd.Flags().BoolP("a", "", false, "Flag indicates the next lookup is for an A or AAAA record.")
	dnsNAPTRCmd.Flags().BoolP("u", "", false, "Flag indicates the next record is the output of the regular expression as a URI.")
	dnsNAPTRCmd.Flags().BoolP("p", "", false, "Flag indicates that processing should continue in a protocol-specific fashion.")

	dnsNAPTRCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsNAPTRCmd.MarkFlagRequired("order")
	dnsNAPTRCmd.MarkFlagRequired("preference")
	dnsNAPTRCmd.MarkFlagRequired("service")
	dnsNAPTRCmd.MarkFlagRequired("regex")
	dnsNAPTRCmd.MarkFlagRequired("replacement")
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
