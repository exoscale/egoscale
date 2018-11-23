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
