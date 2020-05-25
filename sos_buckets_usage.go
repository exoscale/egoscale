package egoscale

// BucketUsage represents the usage (in byte) for a bucket
type BucketUsage struct {
	Created string `json:"created"`
	Name    string `json:"name"`
	Region  string `json:"region"`
	Usage   int    `json:"usage"`
}

// ListBucketsUsage represents a list buckets usage API request
type ListBucketsUsage struct {
	_ bool `name:"listBucketsUsage" description:"List"`
}

// ListBucketsUsageResponse represents a list buckets usage API response
type ListBucketsUsageResponse struct {
	Count        int           `json:"count"`
	BucketsUsage []BucketUsage `json:"bucketsusage"`
}

// Response returns the struct to unmarshal
func (ListBucketsUsage) Response() interface{} {
	return new(ListBucketsUsageResponse)
}
