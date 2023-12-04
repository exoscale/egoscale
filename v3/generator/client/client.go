package client

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"net/url"
	"os"
	"strings"
	"text/template"

	"github.com/exoscale/egoscale/v3/generator/helpers"
	"github.com/pb33f/libopenapi"
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"
)

func Generate(doc libopenapi.Document, path, packageName string) error {
	r, errs := doc.BuildV3Model()
	for _, err := range errs {
		if err != nil {
			return fmt.Errorf("errors %v", errs)
		}
	}

	output := bytes.NewBuffer(helpers.Header(packageName, "v0.0.1"))
	output.WriteString(fmt.Sprintf(`package %s
	import (
		"fmt"
		"net/http"
		"context"
		"strings"
		"runtime"
		"time"
		"errors"

		"github.com/exoscale/egoscale/version"
		api "github.com/exoscale/egoscale/v3/api"
	)
	`, packageName))

	for _, s := range r.Model.Servers {
		api, err := extractAPIName(s)
		if err != nil {
			return err
		}

		if api != "api" {
			// Skip generating code for preprod "ppapi" server.
			continue
		}

		srv, err := renderServer(s)
		if err != nil {
			return err
		}
		output.Write(srv)
	}

	if os.Getenv("GENERATOR_DEBUG") == "client" {
		fmt.Println(output.String())
	}

	content, err := format.Source(output.Bytes())
	if err != nil {
		return err
	}

	return os.WriteFile(path, content, os.ModePerm)
}

func extractAPIName(s *v3.Server) (string, error) {
	var URL string
	for k, v := range s.Variables {
		URL = strings.Replace(s.URL, fmt.Sprintf("{%s}", k), v.Default, 1)
	}

	u, err := url.Parse(URL)
	if err != nil {
		return "", err
	}

	h := strings.Split(u.Host, "-")
	if len(h) < 2 {
		return "", fmt.Errorf("malformed server host: %s", u.Host)
	}

	return h[0], nil
}

type ServerTmpl struct {
	Enum         string
	DefaultZone  string
	ServerURL    string
	RawServerURL string
}

func renderServer(s *v3.Server) ([]byte, error) {
	var srv ServerTmpl
	for k, v := range s.Variables {
		if k != "zone" {
			// Supporting only zone variable for Exoscale
			continue
		}
		enum := ""
		for _, z := range v.Enum {
			enum += "ClientZone" + helpers.ToCamel(z) + " ClientZone = \"" + z + "\"\n  "
		}

		srv = ServerTmpl{
			DefaultZone:  v.Default,
			ServerURL:    strings.Replace(s.URL, "{zone}", v.Default, 1),
			RawServerURL: s.URL,
			Enum:         enum,
		}
	}

	t, err := template.New("client.tmpl").ParseFiles("./client/client.tmpl")
	if err != nil {
		return nil, err
	}

	output := bytes.NewBuffer([]byte{})
	if err := t.Execute(output, srv); err != nil {
		log.Fatal(err)
	}
	return output.Bytes(), nil
}