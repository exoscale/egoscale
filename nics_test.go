package egoscale

import (
	"testing"
)

func TestAddIPToNic(t *testing.T) {
	req := &AddIPToNic{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*NicSecondaryIP)
}

func TestRemoveIPFromNic(t *testing.T) {
	req := &RemoveIPFromNic{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*booleanResponse)
}

func TestListNicsAPIName(t *testing.T) {
	req := &ListNics{}
	_ = req.Response().(*ListNicsResponse)
}

func TestActivateIP6(t *testing.T) {
	req := &ActivateIP6{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Nic)
}
