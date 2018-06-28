package main

import (
	"github.com/exoscale/egoscale/cmd/exo/cmd"
)

func main() {
	cmd.RootCmd.GenBashCompletionFile("bash_completion")
	cmd.RootCmd.GenZshCompletionFile("zsh_completion")
}
