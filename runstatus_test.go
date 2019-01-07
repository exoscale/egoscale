package egoscale

import (
	"context"
	"fmt"
	"testing"
)

func TestRunstatusPage(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{
  "id": 102,
  "url": "https://example.org/pages/testpage",
  "created": "2018-11-14T15:21:10Z",
  "plan": "free",
  "subdomain": "testpage",
  "domain": null,
  "ok_text": "All systems operational",
  "state": "operational",
  "time_zone": "UTC",
  "title": "",
  "support_email": "",
  "services_url": "https://example.org/pages/testpage/services",
  "incidents_url": "https://example.org/pages/testpage/incidents",
  "maintenances_url": "https://example.org/pages/testpage/maintenances",
  "logo": null,
  "twitter_username": "",
  "public_url": "https://testpage.runstat.us",
  "dark_theme": false,
  "gradient_start": "224,224,224,0.9",
  "gradient_end": "255,255,255,0.9",
  "title_color": "204,204,204,1",
  "header_background": null,
  "services": [
    {
      "url": "https://example.org/pages/testpage/services/28",
      "name": "API",
      "state": "operational"
    }
  ],
  "maintenances": [],
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
}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	page, err := cs.GetRunstatusPage(context.TODO(), RunstatusPage{URL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	if page.Subdomain != "testpage" {
		t.Errorf("subpage should be %q, got %q", "testpage", page.Subdomain)
	}
}

func TestCreateRunstatusPage(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{
  "id": 102,
  "url": "https://example.org/pages/testpage",
  "created": "2018-11-14T15:21:10Z",
  "plan": "free",
  "subdomain": "testpage",
  "domain": null,
  "ok_text": "All systems operational",
  "state": "operational",
  "time_zone": "UTC",
  "title": "",
  "support_email": "",
  "services_url": "https://example.org/pages/testpage/services",
  "incidents_url": "https://example.org/pages/testpage/incidents",
  "maintenances_url": "https://example.org/pages/testpage/maintenances",
  "logo": null,
  "twitter_username": "",
  "public_url": "https://testpage.runstat.us",
  "dark_theme": false,
  "gradient_start": "224,224,224,0.9",
  "gradient_end": "255,255,255,0.9",
  "title_color": "204,204,204,1",
  "header_background": null,
  "services": [
    {
      "url": "https://example.org/pages/testpage/services/28",
      "name": "API",
      "state": "operational"
    }
  ],
  "maintenances": [],
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
}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	page, err := cs.CreateRunstatusPage(context.TODO(), RunstatusPage{})
	if err != nil {
		t.Fatal(err)
	}

	if page.Subdomain != "testpage" {
		t.Errorf("subpage should be %q, got %q", "testpage", page.Subdomain)
	}
}

func TestListRunstatusPage(t *testing.T) {
	ps := response{200, jsonContentType, `
{
  "count":9,
  "next":null,
  "previous":null,
  "results":[
    {
      "id": 102,
      "url": "https://example.org/pages/testpage",
      "created": "2018-11-14T15:21:10Z",
      "plan": "free",
      "subdomain": "testpage",
      "domain": null,
      "ok_text": "All systems operational",
      "state": "operational",
      "time_zone": "UTC",
      "title": "",
      "support_email": "",
      "services_url": "https://example.org/pages/testpage/services",
      "incidents_url": "https://example.org/pages/testpage/incidents",
      "maintenances_url": "https://example.org/pages/testpage/maintenances",
      "logo": null,
      "twitter_username": "",
      "public_url": "https://testpage.runstat.us",
      "dark_theme": false,
      "gradient_start": "224,224,224,0.9",
      "gradient_end": "255,255,255,0.9",
      "title_color": "204,204,204,1",
      "header_background": null,
      "services": [
        {
          "url": "https://example.org/pages/testpage/services/28",
          "name": "API",
          "state": "operational"
        }
      ],
      "maintenances": [],
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
    }
  ]
}`}

	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	ts.addResponse(ps)
	pages, err := cs.ListRunstatusPages(context.TODO())
	if err != nil {
		t.Fatal(err)
	}

	if len(pages) != 1 {
		t.Errorf("1 page expected: got %d", len(pages))
	}

	if pages[0].Subdomain != "testpage" {
		t.Errorf("subpage should be %q, got %q", "testpage", pages[0].Subdomain)
	}

	ts.addResponse(ps)
	page, err := cs.GetRunstatusPage(context.TODO(), RunstatusPage{Subdomain: "testpage"})
	if err != nil {
		t.Fatal(err)
	}
	if page.ID != 102 {
		t.Errorf("bad page ID, got %d, wanted 102", page.ID)
	}

	ts.addResponse(ps)
	page, err = cs.GetRunstatusPage(context.TODO(), RunstatusPage{ID: 102})
	if err != nil {
		t.Fatal(err)
	}
	if page.Subdomain != "testpage" {
		t.Errorf(`bad page ID, got %q, wanted "testpage"`, page.Subdomain)
	}
}

func TestRunstatusListService(t *testing.T) {
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

	ts.addResponse(p, ss)
	service, err := cs.GetRunstatusService(context.TODO(), RunstatusService{PageURL: ts.URL, Name: "API"})
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

func TestRunstatusDeleteService(t *testing.T) {
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

func TestRunstatusListMaintenance(t *testing.T) {
	ms := response{200, jsonContentType, `
{
  "maintenances": [
    {
      "url":"https://api.runstatus.com/pages/bauud/maintenances/598",
      "created":"2018-11-27T12:51:05.607060Z",
      "services":["infra"],
      "start_date":"2018-11-27T13:50:00Z",
      "end_date":"2018-11-28T13:50:00Z",
      "title":"hggfh",
      "description":"hgfhgfgf",
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
        }
      ],
      "status":"scheduled",
      "events_url":"https://api.runstatus.com/pages/bauud/maintenances/598/events",
      "real_time":true
    },
    {
      "url":"https://api.runstatus.com/pages/bauud/maintenances/600",
      "created":"2018-11-27T12:51:05.607060Z",
      "services":["infra"],
      "start_date":"2018-11-27T13:50:00Z",
      "end_date":"2018-11-28T13:50:00Z",
      "title":"hggfh",
      "description":"hgfhgfgf",
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
        }
      ],
      "status":"scheduled",
      "events_url":"https://api.runstatus.com/pages/bauud/maintenances/600/events",
      "real_time":true
    }
  ]
}`}
	ts := newServer(ms)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	maintenances, err := cs.ListRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}

	if len(maintenances) != 2 {
		t.Errorf("2 maintenance expected: got %d", len(maintenances))
	}

	if maintenances[0].Title != "hggfh" {
		t.Errorf("title should be %q, got %q", "hggfh", maintenances[0].Title)
	}

	if maintenances[0].ID != 598 {
		t.Errorf("maintenance ID should be 598, got %d", maintenances[0].ID)
	}

	m := response{200, jsonContentType, fmt.Sprintf(`{
  "url": "...",
  "title": "hggfh",
  "status": "scheduled"
}`)}
	p := response{200, jsonContentType, fmt.Sprintf(`{
  "url": "https://api.runstatus.com/pages/bauud",
  "maintenances_url": %q,
  "subdomain": "bauud"
}`, ts.URL)}

	ts.addResponse(m)
	maintenance, err := cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.Title != "hggfh" {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
	}

	ts.addResponse(p, ms)
	maintenance, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL, ID: 598})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.Title != "hggfh" {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
	}

	ts.addResponse(p, ms)
	maintenance, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL, Title: "hggfh"})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.ID != 598 {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
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

func TestRunstatusGenericError(t *testing.T) {
	ts := newServer(
		response{200, jsonContentType, `ERROR`},
		response{400, jsonContentType, `{"detail": "error"}`},
	)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	ts.addResponse()

	_, err := cs.ListRunstatusServices(context.TODO(), RunstatusPage{ServicesURL: "testpage"})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.ListRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: "testpage"})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}
}
