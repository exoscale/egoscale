package egoscale

import (
	"context"
	"fmt"
	"testing"
)

func TestRunstatusServiceGenericError(t *testing.T) {
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	p := response{200, jsonContentType, fmt.Sprintf(`{"subdomain": "testpage", "services_url": %q}`, ts.URL)}
	r200 := response{200, jsonContentType, `ERROR`}
	r400 := response{400, jsonContentType, `{"detail": "error"}`}

	ts.addResponse(r200, r400)
	_, err := cs.ListRunstatusServices(context.TODO(), RunstatusPage{ServicesURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.ListRunstatusServices(context.TODO(), RunstatusPage{ServicesURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, r400)
	_, err = cs.GetRunstatusService(context.TODO(), RunstatusService{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.GetRunstatusService(context.TODO(), RunstatusService{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, p, r200, r400, p, r400)
	_, err = cs.CreateRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	_, err = cs.CreateRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}
}

func TestRunstatusListServices(t *testing.T) {
	ss := response{200, jsonContentType, `
{
  "services": [
    {
      "url": "https://example.org/pages/testpage/services/28",
      "name": "API",
      "state": "operational"
    },
    {
      "url": "https://example.org/pages/testpage/services/29",
      "name": "ABI",
      "state": "hold"
    }
  ]
}`}
	ts := newServer(ss)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	services, err := cs.ListRunstatusServices(context.TODO(), RunstatusPage{ServicesURL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	if len(services) != 2 {
		t.Errorf("2 services expected: got %d", len(services))
	}

	if services[1].URL != "https://example.org/pages/testpage/services/29" {
		t.Errorf("url should be %q, got %q", "https://example.org/pages/testpage/services/29", services[1].URL)
	}

	p := response{200, jsonContentType, fmt.Sprintf(`{
  "url": "https://api.runstatus.com/pages/testpage",
  "services_url": %q,
  "subdomain": "testpage"
}`, ts.URL)}

	ts.addResponse(response{200, jsonContentType, fmt.Sprintf(`{"id": 1, "url": %q, "name": "API"}`, ts.URL)})
	service, err := cs.GetRunstatusService(context.TODO(), RunstatusService{URL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}
	if service.Name != "API" {
		t.Errorf(`bad name, got %q, wanted "API"`, service.Name)
	}

	ts.addResponse(p, ss)
	service, err = cs.GetRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL, Name: "API"})
	if err != nil {
		t.Fatal(err)
	}

	if service.ID != 28 {
		t.Errorf(`bad state, got %d, wanted %d`, service.ID, 28)
	}

	ts.addResponse(p, ss)
	service, err = cs.GetRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL, ID: 29})
	if err != nil {
		t.Fatal(err)
	}

	if service.State != "hold" {
		t.Errorf(`bad state, got %q, wanted "hold"`, service.State)
	}
}

func TestRunstatusServiceDelete(t *testing.T) {
	ts := newServer(response{204, jsonContentType, ""})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if err := cs.DeleteRunstatusService(context.TODO(), RunstatusService{}); err == nil {
		t.Error("service without a status should fail")
	}

	if err := cs.DeleteRunstatusService(context.TODO(), RunstatusService{URL: ts.URL}); err != nil {
		t.Error(err)
	}
}

func TestRunstatusServiceCreate(t *testing.T) {
	ts := newServer()
	defer ts.Close()
	ts.addResponse(
		response{200, jsonContentType, fmt.Sprintf(`{"services_url": %q, "subdomain": "d"}`, ts.URL)},
		response{201, jsonContentType, `{"url": "...", "name": "hello"}`},
	)

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if _, err := cs.CreateRunstatusService(context.TODO(), RunstatusService{}); err == nil {
		t.Error("service without a status should fail")
	}

	if _, err := cs.CreateRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL}); err != nil {
		t.Error(err)
	}
}
