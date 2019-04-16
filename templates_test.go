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

func TestListTemplates(t *testing.T) {
	req := &ListTemplates{}
	_ = req.Response().(*ListTemplatesResponse)
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

func TestDeleteTemplate(t *testing.T) {
	req := &DeleteTemplate{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}
