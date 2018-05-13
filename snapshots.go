package egoscale

// ResourceType returns the type (for tags) of the resource
func (*Snapshot) ResourceType() string {
	return "Snapshot"
}

func (*CreateSnapshot) name() string {
	return "createSnapshot"
}

func (*CreateSnapshot) asyncResponse() interface{} {
	return new(CreateSnapshotResponse)
}

func (*ListSnapshots) name() string {
	return "listSnapshots"
}

func (*ListSnapshots) response() interface{} {
	return new(ListSnapshotsResponse)
}

func (*DeleteSnapshot) name() string {
	return "deleteSnapshot"
}

func (*DeleteSnapshot) asyncResponse() interface{} {
	return new(booleanResponse)
}

func (*RevertSnapshot) name() string {
	return "revertSnapshot"
}

func (*RevertSnapshot) asyncResponse() interface{} {
	return new(booleanResponse)
}

func (*ListSnapshotPolicies) name() string {
	return "listSnapshotPolicies"
}

func (*ListSnapshotPolicies) response() interface{} {
	return new(ListSnapshotPoliciesResponse)
}
