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
	"strconv"
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

//RunstatusIncidentList is a list of incident
type RunstatusIncidentList struct {
	Incidents []RunstatusIncident `json:"incidents"`
}

// RunstatusMaintenance is a runstatus maintenance
type RunstatusMaintenance struct {
	Created     *time.Time       `json:"created,omitempty"`
	Description string           `json:"description,omitempty"`
	EndDate     *time.Time       `json:"end_date"`
	Events      []RunstatusEvent `json:"events,omitempty"`
	EventsURL   string           `json:"events_url,omitempty"`
	RealTime    bool             `json:"real_time,omitempty"`
	Services    []string         `json:"services"`
	StartDate   *time.Time       `json:"start_date"`
	Status      string           `json:"status"`
	Title       string           `json:"title"`
	URL         string           `json:"url,omitempty"`
}

//RunstatusMaintenanceList is a list of incident
type RunstatusMaintenanceList struct {
	Maintenances []RunstatusMaintenance `json:"maintenances"`
}

//RunstatusEvent is a runstatus event
type RunstatusEvent struct {
	Created *time.Time `json:"created,omitempty"`
	State   string     `json:"state,omitempty"`
	Status  string     `json:"status"`
	Text    string     `json:"text"`
}

// RunstatusService is a runstatus service
type RunstatusService struct {
	Name  string `json:"name"`
	State string `json:"state,omitempty"`
	URL   string `json:"url,omitempty"`
}

// RunstatusServiceList service list
type RunstatusServiceList struct {
	Services []RunstatusService `json:"services"`
}

// ID give the service ID
func (s *RunstatusService) ID() (int, error) {
	url := strings.TrimRight(s.URL, "/")
	urlSplited := strings.Split(url, "/")
	id, err := strconv.ParseInt(urlSplited[len(urlSplited)-1], 10, 32)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// DeleteRunstatusService delete runstatus service
func (client *Client) DeleteRunstatusService(ctx context.Context, page string, id int) error {
	// XXX Pick this URL from the API
	// GET /pages/$page | jq .services[name=$name|id=$id].url
	_, err := client.runstatusRequest(ctx, fmt.Sprintf("/pages/%s/services/%d", page, id), nil, "DELETE")
	return err
}

// CreateRunstatusService create runstatus service
func (client *Client) CreateRunstatusService(ctx context.Context, page string, service RunstatusService) error {
	// XXX: GET /pages/$page | jq .services_url
	_, err := client.runstatusRequest(ctx, "/pages/"+page+"/services", service, "POST")
	return err
}

// ListRunstatusService list runstatus service
func (client *Client) ListRunstatusService(ctx context.Context, page string) ([]RunstatusService, error) {
	// XXX: GET /pages/$page | jq .services_url
	resp, err := client.runstatusRequest(ctx, "/pages/"+page+"/services", nil, "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusServiceList
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p.Services, nil
}

// CreateRunstatusEvent create runstatus incident event
func (client *Client) CreateRunstatusEvent(ctx context.Context, page string, incidentID int, event RunstatusEvent) error {
	// XXX: GET /pages/$page | jq .incidents_url
	//      GET incidents_url + "/$id" | jq .events_url
	_, err := client.runstatusRequest(ctx, fmt.Sprintf("/pages/%s/incidents/%d/events", page, incidentID), event, "POST")
	return err
}

// ListRunstatusMaintenance list runstatus Maintenance
func (client *Client) ListRunstatusMaintenance(ctx context.Context, page string) ([]RunstatusMaintenance, error) {
	resp, err := client.runstatusRequest(ctx, "/pages/"+page+"/maintenances", nil, "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusMaintenanceList
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p.Maintenances, nil
}

// CreateRunstatusMaintenance create runstatus Maintenance
func (client *Client) CreateRunstatusMaintenance(ctx context.Context, page string, maintenance RunstatusMaintenance) error {
	// XXX: GET /pages/$page | jq .maintenances_url
	_, err := client.runstatusRequest(ctx, "/pages/"+page+"/maintenances", maintenance, "POST")
	return err
}

// DeleteRunstatusMaintenance delete runstatus Maintenance
func (client *Client) DeleteRunstatusMaintenance(ctx context.Context, page string, id int) error {
	// XXX: GET /pages/$page | jq .maintenances_url
	//      DELETE maintenances_url + "/$id"
	_, err := client.runstatusRequest(ctx, fmt.Sprintf("/pages/%s/maintenances/%d", page, id), nil, "DELETE")
	return err
}

// ID give the maintenance ID
func (m *RunstatusMaintenance) ID() (int, error) {
	url := strings.TrimRight(m.URL, "/")
	urlSplited := strings.Split(url, "/")
	id, err := strconv.ParseInt(urlSplited[len(urlSplited)-1], 10, 64)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// UpdateRunstatusMaintenance update runstatus Maintenance
func (client *Client) UpdateRunstatusMaintenance(ctx context.Context, page string, id int, event RunstatusEvent) error {
	// XXX: GET /pages/$page | jq .maintenances_url
	//      GEt maintenances_url + "/$id" | jq .events_url
	_, err := client.runstatusRequest(ctx, fmt.Sprintf("/pages/%s/maintenances/%d/events", page, id), event, "POST")
	return err
}

// ListRunstatusIncident list runstatus incident
func (client *Client) ListRunstatusIncident(ctx context.Context, page string) ([]RunstatusIncident, error) {
	resp, err := client.runstatusRequest(ctx, "/pages/"+page+"/incidents", nil, "GET")
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
	_, err := client.runstatusRequest(ctx, "/pages/"+page+"/incidents", incident, "POST")
	return err
}

// DeleteRunstatusIncident delete runstatus incident
func (client *Client) DeleteRunstatusIncident(ctx context.Context, page string, id int) error {
	// XXX: GET /pages/$page | jq .incidents_url
	//      DELETE incidents_url + "/$id"
	_, err := client.runstatusRequest(ctx, fmt.Sprintf("/pages/%s/incidents/%d", page, id), nil, "DELETE")
	return err
}

// CreateRunstatusPage create runstatus page
func (client *Client) CreateRunstatusPage(ctx context.Context, page RunstatusPage) (*RunstatusPage, error) {
	resp, err := client.runstatusRequest(ctx, "/pages", page, "POST")
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
	_, err := client.runstatusRequest(ctx, "/pages/"+pageName, nil, "DELETE")
	return err
}

// GetRunstatusPage delete runstatus page
func (client *Client) GetRunstatusPage(ctx context.Context, pageName string) (*RunstatusPage, error) {
	resp, err := client.runstatusRequest(ctx, "/pages/"+pageName, nil, "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusPage
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p, nil
}

// ListRunstatusPage delete runstatus page
func (client *Client) ListRunstatusPage(ctx context.Context) ([]RunstatusPage, error) {
	resp, err := client.runstatusRequest(ctx, "/pages", nil, "GET")
	if err != nil {
		return nil, err
	}

	var p *RunstatusPageList
	if err := json.Unmarshal(resp, &p); err != nil {
		return nil, err
	}

	return p.Results, nil
}

func (client *Client) runstatusRequest(ctx context.Context, uri string, structParam interface{}, method string) (json.RawMessage, error) {
	rawURL := client.Endpoint + uri
	reqURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	var params string
	if structParam != nil {
		m, err := json.Marshal(structParam)
		if err != nil {
			return nil, err
		}
		params = string(m)
	}

	req, err := http.NewRequest(method, reqURL.String(), strings.NewReader(params))
	if err != nil {
		return nil, err
	}

	time := time.Now().Local().Format("2006-01-02T15:04:05-0700")

	//XXX WIP for testing
	payload := fmt.Sprintf("%s%s%s", req.URL.String(), time, params)

	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	_, err = mac.Write([]byte(payload))
	if err != nil {
		return nil, err
	}
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
		// FIXME...
		e := new(DNSErrorResponse)
		if err := json.Unmarshal(b, e); err != nil {
			return nil, err
		}
		return nil, e
	}

	return b, nil
}
