package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"syscall"
	"text/tabwriter"

	"github.com/exoscale/egoscale"
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

func printResponseHelp(out io.Writer, response interface{}) {
	value := reflect.ValueOf(response)
	typeof := reflect.TypeOf(response)

	w := tabwriter.NewWriter(out, 0, 0, 1, ' ', tabwriter.FilterHTML)
	if _, err := fmt.Fprintln(w, "FIELD\tTYPE\tDOCUMENTATION"); err != nil {
		log.Fatal(err)
	}

	for typeof.Kind() == reflect.Ptr {
		typeof = typeof.Elem()
		value = value.Elem()
	}

	for i := 0; i < typeof.NumField(); i++ {
		field := typeof.Field(i)
		tag := field.Tag
		doc := "-"
		if d, ok := tag.Lookup("doc"); ok {
			doc = d
		}

		name := field.Type.Name()
		if name == "" {
			if field.Type.Kind() == reflect.Slice {
				name = "[]" + field.Type.Elem().Name()
			}
		}

		if json, ok := tag.Lookup("json"); ok {
			n, _ := egoscale.ExtractJSONTag(field.Name, json)
			if _, err := fmt.Fprintf(w, "%s\t%s\t%s\n", n, name, doc); err != nil {
				log.Fatal(err)
			}
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
}
