package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const frontmatter = `---
date: %s
title: %s
slug: %s
url: %s
---
`

func writeFlag(buffer *bytes.Buffer, flag cli.Flag) {
	doc := strings.SplitN(flag.String(), "\t", 2)
	if len(doc) != 2 {
		doc = []string{flag.GetName(), ""}
	} else {
		doc[1] = " - " + doc[1]
	}
	buffer.WriteString(fmt.Sprintf("-- `%s`%s\n", doc[0], doc[1]))
}

// generateDocs generates markdown documentation for the commands in app
func generateDocs(app *cli.App, docPath string) {
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
			writeFlag(&buffer, flag)
		}
	}

	globalFlag := buffer.String()

	now := time.Now().Format(time.RFC3339)

	for _, command := range app.Commands {
		base := command.Name
		slug := strings.ToLower(base)
		url := fmt.Sprintf("/cs/%s/", slug)

		buffer = bytes.Buffer{}
		buffer.WriteString(fmt.Sprintf(frontmatter, now, base, slug, url))
		buffer.WriteString(fmt.Sprintf("# `%s`\n\n", command.Name))
		if command.Usage != "" {
			buffer.WriteString(command.Usage)
			buffer.WriteString("\n\n")
		}
		if command.Description != "" {
			buffer.WriteString(command.Description)
			buffer.WriteString("\n\n")
		}
		if len(command.Flags) > 0 {
			buffer.WriteString("## Flags\n\n")
			for _, flag := range command.Flags {
				writeFlag(&buffer, flag)
			}

			buffer.WriteString("\n")

			buffer.WriteString(globalFlag)

			filepath := path.Join(docPath, slug+".md")

			if err := ioutil.WriteFile(filepath, []byte(buffer.String()), 0644); err != nil {
				log.Fatalf("doc could not be written. %s", err)
			}
		}
	}
}
