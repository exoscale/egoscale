package egoscale

import (
	"encoding/json"
)

// QueryAsyncJobResult represents a query to fetch the status of async job
type QueryAsyncJobResult struct {
	JobID string `json:"jobid"`
}

func (req *QueryAsyncJobResult) name() string {
	return "queryAsyncJobResult"
}

func (req *QueryAsyncJobResult) response() interface{} {
	return new(QueryAsyncJobResultResponse)
}

// QueryAsyncJobResultResponse represents the current status of an asynchronous job
type QueryAsyncJobResultResponse struct {
	AccountID       string           `json:"accountid"`
	Cmd             string           `json:"cmd"`
	Created         string           `json:"created"`
	JobInstanceID   string           `json:"jobinstanceid"`
	JobInstanceType string           `json:"jobinstancetype"`
	JobProcStatus   int              `json:"jobprocstatus"`
	JobResult       *json.RawMessage `json:"jobresult"`
	JobResultCode   int              `json:"jobresultcode"`
	JobResultType   string           `json:"jobresulttype"`
	JobStatus       JobStatusType    `json:"jobstatus"`
	UserID          string           `json:"userid"`
	JobID           string           `json:"jobid"`
}
