package egoscale

import (
	"testing"
)

func TestQueryAsyncJobResult(t *testing.T) {
	req := &QueryAsyncJobResult{}
	_ = req.Response().(*AsyncJobResult)
}

func TestListAsyncJobs(t *testing.T) {
	req := &ListAsyncJobs{}
	_ = req.Response().(*ListAsyncJobsResponse)
}

func TestAsyncJobsResultDeepCopy(t *testing.T) {
	req := &AsyncJobResult{
		JobID: MustParseUUID("5d99e658-d0d4-4a91-b194-fc8bb44c272e"),
	}
	copy := req.DeepCopy()

	if copy.JobID == nil {
		t.Errorf("JobID non nil is expected")
	}

	if !copy.JobID.Equal(*MustParseUUID("5d99e658-d0d4-4a91-b194-fc8bb44c272e")) {
		t.Errorf(`uuid is not equal to "5d99e658-d0d4-4a91-b194-fc8bb44c272e": got %q`, copy.JobID.String())
	}
}

func TestAsyncJobsResultDeepCopyInto(t *testing.T) {
	req := &AsyncJobResult{
		JobID: MustParseUUID("5d99e658-d0d4-4a91-b194-fc8bb44c272e"),
	}
	copy := new(AsyncJobResult)
	req.DeepCopyInto(copy)

	if copy.JobID == nil {
		t.Errorf("JobID non nil is expected")
	}

	if !copy.JobID.Equal(*MustParseUUID("5d99e658-d0d4-4a91-b194-fc8bb44c272e")) {
		t.Errorf(`uuid is not equal to "5d99e658-d0d4-4a91-b194-fc8bb44c272e": got %q`, copy.JobID.String())
	}
}
