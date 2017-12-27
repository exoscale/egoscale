package egoscale

import (
	"encoding/json"
)

// QueryAsyncJobResultRequest represents a query to fetch the status of async job
type QueryAsyncJobResultRequest struct {
	JobId string `json:"jobid"`
}

// QueryASyncJobResultResponse represents the current status of an asynchronous job
type QueryAsyncJobResultResponse struct {
	AccountId       string           `json:"accountid"`
	Cmd             string           `json:"cmd"`
	Created         string           `json:"created"`
	JobInstanceId   string           `json:"jobinstanceid"`
	JobInstanceType string           `json:"jobinstancetype"`
	JobProcStatus   int              `json:"jobprocstatus"`
	JobResult       *json.RawMessage `json:"jobresult"`
	JobResultCode   int              `json:"jobresultcode"`
	JobResultType   string           `json:"jobresulttype"`
	JobStatus       JobStatusType    `json:"jobstatus"`
	UserId          string           `json:"userid"`
	JobId           string           `json:"jobid"`
}

// PollAsyncJob is an alias to QueryAsyncJobResult
func (exo *Client) PollAsyncJob(profile QueryAsyncJobResultRequest) (*QueryAsyncJobResultResponse, error) {
	return exo.QueryAsyncJobResult(profile)
}

// QueryAsyncJobResult queries the status of an async job
func (exo *Client) QueryAsyncJobResult(profile QueryAsyncJobResultRequest) (*QueryAsyncJobResultResponse, error) {
	params, err := prepareValues(profile)
	if err != nil {
		return nil, err
	}

	resp, err := exo.request("queryAsyncJobResult", *params)
	if err != nil {
		return nil, err
	}

	var r QueryAsyncJobResultResponse
	if err := json.Unmarshal(resp, &r); err != nil {
		return nil, err
	}

	return &r, nil
}
