package egoscale

import (
	"testing"
)

func TestCreateTags(t *testing.T) {
	req := &CreateTags{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*booleanResponse)
}

func TestDeleteTags(t *testing.T) {
	req := &DeleteTags{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*booleanResponse)
}

func TestListTags(t *testing.T) {
	req := &ListTags{}
	_ = req.Response().(*ListTagsResponse)
}
