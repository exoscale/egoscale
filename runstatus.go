package egoscale

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

//RunstatusPage runstatus page
type RunstatusPage struct {
	Name             string    `json:"name"`
	Created          time.Time `json:"created"`
	DarkTheme        bool      `json:"dark_theme"`
	Domain           string    `json:"domain"`
	GradientEnd      string    `json:"gradient_end"`
	GradientStart    string    `json:"gradient_start"`
	HeaderBackground string    `json:"header_background"`
	ID               int       `json:"id"`
	IncidentsURL     string    `json:"incidents_url"`
	Logo             string    `json:"logo"`
	MaintenancesURL  string    `json:"maintenances_url"`
	OkText           string    `json:"ok_text"`
	Plan             string    `json:"plan"`
	PublicURL        string    `json:"public_url"`
	ServicesURL      string    `json:"services_url"`
	State            string    `json:"state"`
	Subdomain        string    `json:"subdomain"`
	SupportEmail     string    `json:"support_email"`
	TimeZone         string    `json:"time_zone"`
	Title            string    `json:"title"`
	TitleColor       string    `json:"title_color"`
	TwitterUsername  string    `json:"twitter_username"`
	URL              string    `json:"url"`
}

// CreateRunstatusPage creates runstatus page
func (client *Client) CreateRunstatusPage(name string) (*RunstatusPage, error) {
	m, err := json.Marshal(RunstatusPage{
		Name:      name,
		Subdomain: name,
	})
	if err != nil {
		return nil, err
	}

	resp, err := client.runstatusRequest("/pages", nil, string(m), "POST")
	if err != nil {
		return nil, err
	}

	var p *RunstatusPage
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p, nil
}

func (client *Client) runstatusRequest(uri string, urlValues url.Values, params, method string) (json.RawMessage, error) {
	rawURL := client.Endpoint + uri
	url, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	q := url.Query()
	for k, vs := range urlValues {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	url.RawQuery = q.Encode()

	req, err := http.NewRequest(method, url.String(), strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	var hdr = make(http.Header)
	hdr.Add("Authorization", client.APIKey+":"+client.apiSecret)
	hdr.Add("Exoscale-Date", fmt.Sprintf("%s", time.Now()))
	hdr.Add("Accept", "application/json")
	if params != "" {
		hdr.Add("Content-Type", "application/json")
	}
	req.Header = hdr

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	contentType := resp.Header.Get("content-type")
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf(`response content-type expected to be "application/json", got %q`, contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	println(string(b))

	if resp.StatusCode >= 400 {
		e := new(DNSErrorResponse)
		if err := json.Unmarshal(b, e); err != nil {
			return nil, err
		}
		return nil, e
	}

	return b, nil
}
