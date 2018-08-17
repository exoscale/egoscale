// +build linux darwin openbsd freebsd

package cmd

import (
	"fmt"
	"log"
	"os"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

func printJSON(out string) {
	if terminal.IsTerminal(syscall.Stdout) {
		if _, err := fmt.Fprintln(os.Stdout, ""); err != nil {
			log.Fatal(err)
		}
	} else {
		if _, err := fmt.Fprintln(os.Stdout, out); err != nil {
			log.Fatal(err)
		}
	}
}
