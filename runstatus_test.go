package egoscale

import (
	"context"
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
	ts := newServer(response{200, jsonContentType, `
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
}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

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
}

func TestRunstatusListService(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{
  "services": [
    {
      "url": "https://example.org/pages/testpage/services/ERROR",
      "name": "API",
      "state": "operational"
    },
    {
      "url": "https://example.org/pages/testpage/services/29",
      "name": "API",
      "state": "operational"
    }
  ]
}`})
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
}

func TestRunstatusListMaintenance(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
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
      "url":"https://api.runstatus.com/pages/bauud/maintenances/ERROR",
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
      "events_url":"https://api.runstatus.com/pages/bauud/maintenances/ERROR/events",
      "real_time":true
    }
  ]
}`})
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
}

func TestRunstatusListIncident(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
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
}`})
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
}

func TestRunstatusGenericError(t *testing.T) {
	errorCode := 200

	for errorCode <= 400 {
		ts := newServer(response{errorCode, jsonContentType, `
    {
      ERROR
    }
    `})

		cs := NewClient(ts.URL, "KEY", "SECRET")

		i := 1
		_, err := cs.ListRunstatusServices(context.TODO(), RunstatusPage{ServicesURL: "testpage"})
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}
		i++
		_, err = cs.ListRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: "testpage"})
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}
		i++
		_, err = cs.ListRunstatusIncidents(context.TODO(), RunstatusPage{IncidentsURL: "testpage"})
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}
		i++
		_, err = cs.CreateRunstatusPage(context.TODO(), RunstatusPage{})
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}
		i++
		_, err = cs.GetRunstatusPage(context.TODO(), RunstatusPage{Subdomain: "testpage"})
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}
		i++
		_, err = cs.ListRunstatusPages(context.TODO())
		if err == nil {
			t.Errorf("TestRunstatusGenericError %d error expected: got nil", i)
		}

		ts.Close()
		errorCode += 200
	}
}

func TestRunstatusGenericErrorWithoutResp(t *testing.T) {

	ts := newServer(response{400, "ERROR", `
    {
      ERROR
    }
    `})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	i := 1
	err := cs.DeleteRunstatusService(context.TODO(), RunstatusService{URL: ts.URL})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.CreateRunstatusService(context.TODO(), RunstatusPage{ServicesURL: ts.URL}, RunstatusService{})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.CreateRunstatusEvent(context.TODO(), RunstatusIncident{EventsURL: ts.URL}, RunstatusEvent{})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.CreateRunstatusMaintenance(context.TODO(), RunstatusPage{MaintenancesURL: ts.URL}, RunstatusMaintenance{})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.DeleteRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.UpdateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{EventsURL: ts.URL}, RunstatusEvent{})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.CreateRunstatusIncident(context.TODO(), RunstatusPage{IncidentsURL: ts.URL}, RunstatusIncident{})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.DeleteRunstatusIncident(context.TODO(), RunstatusIncident{URL: ts.URL})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
	i++
	err = cs.DeleteRunstatusPage(context.TODO(), RunstatusPage{URL: ts.URL})
	if err == nil {
		t.Errorf("TestRunstatusGenericErrorWithoutResp %d error expected: got nil", i)
	}
}
