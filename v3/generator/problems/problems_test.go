package problems

import (
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	out, err := render("v3", map[string]problemType{
		"not-found":         {Status: 404, Title: "The resource is not found."},
		"kms/key-not-found": {Parent: "not-found", Status: 404, Title: "Key not found", Description: "The key does not exist."},
	})
	if err != nil {
		t.Fatal(err)
	}
	got := string(out)
	if !strings.Contains(got, `ErrKmsKeyNotFound = &ProblemType{ID: "#/problems/kms/key-not-found", Status: 404, Title: "Key not found"}`) {
		t.Errorf("missing the service type:\n%s", got)
	}
	if strings.Contains(got, "ErrNotFound =") {
		t.Errorf("the shared types must be left out:\n%s", got)
	}

	if _, err := render("v3", map[string]problemType{
		"kms/key-gone": {Parent: "not-found"},
		"kms/key.gone": {Parent: "not-found"},
	}); err == nil {
		t.Error("types with the same Go name must be refused")
	}
}
