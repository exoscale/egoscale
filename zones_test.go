package egoscale

import (
	"context"
	"testing"
	"time"
)

func TestListZonesAPIName(t *testing.T) {
	req := &ListZones{}
	_ = req.Response().(*ListZonesResponse)
}

func TestListZonesTypeError(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listzonesresponse": []}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	_, err := cs.List(&Zone{})
	if err == nil {
		t.Errorf("An error was expected")
	}
}

func TestListZonesPaginateBreak(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zone": [
		{
			"id": "1747ef5e-5451-41fd-9f1a-58913bae9702",
			"name": "ch-gva-2",
			"tags": []
		},
		{
			"id": "381d0a95-ed4a-4ad9-b41c-b97073c1a433",
			"name": "ch-dk-2",
			"tags": []
		}
	]
}}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	zone := new(Zone)
	req, _ := zone.ListRequest()

	cs.Paginate(req, func(i interface{}, e error) bool {
		if e != nil {
			t.Error(e)
			return false
		}
		z := i.(*Zone)
		if z.Name == "" {
			t.Errorf("Zone Name not set")
		}
		return false
	})
}

func TestListZonesAsyncError(t *testing.T) {
	ts := newServer(response{431, jsonContentType, `
{
	"listzonesresponse": {
		"cserrorcode": 9999,
		"errorcode": 431,
		"errortext": "Unable to execute API command listzones due to invalid value. Invalid parameter id value=1747ef5e-5451-41fd-9f1a-58913bae9701 due to incorrect long value format, or entity does not exist or due to incorrect parameter annotation for the field in api cmd class.",
		"uuidList": []
	}
}
`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	zone := &Zone{
		ID: MustParseUUID("1747ef5e-5451-41fd-9f1a-58913bae9701"),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	outChan, errChan := cs.AsyncListWithContext(ctx, zone)

	for {
		select {
		case _, ok := <-outChan:
			if ok {
				t.Errorf("no zones were expected")
			} else {
				outChan = nil
			}
		case e, ok := <-errChan:
			if ok {
				t.Errorf("no errors were expected, got %s", e)
			}
			errChan = nil
		}

		if outChan == nil && errChan == nil {
			break
		}
	}
}

func TestListZonesAsync(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zone": [
		{
			"id": "1747ef5e-5451-41fd-9f1a-58913bae9702",
			"name": "ch-gva-2",
			"tags": []
		},
		{
			"id": "381d0a95-ed4a-4ad9-b41c-b97073c1a433",
			"name": "ch-dk-2",
			"tags": []
		},
		{
			"id": "b0fcd72f-47ad-4779-a64f-fe4de007ec72",
			"name": "at-vie-1",
			"tags": []
		},
		{
			"id": "de88c980-78f6-467c-a431-71bcc88e437f",
			"name": "de-fra-1",
			"tags": []
		}
	]
}}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	zone := new(Zone)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	outChan, errChan := cs.AsyncListWithContext(ctx, zone)

	counter := 0
	for {
		select {
		case z, ok := <-outChan:
			if ok {
				zone := z.(*Zone)
				if zone.Name == "" {
					t.Errorf("Zone Name is empty")
				}
				counter++
			} else {
				outChan = nil
			}
		case e, ok := <-errChan:
			if ok {
				t.Error(e)
			}
			errChan = nil
		}

		if outChan == nil && errChan == nil {
			break
		}
	}

	if counter != 4 {
		t.Errorf("Four zones were expected, got %d", counter)
	}
}

func TestListZonesTwoPages(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zone": [
		{
			"id": "1747ef5e-5451-41fd-9f1a-58913bae9702",
			"name": "ch-gva-2",
			"tags": []
		},
		{
			"id": "381d0a95-ed4a-4ad9-b41c-b97073c1a433",
			"name": "ch-dk-2",
			"tags": []
		}
	]
}}`}, response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zone": [
		{
			"id": "b0fcd72f-47ad-4779-a64f-fe4de007ec72",
			"name": "at-vie-1",
			"tags": []
		},
		{
			"id": "de88c980-78f6-467c-a431-71bcc88e437f",
			"name": "de-fra-1",
			"tags": []
		}
	]
}}`}, response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zones": null
}}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")
	cs.PageSize = 2

	zone := new(Zone)
	zones, err := cs.List(zone)
	if err != nil {
		t.Error(err)
	}

	if len(zones) != 4 {
		t.Errorf("Four zones were expected, got %d", len(zones))
	}
}

func TestListZonesError(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listzonesresponse": {
	"count": 4,
	"zone": [
		{
			"id": "1747ef5e-5451-41fd-9f1a-58913bae9702",
			"name": "ch-gva-2",
			"tags": []
		},
		{
			"id": "381d0a95-ed4a-4ad9-b41c-b97073c1a433",
			"name": "ch-dk-2",
			"tags": []
		}
	]
}}`}, response{400, jsonContentType, `
{"listzonesresponse": {
	"cserrorcode": 9999,
	"errorcode": 431,
	"errortext": "Unable to execute API command listzones",
	"uuidList": []
}}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")
	cs.PageSize = 2

	zone := new(Zone)
	_, err := cs.List(zone)
	if err == nil {
		t.Error("An error was expected")
	}
}

func TestListZonesTimeout(t *testing.T) {
	ts := newSleepyServer(`
{"listzonesresponse": {
	"count": 4
}}`)
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")
	cs.HTTPClient.Timeout = time.Millisecond

	zone := new(Zone)
	_, err := cs.List(zone)
	if err == nil {
		t.Errorf("An error was expected")
	}
}
