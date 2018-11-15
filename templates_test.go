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
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Template)
}

func TestCopyTemplate(t *testing.T) {
	req := &CopyTemplate{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Template)
}

func TestUpdateTemplate(t *testing.T) {
	req := &UpdateTemplate{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Template)
}

func TestListTemplates(t *testing.T) {
	req := &ListTemplates{}
	_ = req.response().(*ListTemplatesResponse)
}

func TestDeleteTemplate(t *testing.T) {
	req := &DeleteTemplate{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestPrepareTemplate(t *testing.T) {
	req := &PrepareTemplate{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Template)
}

func TestRegisterTemplate(t *testing.T) {
	req := &RegisterTemplate{}
	_ = req.response().(*Template)
}

func TestListOSCategories(t *testing.T) {
	req := &ListOSCategories{}
	_ = req.response().(*ListOSCategoriesResponse)
}

func TestTemplate(t *testing.T) {
	instance := &Template{}
	if instance.ResourceType() != "Template" {
		t.Errorf("ResourceType doesn't match")
	}
}
