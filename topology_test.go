package egoscale

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetImages(t *testing.T) {
	ts := newServer(`
{
	"listtemplatesresopnse (doesn't matter)": {
		"count": 0,
		"template": [
			{
				"id": "4c0732a0-3df0-4f66-8d16-009f91cf05d6",
				"name": "Linux RedHat 7.4 64-bit",
				"displayText": "Linux RedHat 7.4 64-bit 10G Disk (2017-11-31-dummy)",
				"size": 10737418240
			}
		]
	}
}
	`)
	defer ts.Close()

	cs := NewClient(ts.URL, "TOKEN", "SECRET")
	images, err := cs.GetImages()
	if err != nil {
		log.Fatal(err)
	}

	uuid := "4c0732a0-3df0-4f66-8d16-009f91cf05d6"

	// by short name
	if _, ok := images["redhat-7.4"]; !ok {
		t.Error("expected redhat-7.4 into the map")
	}

	if _, ok := images["redhat-7.4"][10]; !ok {
		t.Error("expected redhat-7.4, 10G into the map")
	}

	if images["redhat-7.4"][10] != uuid {
		t.Error("bad uuid for the redhat-7.4 image")
	}

	// by full name
	if _, ok := images["Linux RedHat 7.4 64-bit"]; !ok {
		t.Error("expected Linux RedHat 7.4 64-bit into the map")
	}

	if _, ok := images["Linux RedHat 7.4 64-bit"][10]; !ok {
		t.Error("expected Linux RedHat 7.4 64-bit, 10G into the map")
	}

	if images["Linux RedHat 7.4 64-bit"][10] != uuid {
		t.Error("bad uuid for the Linux RedHat 7.4 64-bit image")
	}
}

func newServer(response string) *httptest.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(response))
	})
	return httptest.NewServer(mux)
}
