package egoscale

import (
	"encoding/json"
)

const (
	// PENDING represents a job in progress
	PENDING JobStatusType = iota
	// SUCCESS represents a successfully completed job
	SUCCESS
	// FAILURE represents a job that has failed to complete
	FAILURE
)

// JobStatusType represents the status of a Job
type JobStatusType int

// QueryASyncJobResultResponse represents the current status of an asynchronous job
type QueryAsyncJobResultResponse struct {
	Accountid       string           `json:"accountid,omitempty"`
	Cmd             string           `json:"cmd,omitempty"`
	Created         string           `json:"created,omitempty"`
	JobInstanceId   string           `json:"jobinstanceid,omitempty"`
	JobInstanceType string           `json:"jobinstancetype,omitempty"`
	JobProcStatus   int              `json:"jobprocstatus,omitempty"`
	JobResult       *json.RawMessage `json:"jobresult,omitempty"`
	JobResultCode   int              `json:"jobresultcode,omitempty"`
	JobResultType   string           `json:"jobresulttype,omitempty"`
	JobStatus       JobStatusType    `json:"jobstatus,omitempty"`
	Userid          string           `json:"userid,omitempty"`
}

// JobResultResponse represents a generic response to a job task
type JobResultResponse struct {
	AccountId     string           `json:"accountid,omitempty"`
	Cmd           string           `json:"cmd,omitempty"`
	CreatedAt     string           `json:"created,omitempty"`
	JobId         string           `json:"jobid,omitempty"`
	JobProcStatus int              `json:"jobprocstatus,omitempty"`
	JobResult     *json.RawMessage `json:"jobresult,omitempty"`
	JobStatus     JobStatusType    `json:"jobstatus,omitempty"`
	JobResultType string           `json:"jobresulttype,omitempty"`
	UserId        string           `json:"userid,omitempty"`
}

// AssociateIpAddressResponse represents the response to the creation of an IpAddress
type AssociateIpAddressResponse struct {
	IpAddress IpAddress `json:"ipaddress"`
}

// DisassociateIpAddressResponse represents the response to the deletion of an IpAddress
type DisassociateIpAddressResponse struct {
	Success     bool   `json:"success"`
	DisplayText string `json:"diplaytext,omitempty"`
}
