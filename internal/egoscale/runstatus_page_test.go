package egoscale

import (
	"context"
	"fmt"
	"testing"
)

func TestRunstatusError(t *testing.T) {
	ts := newServer(
		response{401, jsonContentType, `{"detail":"Invalid Authorization header"}`},
	)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	_, err := cs.GetRunstatusPage(context.TODO(), RunstatusPage{URL: ts.URL})
	if err == nil {
		t.Errorf("An error was expected")
	}
	if err.Error() != "Runstatus error: Invalid Authorization header" {
		t.Errorf("Bad error response, got %s", err)
	}
}

func TestRunstatusValidationError(t *testing.T) {
	ts := newServer(
		response{401, jsonContentType, `{"foo":["bad mojo"]}`},
	)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	_, err := cs.GetRunstatusPage(context.TODO(), RunstatusPage{URL: ts.URL})
	if err == nil {
		t.Errorf("An error was expected")
	}
	if err.Error() != "Runstatus error: foo: bad mojo" {
		t.Errorf("Bad error response, got %s", err)
	}
}

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
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	p := response{200, jsonContentType, `
{
  "id": 102,
  "url": "https://example.org/pages/testpage",
  "subdomain": "testpage"
}`}
	ps := response{200, jsonContentType, fmt.Sprintf(`
{
  "count":9,
  "next":null,
  "previous":null,
  "results":[
    {
      "id": 102,
      "url": %q,
      "subdomain": "testpage",
      "services": [],
      "maintenances": [],
      "incidents": []
    }
  ]
}`, ts.URL)}

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

	ts.addResponse(ps, p)
	page, err := cs.GetRunstatusPage(context.TODO(), RunstatusPage{Subdomain: "testpage"})
	if err != nil {
		t.Fatal(err)
	}
	if page.ID != 102 {
		t.Errorf("bad page ID, got %d, wanted 102", page.ID)
	}

	ts.addResponse(ps, p)
	page, err = cs.GetRunstatusPage(context.TODO(), RunstatusPage{ID: 102})
	if err != nil {
		t.Fatal(err)
	}
	if page.Subdomain != "testpage" {
		t.Errorf(`bad page ID, got %q, wanted "testpage"`, page.Subdomain)
	}
}

func TestPaginateRunstatusPage(t *testing.T) {
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	p := response{200, jsonContentType, fmt.Sprintf(`
  {
    "count":9,
    "next": %q,
    "previous":null,
    "results":[
      {
        "id": 102,
        "url": "",
        "subdomain": "testpage",
        "services": [],
        "maintenances": [],
        "incidents": []
      }
    ]
  }`, ts.URL)}

	ps := response{200, jsonContentType, `
  {
    "count":64,
    "next":null,
    "previous":null,
    "results":[
      {
        "id": 102,
        "url": "",
        "subdomain": "testpage",
        "services": [],
        "maintenances": [],
        "incidents": []
      }
    ]
  }`}

	ts.addResponse(p, ps)
	cs.PaginateRunstatusPages(context.TODO(), func(pages []RunstatusPage, e error) bool {

		if e != nil {
			t.Errorf(`reply error not expected: %v`, e)
		}

		if pages[0].Subdomain != "testpage" {
			t.Errorf(`testpage2 was expected but got %q`, pages[0].Subdomain)
		}

		return false
	})
}
