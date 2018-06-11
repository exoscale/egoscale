package cmd

import (
	"fmt"
	"log"

	"github.com/exoscale/egoscale"

	"github.com/spf13/cobra"
)

// dnsACmd represents the A command
var dnsUpdateACmd = &cobra.Command{
	Use:   "A <domain name> <record name | id>",
	Short: "Update A record type to a domain",
	Long:  `Update an "A" record that points your domain or a subdomain to an IP address.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateACmd)
	dnsUpdateACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateACmd.Flags().StringP("address", "a", "", "Example: 127.0.0.1")
	dnsUpdateACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsACmd.MarkFlagRequired("address")
}

// AAAACmd represents the AAAA command
var dnsUpdateAAAACmd = &cobra.Command{
	Use:   "AAAA <domain name> <record name | id>",
	Short: "Update AAAA record type to a domain",
	Long:  `Update an "AAAA" record that points your domain to an IPv6 address. These records are the same as A records except they use IPv6 addresses.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateAAAACmd)
	dnsUpdateAAAACmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateAAAACmd.Flags().StringP("address", "a", "", "Example: 2001:0db8:85a3:0000:0000:EA75:1337:BEEF")
	dnsUpdateAAAACmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateAAAACmd.MarkFlagRequired("address")
}

// ALIASCmd represents the ALIAS command
var dnsUpdateALIASCmd = &cobra.Command{
	Use:   "ALIAS <domain name> <record name | id>",
	Short: "Update ALIAS record type to a domain",
	Long: `Update an "ALIAS" record. An ALIAS record is a special record that will
map a domain to another domain transparently. It can be used like a CNAME but
for a name with other records, like the root. When the record is resolved it will
look up the A records for the aliased domain and return those as the records for 
the record name. Note: If you want to redirect to a URL, use a URL record instead.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateALIASCmd)
	dnsUpdateALIASCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateALIASCmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsUpdateALIASCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateALIASCmd.MarkFlagRequired("alias")
}

// CNAMECmd represents the CNAME command
var dnsUpdateCNAMECmd = &cobra.Command{
	Use:   "CNAME <domain name> <record name | id>",
	Short: "Update CNAME record type to a domain",
	Long: `Update a "CNAME" record that aliases a subdomain to another host.
These types of records are used when a server is reached by several names. Only use CNAME records on subdomains.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateCNAMECmd)
	dnsUpdateCNAMECmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsUpdateCNAMECmd.Flags().StringP("alias", "a", "", "Alias for: Example: some-other-site.com")
	dnsUpdateCNAMECmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateCNAMECmd.MarkFlagRequired("alias")
	dnsUpdateCNAMECmd.MarkFlagRequired("name")
}

// HINFOCmd represents the HINFO command
var dnsUpdateHINFOCmd = &cobra.Command{
	Use:   "HINFO <domain name> <record name | id>",
	Short: "Update HINFO record type to a domain",
	Long:  `Update an "HINFO" record is used to describe the CPU and OS of a host.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateHINFOCmd)
	dnsUpdateHINFOCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateHINFOCmd.Flags().StringP("cpu", "c", "", "Example: IBM-PC/AT")
	dnsUpdateHINFOCmd.Flags().StringP("os", "o", "", "The operating system of the machine, example: Linux")
	dnsUpdateHINFOCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateHINFOCmd.MarkFlagRequired("cpu")
	dnsUpdateHINFOCmd.MarkFlagRequired("os")
}

// MXCmd represents the MX command
var dnsUpdateMXCmd = &cobra.Command{
	Use:   "MX <domain name> <record name | id>",
	Short: "Update MX record type to a domain",
	Long: `Update a mail exchange record that points to a mail server or relay.
These types of records are used to describe which servers handle incoming email.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateMXCmd)
	dnsUpdateMXCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>")
	dnsUpdateMXCmd.Flags().StringP("mail-server-host", "m", "", "Example: mail-server.example.com")
	dnsUpdateMXCmd.Flags().IntP("priority", "p", 0, "Common values are for example 1, 5 or 10")
	dnsUpdateMXCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateMXCmd.MarkFlagRequired("mail-server-host")
	dnsUpdateMXCmd.MarkFlagRequired("priority")
}

// NAPTRCmd represents the NAPTR command
var dnsUpdateNAPTRCmd = &cobra.Command{
	Use:   "NAPTR <domain name> <record name | id>",
	Short: "Update NAPTR record type to a domain",
	Long: `Update an "NAPTR" record to provide a means to map a resource that is not in
the domain name syntax to a label that is. More information can be found in RFC 2915.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
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
	dnsUpdateCmd.AddCommand(dnsUpdateNAPTRCmd)
	dnsUpdateNAPTRCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateNAPTRCmd.Flags().IntP("order", "o", 0, "Used to determine the processing order, lowest first.")
	dnsUpdateNAPTRCmd.Flags().IntP("preference", "", 0, "Used to give weight to records with the same value in the 'order' field, low to high.")
	dnsUpdateNAPTRCmd.Flags().StringP("service", "", "", "Service")
	dnsUpdateNAPTRCmd.Flags().StringP("regex", "", "", "The substituion expression.")
	dnsUpdateNAPTRCmd.Flags().StringP("replacement", "", "", "The next record to look up, which must be a fully-qualified domain name.")

	//flags
	dnsUpdateNAPTRCmd.Flags().BoolP("s", "", false, "Flag indicates the next lookup is for an SRV.")
	dnsUpdateNAPTRCmd.Flags().BoolP("a", "", false, "Flag indicates the next lookup is for an A or AAAA record.")
	dnsUpdateNAPTRCmd.Flags().BoolP("u", "", false, "Flag indicates the next record is the output of the regular expression as a URI.")
	dnsUpdateNAPTRCmd.Flags().BoolP("p", "", false, "Flag indicates that processing should continue in a protocol-specific fashion.")

	dnsUpdateNAPTRCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateNAPTRCmd.MarkFlagRequired("order")
	dnsUpdateNAPTRCmd.MarkFlagRequired("preference")
	dnsUpdateNAPTRCmd.MarkFlagRequired("service")
	dnsUpdateNAPTRCmd.MarkFlagRequired("regex")
	dnsUpdateNAPTRCmd.MarkFlagRequired("replacement")
}

// NSCmd represents the NS command
var dnsUpdateNSCmd = &cobra.Command{
	Use:   "NS <domain name> <record name | id>",
	Short: "Update NS record type to a domain",
	Long: `Update an "NS" record the delegates a domain to another name server.
You may only delegate subdomains (for example subdomain.yourdomain.com).`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		mailSrv, err := cmd.Flags().GetString("mail-server-host")
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "NS",
			Name:       name,
			Content:    mailSrv,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsUpdateCmd.AddCommand(dnsUpdateNSCmd)
	dnsUpdateNSCmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsUpdateNSCmd.Flags().StringP("mail-server-host", "m", "", "Example: 'ns1.example.com'")
	dnsUpdateNSCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdateNSCmd.MarkFlagRequired("name")
	dnsUpdateNSCmd.MarkFlagRequired("mail-server-host")
}

// POOLCmd represents the POOL command
var dnsUpdatePOOLCmd = &cobra.Command{
	Use:   "POOL <domain name> <record name | id>",
	Short: "Update POOL record type to a domain",
	Long: `Update a "POOL" record that aliases a subdomain to another host as
part of a pool of available CNAME records. This is a DNSimple custom record type.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "POOL",
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
	dnsUpdateCmd.AddCommand(dnsUpdatePOOLCmd)
	dnsUpdatePOOLCmd.Flags().StringP("name", "n", "", "You may use the * wildcard here.")
	dnsUpdatePOOLCmd.Flags().StringP("alias", "a", "", "Alias for: Example: 'some-other-site.com'")
	dnsUpdatePOOLCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsACmd.MarkFlagRequired("ttl")
	dnsUpdatePOOLCmd.MarkFlagRequired("name")
	dnsUpdatePOOLCmd.MarkFlagRequired("alias")
}

// // SPFCmd represents the SPF command
// var dnsSPFCmd = &cobra.Command{
// 	Use:   "SPF <domain name>",
// 	Short: "Update SPF record type to a domain",
// 	Long: `Update an "SPF" record to indicate what hosts and addresses are allowed to send mail from your domain.
// When creating an SPF record we will automatically create a corresponding TXT record
// for you as some older mail exchanges require a TXT version of the SPF record.

// 	`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		if len(args) < 2 {
// 			cmd.Usage()
// 			return
// 		}
// 		fmt.Println("SPF called")
// 	},
// }

// func init() {
// 	dnsUpdateCmd.AddCommand(dnsSPFCmd)
// 	dnsUpdateCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name> You may use the * wildcard here.")
// 	dnsUpdateCmd.Flags().BoolP("from-this-host", "e", false, "Set this flag if you send email from this host.")
// 	dnsUpdateCmd.Flags().BoolP("from-this-hosts-mx", "x", false, "Set this flag if you send email from this host's MX servers.")
// 	dnsUpdateCmd.Flags().StringP("other-servers", "o", "", "Separate addresses by spaces e.g(127.0.0.1 192.168.1.1)")
// 	dnsUpdateCmd.Flags().StringP("ip-networks", "n", "", "Separate networks by spaces e.g(127.0.0.1 192.168.1.1)")
// 	dnsUpdateCmd.Flags().StringP("ipv6-networks", "6", "", "Separate networks by spaces e.g(2001:0db8:85a3:0000:0000:EA75:1337:BEEF ...)")
// 	dnsUpdateCmd.Flags().StringP("include", "i", "", "If you send mail through your ISP's servers, and the ISP has published an SPF record, name the ISP here. Separate multiple domains with spaces.")
// 	dnsUpdateCmd.Flags().BoolP("hard-fail", "", false, "Soft fail if flag not set")
// 	dnsUpdateCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")

// 	dnsSPFCmd.MarkFlagRequired("")

// }

// SRVCmd represents the SRV command
var dnsUpdateSRVCmd = &cobra.Command{
	Use:   "SRV <domain name> <record name | id>",
	Short: "Update SRV record type to a domain",
	Long:  `Update an "SRV" record to specify the location of servers for a specific service.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}

		if name != "" {
			name = fmt.Sprintf(".%s", name)
		}

		symbName, err := cmd.Flags().GetString("symbolic-name")
		if err != nil {
			log.Fatal(err)
		}
		protocol, err := cmd.Flags().GetString("protocol")
		if err != nil {
			log.Fatal(err)
		}
		prio, err := cmd.Flags().GetInt("priority")
		if err != nil {
			log.Fatal(err)
		}
		weight, err := cmd.Flags().GetInt("weight")
		if err != nil {
			log.Fatal(err)
		}
		port, err := cmd.Flags().GetString("port")
		if err != nil {
			log.Fatal(err)
		}
		target, err := cmd.Flags().GetString("target")
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "SRV",
			Name:       fmt.Sprintf("_%s._%s%s", symbName, protocol, name),
			Content:    fmt.Sprintf("%d %s %s", weight, port, target),
			Prio:       prio,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsUpdateCmd.AddCommand(dnsUpdateSRVCmd)
	dnsUpdateSRVCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateSRVCmd.Flags().StringP("symbolic-name", "s", "", "This will be a symbolic name for the service, like 'sip'. It might also be called Service at other DNS providers.")
	dnsUpdateSRVCmd.Flags().StringP("protocol", "p", "", "This will usually be 'TCP' or 'UDP'.")
	dnsUpdateSRVCmd.Flags().IntP("priority", "", 0, "Priority")
	dnsUpdateSRVCmd.Flags().IntP("weight", "w", 0, "A relative weight for 'SRV' records with the same priority.")
	dnsUpdateSRVCmd.Flags().StringP("port", "P", "", "The 'TCP' or 'UDP' port on which the service is found.")
	dnsUpdateSRVCmd.Flags().StringP("target", "", "", "The canonical hostname of the machine providing the service.")
	dnsUpdateSRVCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsSRVCmd.MarkFlagRequired("ttl")
	dnsUpdateSRVCmd.MarkFlagRequired("symbolic-name")
	dnsUpdateSRVCmd.MarkFlagRequired("protocol")
	dnsUpdateSRVCmd.MarkFlagRequired("priority")
	dnsUpdateSRVCmd.MarkFlagRequired("weight")
	dnsUpdateSRVCmd.MarkFlagRequired("port")
	dnsUpdateSRVCmd.MarkFlagRequired("target")
}

// SSHFPCmd represents the SSHFP command
var dnsUpdateSSHFPCmd = &cobra.Command{
	Use:   "SSHFP <domain name> <record name | id>",
	Short: "Update SSHFP record type to a domain",
	Long:  `Edit an "SSHFP" record to share your SSH fingerprint with others.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		algo, err := cmd.Flags().GetInt("algorithm")
		if err != nil {
			log.Fatal(err)
		}
		fingerIDType, err := cmd.Flags().GetInt("fingerprint-type")
		if err != nil {
			log.Fatal(err)
		}
		fingerprint, err := cmd.Flags().GetString("fingerprint")
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "SSHFP",
			Name:       name,
			Content:    fmt.Sprintf("%d %d %s", algo, fingerIDType, fingerprint),
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsUpdateCmd.AddCommand(dnsUpdateSSHFPCmd)
	dnsUpdateSSHFPCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateSSHFPCmd.Flags().IntP("algorithm", "a", 0, "RSA(1) | DSA(2) | ECDSA(3) | ED25519(4)")
	dnsUpdateSSHFPCmd.Flags().IntP("fingerprint-type", "", 0, "SHA1(1) | SHA256(2)")
	dnsUpdateSSHFPCmd.Flags().StringP("fingerprint", "f", "", "Fingerprint")
	dnsUpdateSSHFPCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsSSHFPCmd.MarkFlagRequired("ttl")
	dnsUpdateSSHFPCmd.MarkFlagRequired("algorithm")
	dnsUpdateSSHFPCmd.MarkFlagRequired("fingerprint-type")
}

// TXTCmd represents the TXT command
var dnsUpdateTXTCmd = &cobra.Command{
	Use:   "TXT <domain name> <record name | id>",
	Short: "Update TXT record type to a domain",
	Long: `Update a "TXT" record. This is useful for domain records that are not covered by
the standard record types. For example, Google uses this type of record for domain verification.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		content, err := cmd.Flags().GetString("content")
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "TXT",
			Name:       name,
			Content:    content,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsUpdateCmd.AddCommand(dnsUpdateTXTCmd)
	dnsUpdateTXTCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateTXTCmd.Flags().StringP("content", "c", "", "Content record")
	dnsUpdateTXTCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsTXTCmd.MarkFlagRequired("ttl")
	dnsUpdateTXTCmd.MarkFlagRequired("content")
}

// URLCmd represents the URL command
var dnsUpdateURLCmd = &cobra.Command{
	Use:   "URL <domain name> <record name | id>",
	Short: "Update URL record type to a domain",
	Long: `Update an URL redirection record that points your domain to a URL.
This type of record uses an HTTP redirect to redirect visitors from a domain to a web site.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Usage()
			return
		}

		recordID, err := getRecordIDByName(csDNS, args[0], args[1])
		if err != nil {
			log.Fatal(err)
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatal(err)
		}
		destURL, err := cmd.Flags().GetString("destination-url")
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

		resp, err := csDNS.UpdateRecord(args[0], egoscale.DNSRecordUpdate{
			DomainID:   domain.ID,
			ID:         recordID,
			TTL:        ttl,
			RecordType: "URL",
			Name:       name,
			Content:    destURL,
		})
		if err != nil {
			log.Fatal(err)
		}
		println(resp.ID)
	},
}

func init() {
	dnsUpdateCmd.AddCommand(dnsUpdateURLCmd)
	dnsUpdateURLCmd.Flags().StringP("name", "n", "", "Leave this blank to create a record for <domain name>, You may use the '*' wildcard here.")
	dnsUpdateURLCmd.Flags().StringP("destination-url", "d", "", "Example: https://www.example.com")
	dnsUpdateURLCmd.Flags().IntP("ttl", "t", 3600, "The time in second to leave (refresh rate) of the record.")
	//dnsURLCmd.MarkFlagRequired("ttl")
	dnsUpdateURLCmd.MarkFlagRequired("destination-url")
}
