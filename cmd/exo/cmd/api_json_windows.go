// +build windows

package cmd

import (
	"fmt"
	"os"
)

// nolint: deadcode
func printJSON(out, theme string) {
	fmt.Fprintln(os.Stdout, out)
}
