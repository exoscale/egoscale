package egoscale

import (
	"context"
	"fmt"
	"testing"
)

func TestRunstatusIncidentGenericError(t *testing.T) { // nolint: dupl
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	p := response{200, jsonContentType, fmt.Sprintf(`{"subdomain": "testpage", "incidents_url": %q}`, ts.URL)}
	r200 := response{200, jsonContentType, `ERROR`}
	r400 := response{400, jsonContentType, `{"detail": "error"}`}

	ts.addResponse(r200, r400)
	_, err := cs.ListRunstatusIncidents(context.TODO(), RunstatusPage{IncidentsURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.ListRunstatusIncidents(context.TODO(), RunstatusPage{IncidentsURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, r400)
	_, err = cs.GetRunstatusIncident(context.TODO(), RunstatusIncident{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.GetRunstatusIncident(context.TODO(), RunstatusIncident{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, p, r200, r400, p, r400)
	_, err = cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	_, err = cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}
}

func TestRunstatusListIncident(t *testing.T) {
	is := response{200, jsonContentType, `
{
  "incidents": [
    {
      "id": 90,
      "url": "https://example.org/pages/testpage/incidents/90",
      "services": [
        "API"
      ],
      "start_date": "2018-11-14T15:37:29Z",
      "end_date": "2018-11-14T15:38:19Z",
      "status": "resolved",
      "events": [
        {
          "created": "2018-11-14T15:38:19Z",
          "text": "fini",
          "status": "resolved",
          "state": "operational"
        },
        {
          "created": "2018-11-14T15:38:09Z",
          "text": "c'est la vie!",
          "status": "identified",
          "state": "degraded_performance"
        },
        {
          "created": "2018-11-14T15:37:29Z",
          "text": "Foo bar",
          "status": "monitoring",
          "state": "operational"
        }
      ],
      "status_text": "fini",
      "state": "degraded_performance",
      "title": "AAAH",
      "events_url": "https://example.org/pages/testpage/incidents/90/events",
      "post_mortem": "# foo bar\n\nc'est la life",
      "real_time": true
    }
  ]
}`}

	ts := newServer(is)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	incidents, err := cs.ListRunstatusIncidents(context.TODO(), RunstatusPage{IncidentsURL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	if len(incidents) != 1 {
		t.Errorf("1 incident expected: got %d", len(incidents))
	}

	if incidents[0].ID != 90 {
		t.Errorf("id 90 expected: got %d", incidents[0].ID)
	}

	i := response{200, jsonContentType, `{
  "id": 90,
  "status": "degraded_performance"
}`}
	p := response{200, jsonContentType, fmt.Sprintf(`{
  "url": "https://api.runstatus.com/pages/testpage",
  "incidents_url": %q,
  "subdomain": "testpage"
}`, ts.URL)}

	ts.addResponse(i)
	incident, err := cs.GetRunstatusIncident(context.TODO(), RunstatusIncident{URL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	if incident.ID != 90 {
		t.Errorf("bad incident, %#v", incident)
	}

	ts.addResponse(p, is)
	incident, err = cs.GetRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL, Title: "AAAH"})
	if err != nil {
		t.Fatal(err)
	}

	if incident.ID != 90 {
		t.Errorf("bad incident, %#v", incident)
	}

	ts.addResponse(p, is)
	incident, err = cs.GetRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL, ID: 90})
	if err != nil {
		t.Fatal(err)
	}

	if incident.Title != "AAAH" {
		t.Errorf("bad incident, %#v", incident)
	}
}

func TestRunstatusIncidentDelete(t *testing.T) {
	ts := newServer(response{204, jsonContentType, ""})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if err := cs.DeleteRunstatusIncident(context.TODO(), RunstatusIncident{}); err == nil {
		t.Error("incident without a status should fail")
	}

	if err := cs.DeleteRunstatusIncident(context.TODO(), RunstatusIncident{URL: ts.URL}); err != nil {
		t.Error(err)
	}
}

func TestRunstatusIncidentCreate(t *testing.T) {
	ts := newServer()
	defer ts.Close()
	ts.addResponse(
		response{200, jsonContentType, fmt.Sprintf(`{"incidents_url": %q, "subdomain": "d"}`, ts.URL)},
		response{201, jsonContentType, `{"url": "...", "name": "hello"}`},
	)

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if _, err := cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{}); err == nil {
		t.Error("incident without a status should fail")
	}

	if _, err := cs.CreateRunstatusIncident(context.TODO(), RunstatusIncident{PageURL: ts.URL}); err != nil {
		t.Error(err)
	}
}
