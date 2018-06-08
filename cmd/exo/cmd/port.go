package cmd

import (
	"strings"
)

//go:generate stringer -type=port

type port uint16

const (
	Daytime   port = 13
	FTP       port = 21
	SSH       port = 22
	Telnet    port = 23
	SMTP      port = 25
	Time      port = 37
	Whois     port = 43
	DNS       port = 53
	TFTP      port = 69
	Gopher    port = 70
	HTTP      port = 80
	Kerberos  port = 88
	Nic       port = 101
	SFTP      port = 115
	NTP       port = 123
	IMAP      port = 143
	SNMP      port = 161
	IRC       port = 194
	HTTPS     port = 443
	RDP       port = 3389
	Minecraft port = 25565
)

func (i port) StringFormatted() string {
	res := i.String()
	if strings.HasPrefix(res, "port(") {
		return ""
	}
	return res
}
