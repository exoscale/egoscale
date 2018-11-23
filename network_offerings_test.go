package egoscale

import (
	"testing"
)

func TestListNetworkOfferings(t *testing.T) {
	req := &ListNetworkOfferings{}
	_ = req.Response().(*ListNetworkOfferingsResponse)
}

func TestUpdateNetworkOffering(t *testing.T) {
	req := &UpdateNetworkOffering{}
	_ = req.Response().(*NetworkOffering)
}
