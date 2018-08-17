// +build windows

package main

import (
	"fmt"
	"os"
)

// nolint: deadcode
func printJSON(out, theme string) {
	fmt.Fprintln(os.Stdout, out)
}
