package egoscale

import (
	"testing"
)

func TestPrepareValues(t *testing.T) {
	profile := struct {
		IgnoreMe string
		Zone     string  `json:"myzone,omitempty"`
		Name     string  `json:"name"`
		Id       int     `json:"id"`
		Uid      uint    `json:"uid"`
		Num      float64 `json:"num"`
	}{
		IgnoreMe: "bar",
		Zone:     "hello",
		Name:     "world",
		Id:       1,
		Uid:      uint(2),
		Num:      3.14,
	}

	params, err := prepareValues(profile)
	if err != nil {
		t.Fatal(err)
	}

	if params.Get("myzone") != "hello" {
		t.Errorf("myzone params wasn't properly set, got %v", params.Get("myzone"))
	}

	if params.Get("name") != "world" {
		t.Errorf("name params wasn't properly set, got %v", params.Get("name"))
	}

	if params.Get("IgnoreMe") != "" || params.Get("ignoreme") != "" {
		t.Errorf("IgnoreMe key was set")
	}
}

func TestPrepareValuesRequired(t *testing.T) {
	profile := struct {
		ReqiredField string `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}
