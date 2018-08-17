// +build windows

package cmd

import (
	"fmt"
	"os"
)

// nolint: deadcode
func printJSON(out string) {
	fmt.Fprintln(os.Stdout, out)
}
