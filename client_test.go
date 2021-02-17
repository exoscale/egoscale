package egoscale

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestClientAPIName(t *testing.T) {
	cs := NewClient("ENDPOINT", "KEY", "SECRET")
	req := &ListAPIs{}
	if cs.APIName(req) != "listApis" {
		t.Errorf("APIName is wrong, wanted listApis")
	}
	if cs.APIName(&AuthorizeSecurityGroupIngress{}) != "authorizeSecurityGroupIngress" {
		t.Errorf("APIName is wrong, wanted Ingress")
	}
	if cs.APIName(&AuthorizeSecurityGroupEgress{}) != "authorizeSecurityGroupEgress" {
		t.Errorf("APIName is wrong, wanted Egress")
	}
	if cs.APIName(AuthorizeSecurityGroupEgress{}) != "authorizeSecurityGroupEgress" {
		t.Errorf("APIName is wrong, wanted Egress")
	}
}

func TestClientAPIDescription(t *testing.T) {
	cs := NewClient("ENDPOINT", "KEY", "SECRET")
	req := &ListAPIs{}
	desc := cs.APIDescription(req)
	if desc != "lists all available apis on the server" {
		t.Errorf("APIDescription of listApis is wrong, got %q", desc)
	}
}

func TestClientResponse(t *testing.T) {
	cs := NewClient("ENDPOINT", "KEY", "SECRET")

	r := cs.Response(&ListAPIs{})
	switch r.(type) {
	case *ListAPIsResponse:
		// do nothing
	default:
		t.Errorf("request is wrong, got %t", r)
	}

	ar := cs.Response(&DeployVirtualMachine{})
	switch ar.(type) {
	case *VirtualMachine:
		// do nothing
	default:
		t.Errorf("asyncRequest is wrong, got %t", ar)
	}
}

func TestClientSyncDelete(t *testing.T) {
	bodySuccessString := `
{"delete%sresponse": {
	"success": "true"
}}`
	bodySuccessBool := `
{"delete%sresponse": {
	"success": true
}}`

	bodyError := `
{"delete%sresponse": {
	"errorcode": 431,
	"cserrorcode": 9999,
	"errortext": "This is a dummy error",
	"uuidList": []
}}`

	things := []struct {
		name      string
		deletable Deletable
	}{
		{"securitygroup", &SecurityGroup{ID: MustParseUUID("09ae3132-3a35-458c-9607-e3c77dd0465b")}},
		{"securitygroup", &SecurityGroup{Name: "test"}},
		{"sshkeypair", &SSHKeyPair{Name: "test"}},
	}

	for _, thing := range things {
		ts := newServer(
			response{200, jsonContentType, fmt.Sprintf(bodySuccessString, thing.name)},
			response{200, jsonContentType, fmt.Sprintf(bodySuccessBool, thing.name)},
			response{431, jsonContentType, fmt.Sprintf(bodyError, thing.name)},
		)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		for i := 0; i < 2; i++ {
			if err := cs.Delete(thing.deletable); err != nil {
				t.Errorf("Deletion of %#v. Err: %s", thing.deletable, err)
			}
		}

		if err := cs.Delete(thing.deletable); err == nil {
			t.Errorf("Deletion of %v an error was expected", thing.deletable)
		}

		ts.Close()
	}
}

func TestClientAsyncDelete(t *testing.T) {
	body := `
{"%sresponse": {
	"jobid": "01ed7adc-8b81-4e33-a0f2-4f55a3b880cd",
	"jobresult": {
		"success": true
	},
	"jobstatus": 1
}}`
	bodyError := `
{"%sresponse": {
	"jobid": "01ed7adc-8b81-4e33-a0f2-4f55a3b880cd",
	"jobresult": {
		"success": false,
		"displaytext": "herp derp",
	},
	"jobstatus": 2
}}`

	id := MustParseUUID("96816f59-9986-499c-91c5-f47bd1122c4b")
	things := []struct {
		name      string
		deletable Deletable
	}{
		{"deleteaffinitygroup", &AffinityGroup{ID: id}},
		{"deleteaffinitygroup", &AffinityGroup{Name: "affinity group name"}},
		{"disassociateipaddress", &IPAddress{ID: id}},
		{"destroyvirtualmachine", &VirtualMachine{ID: id}},
	}

	for _, thing := range things {
		ts := newServer(
			response{200, jsonContentType, fmt.Sprintf(body, thing.name)},
			response{400, jsonContentType, fmt.Sprintf(bodyError, thing.name)},
		)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		if err := cs.Delete(thing.deletable); err != nil {
			t.Errorf("Deletion of %#v. Err: %s", thing, err)
		}
		if err := cs.Delete(thing.deletable); err == nil {
			t.Errorf("Deletion of %#v. An error was expected", thing)
		}

		ts.Close()
	}
}

func TestClientDeleteFailure(t *testing.T) {
	things := []Deletable{
		&AffinityGroup{},
		&SecurityGroup{},
		&SSHKeyPair{},
		&VirtualMachine{},
		&IPAddress{},
	}

	for _, thing := range things {
		ts := newServer()

		cs := NewClient(ts.URL, "KEY", "SECRET")

		if err := cs.Delete(thing); err == nil {
			t.Errorf("Deletion of %#v. Should have failed", thing)
		}

		ts.Close()
	}
}

func TestClientGetFailure(t *testing.T) {
	things := []Listable{
		nil,
		&Account{},
		(*ListAccounts)(nil),
		&AffinityGroup{},
		(*ListAffinityGroups)(nil),
		&IPAddress{},
		(*ListPublicIPAddresses)(nil),
		&Nic{},
		(*ListNics)(nil),
		&SSHKeyPair{},
		(*ListSSHKeyPairs)(nil),
		&SecurityGroup{},
		(*ListSecurityGroups)(nil),
		&ServiceOffering{},
		(*ListServiceOfferings)(nil),
		&Snapshot{},
		(*ListSnapshots)(nil),
		&VirtualMachine{},
		(*ListVirtualMachines)(nil),
		&Volume{},
		(*ListVolumes)(nil),
		&Zone{},
		(*ListZones)(nil),
		&ISO{},
		(*ListISOs)(nil),
	}

	for _, thing := range things {
		ts := newServer()

		cs := NewClient(ts.URL, "KEY", "SECRET")

		if _, err := cs.Get(thing); err == nil {
			t.Errorf("Get of %#v. Should have failed", thing)
		}

		ts.Close()
	}
}

func TestClientGetNone(t *testing.T) {
	body := `{"list%sresponse": {}}`
	bodyError := `{"errorresponse": {
		"cserrorcode": 9999,
		"errorcode": 431,
		"errortext": "Unable to execute API command due to invalid value.",
		"uuidList": []
	}}`
	id := MustParseUUID("4557261a-c4b9-45a3-91b3-e48ef55857ed")
	things := []struct {
		name     string
		listable Listable
	}{
		{"zones", &Zone{ID: id}},
		{"zones", &Zone{Name: "test zone"}},
		{"publicipaddresses", &IPAddress{ID: id, IsElastic: true, ForVirtualNetwork: true}},
		{"publicipaddresses", &IPAddress{IPAddress: net.ParseIP("127.0.0.1"), IsSourceNat: true}},
		{"sshkeypairs", &SSHKeyPair{Name: "1"}},
		{"sshkeypairs", &SSHKeyPair{Fingerprint: "test ssh keypair"}},
		{"affinitygroups", &AffinityGroup{ID: id}},
		{"affinitygroups", &AffinityGroup{Name: "test affinity group"}},
		{"securitygroups", &SecurityGroup{ID: id}},
		{"securitygroups", &SecurityGroup{Name: "test affinity group"}},
		{"virtualmachines", &VirtualMachine{ID: id}},
		{"volumes", &Volume{ID: id}},
		{"templates", &Template{ID: id, IsFeatured: true}},
		{"serviceofferings", &ServiceOffering{ID: id}},
		{"accounts", &Account{}},
		{"nics", &Nic{}},
		{"snapshots", &Snapshot{}},
		{"isos", &ISO{}},
	}

	for _, thing := range things {
		ts := newServer(
			response{200, jsonContentType, fmt.Sprintf(body, thing.name)},
			response{431, jsonContentType, bodyError},
		)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		_, err := cs.Get(thing.listable)
		if err == nil {
			t.Error("an error was expected")
			continue
		}

		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}

		ts.Close()
	}
}

func TestClientGetZero(t *testing.T) {
	body := `
	{"list%sresponse": {
		"%s": []
	}}`

	id := MustParseUUID("4557261a-c4b9-45a3-91b3-e48ef55857ed")
	things := []struct {
		name     string
		listable Listable
	}{
		{"zone", &Zone{ID: id}},
		{"zone", &Zone{Name: "test zone"}},
		{"publicipaddress", &IPAddress{ID: id}},
		{"publicipaddress", &IPAddress{IPAddress: net.ParseIP("127.0.0.1")}},
		{"sshkeypair", &SSHKeyPair{Name: "1"}},
		{"sshkeypair", &SSHKeyPair{Fingerprint: "test ssh keypair"}},
		{"affinitygroup", &AffinityGroup{ID: id}},
		{"affinitygroup", &AffinityGroup{Name: "test affinity group"}},
		{"securitygroup", &SecurityGroup{ID: id}},
		{"securitygroup", &SecurityGroup{Name: "test affinity group"}},
		{"virtualmachine", &VirtualMachine{ID: id}},
		{"volume", &Volume{ID: id}},
		{"template", &Template{ID: id, IsFeatured: true}},
		{"serviceoffering", &ServiceOffering{ID: id}},
		{"account", &Account{}},
		{"nic", &Nic{ID: id}},
		{"snapshot", &Snapshot{ID: id}},
		{"iso", &ISO{ID: id}},
	}

	for _, thing := range things {
		plural := thing.name
		if strings.HasSuffix(plural, "s") {
			plural += "es"
		} else {
			plural += "s"
		}
		resp := response{200, jsonContentType, fmt.Sprintf(body, plural, thing.name)}
		ts := newServer(resp)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		// fake 431
		_, err := cs.Get(thing.listable)
		if err == nil {
			t.Errorf("an error was expected")
		}

		if err != ErrNotFound {
			t.Fatalf("expected ErrNotFound, got %v", err)
		}

		ts.Close()
	}
}

func TestClientGetTooMany(t *testing.T) {
	body := `
	{"list%sresponse": {
		"count": 2,
		"affinitygroup": [{}, {}],
		"publicipaddress": [{}, {}],
		"securitygroup": [{}, {}],
		"sshkeypair": [{}, {}],
		"virtualmachine": [{}, {}],
		"volume": [{}, {}],
		"zone": [{}, {}],
		"template": [{}, {}],
		"serviceoffering": [{}, {}],
		"account": [{}, {}],
		"nic": [{}, {}],
		"snapshot": [{}, {}],
		"iso": [{}, {}]
	}}`

	id := MustParseUUID("4557261a-c4b9-45a3-91b3-e48ef55857ed")
	things := []struct {
		name     string
		listable Listable
	}{
		{"zones", &Zone{ID: id}},
		{"zones", &Zone{Name: "test zone"}},
		{"publicipaddresses", &IPAddress{ID: id}},
		{"publicipaddresses", &IPAddress{IPAddress: net.ParseIP("127.0.0.1")}},
		{"sshkeypairs", &SSHKeyPair{Name: "1"}},
		{"sshkeypairs", &SSHKeyPair{Fingerprint: "test ssh keypair"}},
		{"affinitygroups", &AffinityGroup{ID: id}},
		{"affinitygroups", &AffinityGroup{Name: "test affinity group"}},
		{"securitygroups", &SecurityGroup{ID: id}},
		{"securitygroups", &SecurityGroup{Name: "test affinity group"}},
		{"virtualmachines", &VirtualMachine{ID: id}},
		{"volumes", &Volume{ID: id}},
		{"templates", &Template{ID: id, IsFeatured: true}},
		{"serviceofferings", &ServiceOffering{ID: id}},
		{"accounts", &Account{}},
		{"nics", &Nic{}},
		{"snapshots", &Snapshot{}},
		{"isos", &ISO{}},
	}

	for _, thing := range things {
		resp := response{200, jsonContentType, fmt.Sprintf(body, thing.name)}
		ts := newServer(resp)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		// Too many
		_, err := cs.Get(thing.listable)
		if err == nil {
			t.Errorf("an error was expected")
		}

		if err != ErrTooManyFound {
			t.Fatalf("expected ErrTooManyFound, got %v", err)
		}

		ts.Close()
	}
}

func TestClientTrace(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `{"listzonesresponse":{ "count": 0, "zone": []}}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	// XXX test something... this only increases the coverage
	cs.TraceOn()

	_, err := cs.Request(&ListZones{})

	cs.TraceOff()

	if err != nil {
		t.Error(err)
	}
}

// Things that can be listed, paginated

type lsTest struct {
	name      string
	fieldName string
	listables []Listable
}

func lsTests() []lsTest {
	ts := []lsTest{
		{"zones", "", []Listable{
			&Zone{},
			&ListZones{},
		}},
		{"publicipaddresses", "", []Listable{
			&IPAddress{},
			&ListPublicIPAddresses{},
		}},
		{"sshkeypairs", "", []Listable{
			&SSHKeyPair{},
			&ListSSHKeyPairs{},
		}},
		{"affinitygroups", "", []Listable{
			&AffinityGroup{},
			&ListAffinityGroups{},
		}},
		{"securitygroups", "", []Listable{
			&SecurityGroup{},
			&ListSecurityGroups{},
		}},
		{"virtualmachines", "", []Listable{
			&VirtualMachine{},
			&ListVirtualMachines{},
		}},
		{"volumes", "", []Listable{
			&Volume{},
			&ListVolumes{},
		}},
		{"templates", "", []Listable{
			&Template{IsFeatured: true},
			&ListTemplates{TemplateFilter: "featured"},
		}},
		{"serviceofferings", "", []Listable{
			&ServiceOffering{},
			&ListServiceOfferings{},
		}},
		{"networks", "", []Listable{
			&Network{},
			&ListNetworks{},
		}},
		{"accounts", "", []Listable{
			&Account{},
			&ListAccounts{},
		}},
		{"nics", "", []Listable{
			&Nic{},
			&ListNics{},
		}},
		{"snapshots", "", []Listable{
			&Snapshot{},
			&ListSnapshots{},
		}},
		{"events", "", []Listable{
			&Event{},
			&ListEvents{},
		}},
		{"eventtypes", "", []Listable{
			&EventType{},
			&ListEventTypes{},
		}},
		{"resourcelimits", "", []Listable{
			&ResourceLimit{},
			&ListResourceLimits{},
		}},
		{"resourcedetails", "", []Listable{
			&ResourceDetail{
				ResourceType: "UserVM",
			},
			&ListResourceDetails{
				ResourceType: "UserVM",
			},
		}},
		{"tags", "", []Listable{
			&ResourceTag{},
			&ListTags{},
		}},
		{"users", "", []Listable{
			&User{},
			&ListUsers{},
		}},
		{"instancegroups", "", []Listable{
			&InstanceGroup{},
			&ListInstanceGroups{},
		}},
		{"asyncjobs", "", []Listable{
			&AsyncJobResult{},
			&ListAsyncJobs{},
		}},
		{"oscategories", "", []Listable{
			&OSCategory{},
			&ListOSCategories{},
		}},
		{"isos", "", []Listable{
			&ISO{},
			&ListISOs{},
		}},
	}

	for i, t := range ts {
		end := len(t.name) - 1
		if strings.HasSuffix(t.name, "ses") {
			end--
		}
		if strings.HasSuffix(t.name, "jobs") {
			end++
		}
		fieldName := t.name[:end]
		if strings.HasSuffix(fieldName, "ie") {
			fieldName = t.name[:end-2] + "y"
		}

		ts[i].fieldName = fieldName
	}

	return ts
}

func TestClientList(t *testing.T) {
	body := `
	{"list%sresponse": {
		"count": 4,
		"%s": [{}, {}, {}, {}]
	}}`

	for _, tt := range lsTests() {
		responses := make([]response, len(tt.listables))
		for i := range tt.listables {
			responses[i] = response{200, jsonContentType, fmt.Sprintf(body, tt.name, tt.fieldName)}
		}
		ts := newServer(responses...)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		for _, ls := range tt.listables {
			things, err := cs.List(ls)
			if err != nil {
				t.Error(err)
			}

			if len(things) != 4 {
				t.Errorf("four %T were expected, got %d", ls, len(things))
			}

		}

		ts.Close()
	}
}

func TestClientPaginate(t *testing.T) {
	body := `
	{"list%sresponse": {
		"count": 4,
		"%s": [{}, {}, {}, {}]
	}}`

	for _, tt := range lsTests() {
		end := len(tt.name)
		switch {
		case strings.HasSuffix(tt.name, "ses"):
			end -= 2 // nolint: ineffassign
		case strings.HasSuffix(tt.name, "jobs"):
			break
		default:
			end-- // nolint: ineffassign
		}
		responses := make([]response, len(tt.listables))
		for i := range tt.listables {
			responses[i] = response{200, jsonContentType, fmt.Sprintf(body, tt.name, tt.fieldName)}
		}
		ts := newServer(responses...)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		for _, ls := range tt.listables {
			req, _ := ls.ListRequest()
			counter := 0

			cs.Paginate(req, func(i interface{}, e error) bool {
				if e != nil {
					t.Error(e)
					return false
				}

				counter++
				return true
			})

			if counter != 4 {
				t.Errorf("Four %s were expected, got %d", tt.name, counter)
			}
		}

		ts.Close()
	}
}

func TestClientPaginateError(t *testing.T) {
	body := `
	{"list%sresponse": {
		"cserrorcode": 9999,
		"errorcode": 431,
		"errortext": "Unable to execute API command listzones due to invalid value. Invalid parameter id value=1747ef5e-5451-41fd-9f1a-58913bae9701 due to incorrect long value format, or entity does not exist or due to incorrect parameter annotation for the field in api cmd class.",
		"uuidList": []
	}}
`
	for _, tt := range lsTests() {
		responses := make([]response, len(tt.listables))
		for i := range tt.listables {
			responses[i] = response{431, jsonContentType, fmt.Sprintf(body, tt.name)}
		}
		ts := newServer(responses...)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		for i := range tt.listables {
			listable := tt.listables[i]
			cs.Paginate(listable, func(i interface{}, e error) bool {
				t.Errorf("no %T were expected %v %s", listable, i, e)
				return false
			})
		}

		ts.Close()
	}
}

func TestClient_Do(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(http.MethodGet, "/",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusNotFound, ""), nil
		})

	httpmock.RegisterResponder(http.MethodPost, "/",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusBadRequest, `{"message":"not this way"}`), nil
		})

	httpmock.RegisterResponder(http.MethodPost, "/broken",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusInternalServerError, "API is broken"), nil
		})

	httpmock.RegisterResponder(http.MethodGet, "/ok",
		func(_ *http.Request) (*http.Response, error) {
			return httpmock.NewStringResponse(http.StatusOK, "test"), nil
		})

	client := NewClient("x", "x", "x")

	// Test for ErrNotFound when receiving a http.StatusNotFound status
	req, err := http.NewRequest(http.MethodGet, "http://example.net/", nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.True(t, errors.Is(err, ErrNotFound))
	require.Nil(t, resp)

	// Test for ErrInvalidRequest when receiving a http.StatusBadRequest status
	req, err = http.NewRequest(http.MethodPost, "http://example.net/", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.True(t, errors.Is(err, ErrInvalidRequest))
	require.True(t, strings.Contains(err.Error(), "not this way"))
	require.Nil(t, resp)

	// Test for ErrAPIError when receiving a http.StatusInternalServerError status
	req, err = http.NewRequest(http.MethodPost, "http://example.net/broken", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.True(t, errors.Is(err, ErrAPIError))
	require.Nil(t, resp)

	// Test for successful request
	req, err = http.NewRequest(http.MethodGet, "http://example.net/ok", nil)
	require.NoError(t, err)
	resp, err = client.Do(req)
	require.NoError(t, err)
	actual, _ := ioutil.ReadAll(resp.Body)
	require.Equal(t, []byte("test"), actual)
}

func TestNewClient(t *testing.T) {
	var (
		testHTTPTransport = http.Transport{}
		testHTTPClient    = &http.Client{Transport: &testHTTPTransport}
		testTimeout       = 5 * time.Second
	)

	client := NewClient("x", "x", "x",
		WithHTTPClient(testHTTPClient),
		WithTimeout(testTimeout),
		WithTrace(),
	)

	require.Equal(t, testHTTPClient, client.HTTPClient)
	require.Equal(t, testTimeout, client.Timeout)
	require.IsType(t, &traceTransport{}, client.HTTPClient.Transport)
	require.NotNil(t, client.Client)

	// Test embeded v2.Client disabling
	client = NewClient("x", "x", "x", WithoutV2Client())
	require.Nil(t, client.Client)
}
