package egoscale

import (
	"net/url"
	"testing"
)

func TestNetwork(t *testing.T) {
	instance := &Network{}
	if instance.ResourceType() != "Network" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestListNetworks(t *testing.T) {
	req := &ListNetworks{}
	_ = req.Response().(*ListNetworksResponse)
}

func TestCreateNetwork(t *testing.T) {
	req := &CreateNetwork{}
	_ = req.Response().(*Network)
}

func TestRestartNetwork(t *testing.T) {
	req := &RestartNetwork{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Network)
}

func TestUpdateNetwork(t *testing.T) {
	req := &UpdateNetwork{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Network)
}

func TestDeleteNetwork(t *testing.T) {
	req := &DeleteNetwork{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestCreateNetworkOnBeforeSend(t *testing.T) {
	req := &CreateNetwork{}
	params := url.Values{}

	if err := req.onBeforeSend(params); err != nil {
		t.Error(err)
	}

	if _, ok := params["name"]; !ok {
		t.Errorf("name should have been set")
	}
	if _, ok := params["displaytext"]; !ok {
		t.Errorf("displaytext should have been set")
	}
}

func TestListNetworkEmpty(t *testing.T) {
	ts := newTestServer(response{200, jsonContentType, `
{"listnetworksresponse": {
	"count": 0,
	"network": []
  }}`})
	defer ts.Close()

	cs := newTestClient(ts.URL)

	network := new(Network)
	networks, err := cs.List(network)
	if err != nil {
		t.Fatal(err)
	}

	if len(networks) != 0 {
		t.Errorf("zero networks were expected, got %d", len(networks))
	}
}

func TestListNetworkFailure(t *testing.T) {
	ts := newTestServer(response{200, jsonContentType, `
{"listnetworksresponse": {
	"count": 3456,
	"network": {}
  }}`})
	defer ts.Close()

	cs := newTestClient(ts.URL)

	network := new(Network)
	networks, err := cs.List(network)
	if err == nil {
		t.Errorf("error was expected, got %v", err)
	}

	if len(networks) != 0 {
		t.Errorf("zero networks were expected, got %d", len(networks))
	}
}

func TestFindNetwork(t *testing.T) {
	ts := newTestServer(response{200, jsonContentType, `
{"listnetworksresponse": {
	"count": 1,
	"network": [
	  {
		"account": "exoscale-1",
		"acltype": "Account",
		"broadcastdomaintype": "Vxlan",
		"canusefordeploy": true,
		"displaytext": "fddd",
		"domain": "exoscale-1",
		"id": "772cee7a-631f-4c0e-ad2b-27776f260d71",
		"ispersistent": true,
		"issystem": false,
		"name": "klmfsdvds",
		"physicalnetworkid": "07f747f5-b445-487f-b2d7-81a5a512989e",
		"related": "772cee7a-631f-4c0e-ad2b-27776f260d71",
		"restartrequired": false,
		"service": [
		  {
			"name": "PrivateNetwork"
		  }
		],
		"specifyipranges": false,
		"state": "Implemented",
		"strechedl2subnet": false,
		"tags": [],
		"traffictype": "Guest",
		"type": "Isolated",
		"zoneid": "1128bd56-b4d9-4ac6-a7b9-c715b187ce11",
		"zonename": "ch-gva-2"
	  }
	]
  }}`})
	defer ts.Close()

	cs := newTestClient(ts.URL)

	networks, err := cs.List(&Network{Name: "klmfsdvds", CanUseForDeploy: true, Type: "Isolated"})
	if err != nil {
		t.Fatal(err)
	}

	if len(networks) != 1 {
		t.Fatalf("One network was expected, got %d", len(networks))
	}

	net, ok := networks[0].(*Network)
	if !ok {
		t.Errorf("unable to type inference *Network, got %v", net)
	}

	if networks[0].(*Network).Name != "klmfsdvds" {
		t.Errorf("klmfsdvds network name was expected, got %s", networks[0].(*Network).Name)
	}
}
