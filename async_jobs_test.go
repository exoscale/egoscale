package egoscale

import (
	"testing"
)

func TestAsyncJobs(t *testing.T) {
	var _ Command = (*QueryAsyncJobResult)(nil)
	var _ Command = (*ListAsyncJobs)(nil)
}
