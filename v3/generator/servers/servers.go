package servers

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
	output.WriteString(serversStaticTemplate)

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

	if os.Getenv("GENERATOR_DEBUG") == "servers" {
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

func renderServer(s *v3.Server) ([]byte, error) {
	var srv Server
	for k, v := range s.Variables {
		if k != "zone" {
			// Supporting only zone variable for Exoscale
			continue
		}
		enum := ""
		for _, z := range v.Enum {
			enum += "APIZone" + helpers.ToCamel(z) + " = \"" + z + "\"\n  "
		}

		srv = Server{
			DefaultZone:  v.Default,
			ServerURL:    strings.Replace(s.URL, "{zone}", v.Default, 1),
			RawServerURL: s.URL,
			Enum:         enum,
			APIName:      "API",
		}
	}

	t, err := template.New("Server").Parse(serverTemplate)
	if err != nil {
		return nil, err
	}

	output := bytes.NewBuffer([]byte{})
	if err := t.Execute(output, srv); err != nil {
		log.Fatal(err)
	}
	return output.Bytes(), nil
}

type Server struct {
	APIName      string
	Enum         string
	DefaultZone  string
	ServerURL    string
	RawServerURL string
}

const serverTemplate = `
type {{ .APIName }}Zone string

const (
  {{ .Enum }}
)

type Client{{.APIName}} struct {
  apiKey			string
  apiSecret			string
  serverURL     	string
  rawServerURL 		string
  httpClient    	*http.Client
  timeout			time.Duration
  pollingInterval	time.Duration
  trace				bool
  middlewares		[]RequestMiddlewareFn
}

// Client{{.APIName}}Opt represents a function setting Exoscale API client option.
type Client{{.APIName}}Opt func(*Client{{.APIName}}) error


// Client{{.APIName}}OptWithTimeout returns a Client{{.APIName}}Opt overriding the default client timeout.
func Client{{.APIName}}OptWithTimeout(v time.Duration) Client{{.APIName}}Opt {
	return func(c *Client{{.APIName}}) error {
		if v <= 0 {
			return errors.New("timeout value must be greater than 0")
		}
		c.timeout = v

		return nil
	}
}

// Client{{.APIName}}OptWithTrace returns a Client{{.APIName}}Opt enabling HTTP request/response tracing.
func Client{{.APIName}}OptWithTrace() Client{{.APIName}}Opt {
	return func(c *Client{{.APIName}}) error {
		c.trace = true
		return nil
	}
}

func Client{{.APIName}}OptWithEnvironment(env string) Client{{.APIName}}Opt {
	return func(c *Client{{.APIName}}) error {
		c.rawServerURL = strings.Replace(c.rawServerURL, "api", env, 1)
		c.serverURL = strings.Replace(c.serverURL, "api", env, 1)
		return nil
	}
}

// Client{{.APIName}}OptWithZone returns a Client{{.APIName}}Opt With a given zone.
func Client{{.APIName}}OptWithZone(zone {{ .APIName }}Zone) Client{{.APIName}}Opt {
	return func(c *Client{{.APIName}}) error {
		c.serverURL = strings.Replace(c.rawServerURL, "{zone}", string(zone), 1)
		return nil
	}
}

// Client{{.APIName}}OptWithHTTPClient returns a Client{{.APIName}}Opt overriding the default http.Client.
// Note: the Exoscale API client will chain additional middleware
// (http.RoundTripper) on the HTTP client internally, which can alter the HTTP
// requests and responses. If you don't want any other middleware than the ones
// currently set to your HTTP client, you should duplicate it and pass a copy
// instead.
func Client{{.APIName}}OptWithHTTPClient(v *http.Client) Client{{.APIName}}Opt {
	return func(c *Client{{.APIName}}) error {
		c.httpClient = v

		return nil
	}
}

func NewClient{{.APIName}}(apiKey, apiSecret string, opts ...Client{{.APIName}}Opt) (Client, error) {
  if apiKey == "" || apiSecret == "" {
	return nil, fmt.Errorf("missing or incomplete API credentials")
  }

  client := &Client{{.APIName}}{
	apiKey:				apiKey,
	apiSecret:			apiSecret,
	serverURL:			"{{.ServerURL}}",
    rawServerURL:		"{{.RawServerURL}}",
	httpClient:			http.DefaultClient,
	pollingInterval:	pollingInterval,
  }

  for _, opt := range opts {
	if err := opt(client); err != nil {
		return nil, fmt.Errorf("client configuration error: %s", err)
	}
  }

  security, err := api.NewSecurityProvider(apiKey, apiSecret)
  if err != nil {
	  return nil, fmt.Errorf("unable to initialize API security provider: %w", err)
  }

  // Tracing must be performed before API error handling in the middleware chain,
  // otherwise the response won't be dumped in case of an API error.
  if client.trace {
    client.httpClient.Transport = api.NewTraceMiddleware(client.httpClient.Transport)
  }

  client.httpClient.Transport = api.NewAPIErrorHandlerMiddleware(client.httpClient.Transport)

  return client.WithRequestMiddleware(security.Intercept), nil
}

func (c *Client{{.APIName}}) WithZone(z {{ .APIName }}Zone) Client {
  return &Client{{.APIName}}{
	serverURL:			strings.Replace(c.rawServerURL, "{zone}", string(z), 1),
    rawServerURL:  		c.rawServerURL,
	httpClient:			c.httpClient,
	middlewares:  		c.middlewares,
	pollingInterval:	c.pollingInterval,
  }
}

func (c *Client{{.APIName}}) WithContext(ctx context.Context) Client {
	return &Client{{.APIName}}{
	  serverURL:		c.serverURL,
	  rawServerURL:		c.rawServerURL,
	  httpClient:		c.httpClient,
	  middlewares:  	c.middlewares,
	  pollingInterval:	c.pollingInterval,
	}
  }

  func (c *Client{{.APIName}}) WithHttpClient(client *http.Client) Client {
	return &Client{{.APIName}}{
	  serverURL:		c.serverURL,
	  rawServerURL:		c.rawServerURL,
	  httpClient:		client,
	  middlewares:  	c.middlewares,
	  pollingInterval:	c.pollingInterval,
	}
  }

  func (c *Client{{.APIName}}) WithRequestMiddleware(f RequestMiddlewareFn) Client {
	return &Client{{.APIName}}{
		serverURL:    		c.serverURL,
		rawServerURL: 		c.rawServerURL,
		httpClient:   		c.httpClient,
		middlewares:  		append(c.middlewares, f),
		pollingInterval:	c.pollingInterval,
	}
  }

func registerRequestMiddlewares(c *Client{{.APIName}}, ctx context.Context, req *http.Request) error {
	for _, fn := range c.middlewares {
		if err := fn(ctx, req); err != nil {
			return err
		}
	}

	return nil
}
`

const serversStaticTemplate = `
type RequestMiddlewareFn func(ctx context.Context, req *http.Request) error

// UserAgent is the "User-Agent" HTTP request header added to outgoing HTTP requests.
var UserAgent = fmt.Sprintf("egoscale/%s (%s; %s/%s)",
	version.Version,
	runtime.Version(),
	runtime.GOOS,
	runtime.GOARCH)

const pollingInterval = 3 * time.Second
`
