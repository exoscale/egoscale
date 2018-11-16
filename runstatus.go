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

// CreateRunstatusPage creates runstatus page
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

	println("<PAYLOAD>", payload, "<PAYLOAD>")

	// val, err := url.ParseQuery(payload)
	// if err != nil {
	// 	return nil, err
	// }

	query := strings.ToLower(payload)
	mac := hmac.New(sha256.New, []byte(client.apiSecret))
	_, err = mac.Write([]byte(query))
	if err != nil {
		return nil, err
	}

	signature := hex.EncodeToString(mac.Sum(nil))

	var hdr = make(http.Header)

	hdr.Add("Authorization", fmt.Sprintf("Exoscale-HMAC-SHA256 %s:%s", client.APIKey, signature))
	hdr.Add("Exoscale-Date", fmt.Sprintf("%s", time))
	//hdr.Add("User-Agent", fmt.Sprintf("exoscale/egoscale (%v)", Version))
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

// class ExoscaleAuth(AuthBase):
//     def __init__(self, key, secret):
//         self.key = key
//         self.secret = secret.encode('utf-8')

//     def __call__(self, request):
//         body = request.body or b''
//         if hasattr(body, 'encode'):
//             body = body.encode('utf-8')
//         date = datetime.utcnow().strftime('%Y-%m-%dT%H:%M:%SZ')
//         string_to_sign = '{0}{1}'.format(request.url,
//                                          date).encode('utf-8') + body
//         signature = hmac.new(self.secret,
//                              msg=string_to_sign,
//                              digestmod=hashlib.sha256).hexdigest()
//         auth = u'Exoscale-HMAC-SHA256 {0}:{1}'.format(self.key, signature)
//         request.headers.update({
//             'Exoscale-Date': date,
//             'Authorization': auth,
//         })
//         return request
