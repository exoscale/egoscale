package egoscale

import (
	"testing"
)

func TestSnapshot(t *testing.T) {
	instance := &Snapshot{}
	if instance.ResourceType() != "Snapshot" {
		t.Errorf("ResourceType doesn't match")
	}
}

func TestCreateSnapshot(t *testing.T) {
	req := &CreateSnapshot{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*Snapshot)
}

func TestListSnapshots(t *testing.T) {
	req := &ListSnapshots{}
	_ = req.response().(*ListSnapshotsResponse)
}

func TestDeleteSnapshot(t *testing.T) {
	req := &DeleteSnapshot{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestRevertSnapshot(t *testing.T) {
	req := &RevertSnapshot{}
	_ = req.response().(*AsyncJobResult)
	_ = req.asyncResponse().(*booleanResponse)
}

func TestListSnapshotEmpty(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listsnapshotsresponse": {
	"count": 0,
	"snapshot": []
  }}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	snapshot := new(Snapshot)
	snapshots, err := cs.List(snapshot)
	if err != nil {
		t.Fatal(err)
	}

	if len(snapshots) != 0 {
		t.Errorf("zero networks were expected, got %d", len(snapshots))
	}
}

func TestListSnapshotFailure(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listsnapshotsresponse": {
	"count": 2345,
	"snapshot": {}
  }}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	snapshot := new(Snapshot)
	_, err := cs.List(snapshot)
	if err == nil {
		t.Errorf("error was expected, got %v", err)
	}
}

func TestFindSnapshoy(t *testing.T) {
	ts := newServer(response{200, jsonContentType, `
{"listsnapshotsresponse": {
	"count": 1,
	"snapshot": [
	  {
			"account": "exoscale-1",
      "accountid": "307dd827-aa02-42fd-9b43-e7de55c6e7c5",
      "created": "2018-11-08T10:52:39+0100",
      "domain": "exoscale-1",
      "domainid": "5b2f621e-3eb6-4a14-a315-d4d7d62f28ff",
      "id": "e09c6194-8976-493a-a57c-99650d51da9b",
      "intervaltype": "MANUAL",
      "name": "rt_ROOT-690946_20181108095239",
      "revertable": true,
      "size": 53687091200,
      "snapshottype": "MANUAL",
      "state": "BackedUp",
      "tags": [],
      "volumeid": "74073768-9502-49b4-93fb-f11c5b106976",
      "volumename": "ROOT-690946",
      "volumetype": "ROOT"
	  }
	]
  }}`})
	defer ts.Close()

	cs := NewClient(ts.URL, "KEY", "SECRET")

	snapshots, err := cs.List(&Snapshot{Name: "rt_ROOT-690946_20181108095239"})
	if err != nil {
		t.Fatal(err)
	}

	if len(snapshots) != 1 {
		t.Fatalf("One snapshot was expected, got %d", len(snapshots))
	}

	net, ok := snapshots[0].(*Snapshot)
	if !ok {
		t.Errorf("unable to type inference *Snapshot, got %v", net)
	}

	if snapshots[0].(*Snapshot).Name != "rt_ROOT-690946_20181108095239" {
		t.Errorf("rt_ROOT-690946_20181108095239 snapshot name was expected, got %s", snapshots[0].(*Network).Name)
	}

}
