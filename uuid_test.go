package egoscale

import (
	"encoding/json"
	"testing"
)

func recoverFromPanicing(t *testing.T) {
	if r := recover(); r != nil {
		_, ok := r.(error)
		if !ok {
			t.Error(r)
		}
	}
}

func TestUUIDDeepCopy(t *testing.T) {
	u := MustParseUUID("5361a11b-615c-42bf-9bdb-e2c3790ada14")
	v := u.DeepCopy()

	if !u.Equal(*v) {
		t.Errorf("copies should be identical")
	}

	if u == v {
		t.Errorf("copies shouldn't be the same")
	}
}

func TestUUIDDeepCopyInto(t *testing.T) {
	u := MustParseUUID("5361a11b-615c-42bf-9bdb-e2c3790ada14")
	v := MustParseUUID("c39c7da0-b879-4a71-a21c-d78b83682e65")

	if u.Equal(*v) {
		t.Errorf("u and v should be different")
	}

	u.DeepCopyInto(v)

	if !u.Equal(*v) {
		t.Errorf("copies should be identical")
	}

	if u == v {
		t.Errorf("copies shouldn't be the same")
	}
}

func TestUUIDMustParse(t *testing.T) {
	defer recoverFromPanicing(t)
	MustParseUUID("foo")
	t.Error("invalid uuid should panic")
}

func TestUUIDMarshalJSON(t *testing.T) {
	zone := &Zone{
		ID: MustParseUUID("5361a11b-615c-42bf-9bdb-e2c3790ada14"),
	}
	j, err := json.Marshal(zone)
	if err != nil {
		t.Fatal(err)
	}
	s := string(j)
	expected := `{"id":"5361a11b-615c-42bf-9bdb-e2c3790ada14"}`
	if expected != s {
		t.Errorf("bad json serialization, got %q, expected %s", s, expected)
	}
}

func TestUUIDUnmarshalJSONFailure(t *testing.T) {
	ss := []string{
		`{"id": 1}`,
		`{"id": "1"}`,
		`{"id": "5361a11b-615c-42bf-9bdb-e2c3790ada1"}`,
	}
	zone := &Zone{}
	for _, s := range ss {
		if err := json.Unmarshal([]byte(s), zone); err == nil {
			t.Errorf("an error was expected, %#v", zone)
		}
	}
}
