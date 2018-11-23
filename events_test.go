package egoscale

import (
	"testing"
)

func TestListEvents(t *testing.T) {
	req := &ListEvents{}
	_ = req.Response().(*ListEventsResponse)
}

func TestListEventTypes(t *testing.T) {
	req := &ListEventTypes{}
	_ = req.Response().(*ListEventTypesResponse)
}
