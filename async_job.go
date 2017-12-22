package egoscale

import (
	"encoding/json"
)

// PollAsyncJob is an alias to QueryAsyncJobResult
func (exo *Client) PollAsyncJob(profile AsyncJobResultProfile) (*QueryAsyncJobResultResponse, error) {
	return exo.QueryAsyncJobResult(profile)
}

// QueryAsyncJobResult queries the status of an async job
func (exo *Client) QueryAsyncJobResult(profile AsyncJobResultProfile) (*QueryAsyncJobResultResponse, error) {
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
