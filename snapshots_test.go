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
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*Snapshot)
}

func TestListSnapshots(t *testing.T) {
	req := &ListSnapshots{}
	_ = req.Response().(*ListSnapshotsResponse)
}

func TestDeleteSnapshot(t *testing.T) {
	req := &DeleteSnapshot{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestRevertSnapshot(t *testing.T) {
	req := &RevertSnapshot{}
	_ = req.Response().(*AsyncJobResult)
	_ = req.AsyncResponse().(*BooleanResponse)
}

func TestExportSnapshot(t *testing.T) {
	req := &ExportSnapshot{}
	_ = req.Response().(*ExportSnapshotResponse)
}
