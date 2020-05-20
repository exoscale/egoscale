package egoscale

import (
	"context"
	"fmt"
	"testing"
)

func TestRunstatusMaintenanceGenericError(t *testing.T) { // nolint: dupl
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	p := response{200, jsonContentType, fmt.Sprintf(`{"subdomain": "testpage", "services_url": %q}`, ts.URL)}
	r200 := response{200, jsonContentType, `ERROR`}
	r400 := response{400, jsonContentType, `{"detail": "error"}`}

	ts.addResponse(r200, r400)
	_, err := cs.ListRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.ListRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, r400)
	_, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	ts.addResponse(r200, p, r200, r400, p, r400)
	_, err = cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 200 bad json: got nil")
	}

	_, err = cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}

	_, err = cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL})
	if err == nil {
		t.Errorf("error expected on 400: got nil")
	}
}

func TestRunstatusListMaintenances(t *testing.T) {
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	ms := response{200, jsonContentType, `
{
  "next": null,
  "previous": null,
  "results": [
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

	m := response{200, jsonContentType, `{
  "url": "...",
  "title": "hggfh",
  "status": "scheduled"
}`}
	p := response{200, jsonContentType, fmt.Sprintf(`{
  "url": "https://api.runstatus.com/pages/bauud",
  "maintenances_url": %q,
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
  ],
  "subdomain": "bauud"
}`, ts.URL)}

	ts.addResponse(ms)
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

	ts.addResponse(m)
	maintenance, err := cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.Title != "hggfh" {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
	}

	ts.addResponse(p)
	maintenance, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL, ID: 598})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.Title != "hggfh" {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
	}

	ts.addResponse(p)
	maintenance, err = cs.GetRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL, Title: "hggfh"})
	if err != nil {
		t.Fatal(err)
	}
	if maintenance.ID != 598 {
		t.Errorf("bad maintenance fetched, got %#v", maintenance)
	}
}

func TestRunstatusPaginateMaintenances(t *testing.T) {
	ts := newServer()
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	ps := response{200, jsonContentType, fmt.Sprintf(`{
    "next": %q,
    "previous":null,
    "results": [
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
  }`, ts.URL)}

	ms := response{200, jsonContentType, `
{
  "next": null,
  "previous":null,
  "results": [
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

	ts.addResponse(ps, ms)
	cs.PaginateRunstatusMaintenances(context.TODO(), RunstatusPage{MaintenancesURL: ts.URL}, func(maintenance *RunstatusMaintenance, e error) bool {
		if e != nil {
			t.Errorf(`reply error not expected: %v`, e)
		}

		if maintenance.Title != "hggfh" {
			t.Errorf(`hggfh was expected but got %q`, maintenance.Title)
		}

		return false
	})
}

func TestRunstatusMaintenanceDelete(t *testing.T) {
	ts := newServer(response{204, jsonContentType, ""})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if err := cs.DeleteRunstatusMaintenance(context.TODO(), RunstatusMaintenance{}); err == nil {
		t.Error("service without a status should fail")
	}

	if err := cs.DeleteRunstatusMaintenance(context.TODO(), RunstatusMaintenance{URL: ts.URL}); err != nil {
		t.Error(err)
	}
}

func TestRunstatusMaintenanceCreate(t *testing.T) {
	ts := newServer()
	defer ts.Close()
	ts.addResponse(
		response{200, jsonContentType, fmt.Sprintf(`{"maintenances_url": %q, "subdomain": "d"}`, ts.URL)},
		response{201, jsonContentType, `{"id": 1, "url": "...", "name": "hello"}`},
	)

	cs := NewClient(ts.URL, "KEY", "SECRET")

	if _, err := cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{}); err == nil {
		t.Error("service without a status should fail")
	}

	if _, err := cs.CreateRunstatusMaintenance(context.TODO(), RunstatusMaintenance{PageURL: ts.URL}); err != nil {
		t.Error(err)
	}
}
