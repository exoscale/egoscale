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
		Bytes    []byte  `json:"bytes"`
	}{
		IgnoreMe: "bar",
		Name:     "world",
		Id:       1,
		Uid:      uint(2),
		Num:      3.14,
		Bytes:    []byte("exo"),
	}

	params, err := prepareValues(profile)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := (*params)["myzone"]; ok {
		t.Errorf("myzone params shouldn't be set, got %v", params.Get("myzone"))
	}

	if params.Get("name") != "world" {
		t.Errorf("name params wasn't properly set, got %v", params.Get("name"))
	}

	if params.Get("bytes") != "ZXhv" {
		t.Errorf("bytes params wasn't properly encoded in base 64, got %v", params.Get("bytes"))
	}

	if _, ok := (*params)["ignoreme"]; ok {
		t.Errorf("IgnoreMe key was set")
	}
}

func TestPrepareValuesStringRequired(t *testing.T) {
	profile := struct {
		RequiredField string `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}

func TestPrepareValuesBoolRequired(t *testing.T) {
	profile := struct {
		RequiredField bool `json:"requiredfield"`
	}{}

	params, err := prepareValues(&profile)
	if err != nil {
		t.Fatal(nil)
	}
	if params.Get("requiredfield") != "false" {
		t.Errorf("bool params wasn't set to false (default value)")
	}
}

func TestPrepareValuesIntRequired(t *testing.T) {
	profile := struct {
		RequiredField int64 `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}

func TestPrepareValuesUintRequired(t *testing.T) {
	profile := struct {
		RequiredField uint64 `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}

func TestPrepareValuesBytesRequired(t *testing.T) {
	profile := struct {
		RequiredField []byte `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}

func TestPrepareValuesSliceString(t *testing.T) {
	profile := struct {
		RequiredField []string `json:"requiredfield"`
	}{}

	_, err := prepareValues(&profile)
	if err == nil {
		t.Errorf("It should have failed")
	}
}
