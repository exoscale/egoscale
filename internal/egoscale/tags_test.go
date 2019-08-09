package egoscale

import (
	"testing"
)

func TestCreateTags(t *testing.T) {
	req := &CreateTags{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestDeleteTags(t *testing.T) {
	req := &DeleteTags{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestListTags(t *testing.T) {
	req := &ListTags{}
	_ = req.Response().(*ListTagsResponse)
}
