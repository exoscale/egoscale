package egoscale

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
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
	Name             string     `json:"name"`
	Created          *time.Time `json:"created,omitempty"`
	DarkTheme        bool       `json:"dark_theme,omitempty"`
	Domain           string     `json:"domain,omitempty"`
	GradientEnd      string     `json:"gradient_end,omitempty"`
	GradientStart    string     `json:"gradient_start,omitempty"`
	HeaderBackground string     `json:"header_background,omitempty"`
	ID               int        `json:"id,omitempty"`
	IncidentsURL     string     `json:"incidents_url,omitempty"`
	Logo             string     `json:"logo,omitempty"`
	MaintenancesURL  string     `json:"maintenances_url,omitempty"`
	OkText           string     `json:"ok_text,omitempty"`
	Plan             string     `json:"plan,omitempty"`
	PublicURL        string     `json:"public_url,omitempty"`
	ServicesURL      string     `json:"services_url,omitempty"`
	State            string     `json:"state,omitempty"`
	Subdomain        string     `json:"subdomain"`
	SupportEmail     string     `json:"support_email,omitempty"`
	TimeZone         string     `json:"time_zone,omitempty"`
	Title            string     `json:"title,omitempty"`
	TitleColor       string     `json:"title_color,omitempty"`
	TwitterUsername  string     `json:"twitter_username,omitempty"`
	URL              string     `json:"url,omitempty"`
}

//RunstatusPageList runstatus page list
type RunstatusPageList struct {
	Count    int             `json:"count"`
	Next     string          `json:"next"`
	Previous string          `json:"previous"`
	Results  []RunstatusPage `json:"results"`
}

//RunstatusIncident is a runstatus incident
type RunstatusIncident struct {
	EndDate    *time.Time       `json:"end_date,omitempty"`
	Events     []RunstatusEvent `json:"events,omitempty"`
	EventsURL  string           `json:"events_url,omitempty"`
	ID         int              `json:"id,omitempty"`
	PostMortem string           `json:"post_mortem,omitempty"`
	RealTime   bool             `json:"real_time,omitempty"`
	Services   []string         `json:"services"`
	StartDate  *time.Time       `json:"start_date,omitempty"`
	State      string           `json:"state"`
	Status     string           `json:"status"`
	StatusText string           `json:"status_text"`
	Title      string           `json:"title"`
	URL        string           `json:"url,omitempty"`
}

//RunstatusEvent is a runstatus event
type RunstatusEvent struct {
	Created *time.Time `json:"created"`
	State   string     `json:"state"`
	Status  string     `json:"status"`
	Text    string     `json:"text"`
}

//RunstatusIncidentList is a list of incident
type RunstatusIncidentList struct {
	Incidents []RunstatusIncident `json:"incidents"`
}

// ListRunstatusIncident list runstatus incident
func (client *Client) ListRunstatusIncident(ctx context.Context, page string) ([]RunstatusIncident, error) {
	resp, err := client.runstatusRequest(ctx, "/pages/"+page+"/incidents", nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusIncidentList
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p.Incidents, nil
}

// CreateRunstatusIncident create runstatus incident
func (client *Client) CreateRunstatusIncident(ctx context.Context, page string, incident RunstatusIncident) error {
	m, err := json.Marshal(incident)
	if err != nil {
		return err
	}

	_, err = client.runstatusRequest(ctx, "/pages/"+page+"/incidents", nil, string(m), "POST")
	return err
}

// DeleteRunstatusIncident delete runstatus incident
func (client *Client) DeleteRunstatusIncident(ctx context.Context, page, id string) error {
	_, err := client.runstatusRequest(ctx, "/pages/"+page+"/incidents/"+id, nil, "", "DELETE")
	return err
}

// CreateRunstatusPage create runstatus page
func (client *Client) CreateRunstatusPage(ctx context.Context, page RunstatusPage) (*RunstatusPage, error) {
	m, err := json.Marshal(page)
	if err != nil {
		return nil, err
	}

	resp, err := client.runstatusRequest(ctx, "/pages", nil, string(m), "POST")
	if err != nil {
		return nil, err
	}

	var p *RunstatusPage
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p, nil
}

// DeleteRunstatusPage delete runstatus page
func (client *Client) DeleteRunstatusPage(ctx context.Context, pageName string) error {
	_, err := client.runstatusRequest(ctx, "/pages/"+pageName, nil, "", "DELETE")
	return err
}

// ListRunstatusPage delete runstatus page
func (client *Client) ListRunstatusPage(ctx context.Context) ([]RunstatusPage, error) {
	resp, err := client.runstatusRequest(ctx, "/pages", nil, "", "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusPageList
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p.Results, nil
}

func (client *Client) runstatusRequest(ctx context.Context, uri string, urlValues url.Values, params, method string) (json.RawMessage, error) {
	rawURL := client.Endpoint + uri
	reqURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	q := reqURL.Query()
	for k, vs := range urlValues {
		for _, v := range vs {
			q.Add(k, v)
		}
	}

	reqURL.RawQuery = q.Encode()

	req, err := http.NewRequest(method, reqURL.String(), strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	time := time.Now().Local().Format("2006-01-02T15:04:05-0700")

	//XXX WIP for testing
	payload := fmt.Sprintf("%s%s%s", req.URL.String(), time, params)

	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	mac.Write([]byte(payload))
	signature := hex.EncodeToString(mac.Sum(nil))

	var hdr = make(http.Header)

	hdr.Add("Authorization", fmt.Sprintf("Exoscale-HMAC-SHA256 %s:%s", client.APIKey, signature))
	hdr.Add("Exoscale-Date", time)
	hdr.Add("User-Agent", fmt.Sprintf("exoscale/egoscale (%v)", Version))
	hdr.Add("Accept", "application/json")
	if params != "" {
		hdr.Add("Content-Type", "application/json")
	}
	req.Header = hdr

	req = req.WithContext(ctx)

	resp, err := client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	contentType := resp.Header.Get("content-type")
	if resp.StatusCode < 400 && contentType == "" {
		return nil, nil
	}
	if !strings.Contains(contentType, "application/json") {
		return nil, fmt.Errorf(`response content-type expected to be "application/json", got %q`, contentType)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		e := new(DNSErrorResponse)
		if err := json.Unmarshal(b, e); err != nil {
			return nil, err
		}
		return nil, e
	}

	return b, nil
}
