package egoscale

import (
	"testing"
)

func TestAddIPToNic(t *testing.T) {
	req := &AddIPToNic{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*NicSecondaryIP)
}

func TestRemoveIPFromNic(t *testing.T) {
	req := &RemoveIPFromNic{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestListNicsAPIName(t *testing.T) {
	req := &ListNics{}
	_ = req.response().(*ListNicsResponse)
}

func TestActivateIP6(t *testing.T) {
	req := &ActivateIP6{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Nic)
}
