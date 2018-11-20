package egoscale

import (
	"fmt"
	"net"
	"strings"
	"testing"
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
		&NetworkOffering{},
		(*ListNetworkOfferings)(nil),
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
		{"networkofferings", &NetworkOffering{}},
		{"nics", &Nic{}},
		{"snapshots", &Snapshot{}},
	}

	for _, thing := range things {
		ts := newServer(
			response{200, jsonContentType, fmt.Sprintf(body, thing.name)},
			response{431, jsonContentType, bodyError},
		)

		cs := NewClient(ts.URL, "KEY", "SECRET")

		errText := "not found"
		_, err := cs.Get(thing.listable)
		if err == nil {
			t.Error("an error was expected")
			continue
		}

		e, ok := err.(*ErrorResponse)
		if !ok {
			t.Errorf("an ErrorResponse was expected, got %T", err)
			continue
		}

		if !strings.Contains(e.ErrorText, errText) {
			t.Errorf("bad error test, got %q", e.ErrorText)
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
		{"networkoffering", &NetworkOffering{ID: id}},
		{"nic", &Nic{ID: id}},
		{"snapshot", &Snapshot{ID: id}},
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

		if !strings.HasPrefix(err.Error(), "API error ParamError 431") {
			t.Errorf("bad error %q", err)
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
		"networkoffering": [{}, {}],
		"nic": [{}, {}],
		"snapshot": [{}, {}]
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

		if !strings.HasPrefix(err.Error(), "more than one") {
			t.Errorf("bad error %s", err)
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
	listables []Listable
}

func lsTests() []lsTest {
	return []lsTest{
		{"zones", []Listable{
			&Zone{},
			&ListZones{},
		}},
		{"publicipaddresses", []Listable{
			&IPAddress{},
			&ListPublicIPAddresses{},
		}},
		{"sshkeypairs", []Listable{
			&SSHKeyPair{},
			&ListSSHKeyPairs{},
		}},
		{"affinitygroups", []Listable{
			&AffinityGroup{},
			&ListAffinityGroups{},
		}},
		{"securitygroups", []Listable{
			&SecurityGroup{},
			&ListSecurityGroups{},
		}},
		{"virtualmachines", []Listable{
			&VirtualMachine{},
			&ListVirtualMachines{},
		}},
		{"volumes", []Listable{
			&Volume{},
			&ListVolumes{},
		}},
		{"templates", []Listable{
			&Template{IsFeatured: true},
			&ListTemplates{TemplateFilter: "featured"},
		}},
		{"serviceofferings", []Listable{
			&ServiceOffering{},
			&ListServiceOfferings{},
		}},
		{"networks", []Listable{
			&Network{},
			&ListNetworks{},
		}},
		{"networkofferings", []Listable{
			&NetworkOffering{},
			&ListNetworkOfferings{},
		}},
		{"accounts", []Listable{
			&Account{},
			&ListAccounts{},
		}},
		{"nics", []Listable{
			&Nic{},
			&ListNics{},
		}},
		{"snapshots", []Listable{
			&Snapshot{},
			&ListSnapshots{},
		}},
		{"events", []Listable{
			&Event{},
			&ListEvents{},
		}},
		{"resourcelimits", []Listable{
			&ResourceLimit{},
			&ListResourceLimits{},
		}},
		{"tags", []Listable{
			&ResourceTag{},
			&ListTags{},
		}},
		{"users", []Listable{
			&User{},
			&ListUsers{},
		}},
	}
}

func TestClientList(t *testing.T) {
	body := `
	{"list%sresponse": {
		"count": 4,
		"%s": [{}, {}, {}, {}]
	}}`

	for _, tt := range lsTests() {
		end := len(tt.name) - 1
		if strings.HasSuffix(tt.name, "ses") {
			end--
		}
		responses := make([]response, len(tt.listables))
		for i := range tt.listables {
			responses[i] = response{200, jsonContentType, fmt.Sprintf(body, tt.name, tt.name[:end])}
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
		end := len(tt.name) - 1
		if strings.HasSuffix(tt.name, "ses") {
			end--
		}
		responses := make([]response, len(tt.listables))
		for i := range tt.listables {
			responses[i] = response{200, jsonContentType, fmt.Sprintf(body, tt.name, tt.name[:end])}
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
				t.Errorf("Four zones were expected, got %d", counter)
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
