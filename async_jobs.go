package egoscale

import (
	"encoding/json"
)

// QueryAsyncJobResultRequest represents a query to fetch the status of async job
type QueryAsyncJobResultRequest struct {
	JobID string `json:"jobid"`
}

// Command returns the CloudStack API command
func (req *QueryAsyncJobResultRequest) Command() string {
	return "queryAsyncJobResult"
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
