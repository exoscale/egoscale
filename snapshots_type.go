package egoscale

// SnapshotState represents the Snapshot.State enum
//
// See: https://github.com/apache/cloudstack/blob/master/api/src/main/java/com/cloud/storage/Snapshot.java
type SnapshotState int

//go:generate stringer -type SnapshotState
const (
	// Allocated ... (TODO)
	Allocated SnapshotState = iota
	// Creating ... (TODO)
	Creating
	// CreatedOnPrimary ... (TODO)
	CreatedOnPrimary
	// BackingUp ... (TODO)
	BackingUp
	// BackedUp ... (TODO)
	BackedUp
	// Copying ... (TODO)
	Copying
	// Destroying ... (TODO)
	Destroying
	// Destroyed ... (TODO)
	Destroyed
	// Error ... (TODO)
	Error
)

// Snapshot represents a volume snapshot
type Snapshot struct {
	Account      string        `json:"account,omitempty" doc:"the account associated with the snapshot"`
	Created      string        `json:"created,omitempty" doc:"  the date the snapshot was created"`
	Domain       string        `json:"domain,omitempty" doc:"the domain name of the snapshot's account"`
	DomainID     string        `json:"domainid,omitempty" doc:"the domain ID of the snapshot's account"`
	ID           string        `json:"id,omitempty" doc:"ID of the snapshot"`
	IntervalType string        `json:"intervaltype,omitempty" doc:"valid types are hourly, daily, weekly, monthy, template, and none."`
	Name         string        `json:"name,omitempty" doc:"name of the snapshot"`
	PhysicalSize int64         `json:"physicalsize,omitempty" doc:"physical size of the snapshot on image store"`
	Project      string        `json:"project,omitempty" doc:"the project name of the snapshot"`
	ProjectID    string        `json:"projectid,omitempty" doc:"the project id of the snapshot"`
	Revertable   *bool         `json:"revertable,omitempty" doc:"indicates whether the underlying storage supports reverting the volume to this snapshot"`
	Size         int64         `json:"size,omitempty" doc:"the size of original volume"`
	SnapshotType string        `json:"snapshottype,omitempty" doc:"the type of the snapshot"`
	State        SnapshotState `json:"state,omitempty" doc:"the state of the snapshot. BackedUp means that snapshot is ready to be used; Creating - the snapshot is being allocated on the primary storage; BackingUp - the snapshot is being backed up on secondary storage"`
	Tags         []ResourceTag `json:"tags,omitempty" doc:"the list of resource tags associated with snapshot"`
	VolumeID     string        `json:"volumeid,omitempty" doc:"ID of the disk volume"`
	VolumeName   string        `json:"volumename,omitempty" doc:"name of the disk volume"`
	VolumeType   string        `json:"volumetype,omitempty" doc:"type of the disk volume"`
	ZoneID       string        `json:"zoneid,omitempty" doc:"id of the availability zone"`
}

// SnapshotPolicy represents a snapshot policy
type SnapshotPolicy struct {
	ForDisplay   *bool  `json:"fordisplay,omitempty" doc:"is this policy for display to the regular user"`
	ID           string `json:"id,omitempty" doc:"the ID of the snapshot policy"`
	IntervalType int16  `json:"intervaltype,omitempty" doc:"the interval type of the snapshot policy"`
	MaxSnaps     int    `json:"maxsnaps,omitempty" doc:"maximum number of snapshots retained"`
	Schedule     string `json:"schedule,omitempty" doc:"time the snapshot is scheduled to be taken."`
	Timezone     string `json:"timezone,omitempty" doc:"the time zone of the snapshot policy"`
	VolumeID     string `json:"volumeid,omitempty" doc:"the ID of the disk volume"`
}

// CreateSnapshot (Async) creates an instant snapshot of a volume
//
// CloudStackAPI: http://cloudstack.apache.org/api/apidocs-4.10/apis/createSnapshot.html
type CreateSnapshot struct {
	VolumeID  string `json:"volumeid" doc:"The ID of the disk volume"`
	Account   string `json:"account,omitempty" doc:"The account of the snapshot. The account parameter must be used with the domainId parameter."`
	Name      string `json:"name,omitempty" doc:"the name of the snapshot"`
	DomainID  string `json:"domainid,omitempty" doc:"The domain ID of the snapshot. If used with the account parameter, specifies a domain for the account associated with the disk volume."`
	PolicyID  string `json:"policyid,omitempty" doc:"policy id of the snapshot, if this is null, then use MANUAL_POLICY."`
	QuiesceVM *bool  `json:"quiescevm,omitempty" doc:"quiesce vm if true"`
}

// CreateSnapshotResponse represents a freshly created snapshot
type CreateSnapshotResponse struct {
	Snapshot Snapshot `json:"snapshot"`
}

// ListSnapshots lists the volume snapshots
//
// CloudStackAPI: http://cloudstack.apache.org/api/apidocs-4.10/apis/listSnapshots.html
type ListSnapshots struct {
	Account      string        `json:"account,omitempty" doc:"list resources by account. Must be used with the domainId parameter."`
	DomainID     string        `json:"domainid,omitempty" doc:"list only resources belonging to the domain specified"`
	ID           string        `json:"id,omitempty" doc:"lists snapshot by snapshot ID"`
	IDs          []string      `json:"ids,omitempty" doc:"the IDs of the snapshots, mutually exclusive with id"`
	IntervalType string        `json:"intervaltype,omitempty" doc:"valid values are HOURLY, DAILY, WEEKLY, and MONTHLY."`
	IsRecursive  *bool         `json:"isrecursive,omitempty" doc:"defaults to false, but if true, lists all resources from the parent specified by the domainId till leaves."`
	Keyword      string        `json:"keyword,omitempty" doc:"List by keyword"`
	ListAll      *bool         `json:"listall,omitempty" doc:"If set to false, list only resources belonging to the command's caller; if set to true - list resources that the caller is authorized to see. Default value is false"`
	Name         string        `json:"name,omitempty" doc:"lists snapshot by snapshot name"`
	Page         int           `json:"page,omitempty"`
	PageSize     int           `json:"pagesize,omitempty"`
	ProjectID    string        `json:"projectid,omitempty" doc:"list objects by project"`
	SnapshotType string        `json:"snapshottype,omitempty" doc:"valid values are MANUAL or RECURRING."`
	Tags         []ResourceTag `json:"tags,omitempty" doc:"List resources by tags (key/value pairs)"`
	VolumeID     string        `json:"volumeid,omitempty" doc:"the ID of the disk volume"`
	ZoneID       string        `json:"zoneid,omitempty" doc:"list snapshots by zone id"`
}

// ListSnapshotsResponse represents a list of volume snapshots
type ListSnapshotsResponse struct {
	Count    int        `json:"count"`
	Snapshot []Snapshot `json:"snapshot"`
}

// DeleteSnapshot (Async) deletes a snapshot of a disk volume
//
// CloudStackAPI: http://cloudstack.apache.org/api/apidocs-4.10/apis/deleteSnapshot.html
type DeleteSnapshot struct {
	ID string `json:"id" doc:"The ID of the snapshot"`
}

// RevertSnapshot (Async) reverts a volume snapshot
//
// CloudStackAPI: http://cloudstack.apache.org/api/apidocs-4.10/apis/revertSnapshot.html
type RevertSnapshot struct {
	ID string `json:"id" doc:"The ID of the snapshot"`
}

// ListSnapshotPolicies lists snapshot policies
//
// CloudStack API: https://cloudstack.apache.org/api/apidocs-4.4/user/listSnapshotPolicies.html
type ListSnapshotPolicies struct {
	ForDisplay *bool  `json:"fordisplay,omitempty" doc:"list resources by display flag; only ROOT admin is eligible to pass this parameter"`
	ID         string `json:"id,omitempty" doc:"the ID of the snapshot policy"`
	Keyword    string `json:"keyword,omitempty" doc:"List by keyword"`
	Page       int    `json:"page,omitempty"`
	PageSize   int    `json:"pagesize,omitempty"`
	VolumeID   string `json:"volumeid,omitempty" doc:"the ID of the disk volume"`
}

// ListSnapshotPoliciesResponse represents a list of snapshot policies
type ListSnapshotPoliciesResponse struct {
	Count          int              `json:"count"`
	SnapshotPolicy []SnapshotPolicy `json:"snapshotpolicy"`
}
