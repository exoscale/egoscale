package egoscale

import (
	"testing"
)

func TestTemplateResourceType(t *testing.T) {
	instance := &Template{}
	if instance.ResourceType() != "Template" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestCreateTemplate(t *testing.T) {
	req := &CreateTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Template)
}

func TestCopyTemplate(t *testing.T) {
	req := &CopyTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Template)
}

func TestUpdateTemplate(t *testing.T) {
	req := &UpdateTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Template)
}

func TestListTemplates(t *testing.T) {
	req := &ListTemplates{}
	_ = req.Response().(*ListTemplatesResponse)
}

func TestDeleteTemplate(t *testing.T) {
	req := &DeleteTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*booleanResponse)
}

func TestPrepareTemplate(t *testing.T) {
	req := &PrepareTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Template)
}

func TestRegisterTemplate(t *testing.T) {
	req := &RegisterTemplate{}
	_ = req.Response().(*Template)
}

func TestListOSCategories(t *testing.T) {
	req := &ListOSCategories{}
	_ = req.Response().(*ListOSCategoriesResponse)
}

func TestTemplate(t *testing.T) {
	instance := &Template{}
	if instance.ResourceType() != "Template" {
		t.Errorf("ResourceType doesn't match")
	}
}
