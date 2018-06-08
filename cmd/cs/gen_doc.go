package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/urfave/cli"
)

const frontmatter = `---
date: %s
title: %s
slug: %s
url: %s
---
`

// GenerateDocs generates markdown documentation for the commands in app
func GenerateDocs(app *cli.App, docPath string) {
	buffer := bytes.Buffer{}

	var appDescription string

	buffer.WriteString(fmt.Sprintf("# `%s`\n%s - %s <%s>\n\n", app.Name, app.Version, app.Author, app.Email))

	if app.Description != "" {
		buffer.WriteString(app.Description)
		buffer.WriteString("\n\n")
	}

	appDescription = buffer.String()

	filepath := path.Join(docPath, app.Name+".md")

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		if err := ioutil.WriteFile(filepath, []byte(appDescription), 0644); err != nil {
			log.Fatalf("doc could not be written. %s", err)
		}
	}

	buffer = bytes.Buffer{}

	if len(app.Flags) > 0 {
		buffer.WriteString("## Global Flags\n\n")
		for _, flag := range app.Flags {
			buffer.WriteString(fmt.Sprintf("- `--%s`\n", flag.GetName()))
		}
	}

	globalFlag := buffer.String()

	buffer = bytes.Buffer{}

	//buffer.WriteString("## Subcommands\n\n")

	for _, command := range app.Commands {
		buffer.WriteString(fmt.Sprintf("### `%s`\n\n", command.Name))
		if command.Usage != "" {
			buffer.WriteString(command.Usage)
			buffer.WriteString("\n\n")
		}
		if command.Description != "" {
			buffer.WriteString(command.Description)
			buffer.WriteString("\n\n")
		}
		if len(command.Flags) > 0 {
			buffer.WriteString("#### Flags\n\n")
			for _, flag := range command.Flags {
				buffer.WriteString(fmt.Sprintf("- `--%s`\n", flag.GetName()))
			}

			buffer.WriteString("\n")

			buffer.WriteString(globalFlag)

			filepath := path.Join(docPath, app.Name+"_"+command.Name+".md")

			if _, err := os.Stat(filepath); os.IsNotExist(err) {
				if err := ioutil.WriteFile(filepath, []byte(buffer.String()), 0644); err != nil {
					log.Fatalf("doc could not be written. %s", err)
				}
			}

			buffer = bytes.Buffer{}
		}
	}
}
