package storage

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/exoscale/egoscale/api"
)

// Bucket represents a Storage bucket.
type Bucket struct {
	api.Resource

	Name string

	c *Client
}

// TODO: Bucket.SetACL()

// TODO: Bucket.ACL()

// TODO: Bucket.SetCORS()

// TODO: Bucket.CORS()

// TODO: Bucket.Zone()

// TODO: Bucket.Files()

// TODO: Bucket.File()

// TODO: Bucket.PutFile()

// TODO: Bucket.DeleteFiles()

// TODO: Bucket.Delete()

// TODO: CreateBucket()

// ListBuckets returns the list of Storage buckets owned, or an error if the API query failed.
func (c *Client) ListBuckets() ([]*Bucket, error) {
	buckets := make([]*Bucket, 0)

	res, err := c.c.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		return nil, err
	}

	for _, bucket := range res.Buckets {
		buckets = append(buckets, c.bucketFromAPI(bucket))
	}

	return buckets, nil
}

// TODO: GetBucket()

func (c *Client) bucketFromAPI(bucket *s3.Bucket) *Bucket {
	return &Bucket{
		Resource: api.MarshalResource(bucket),
		Name:     aws.StringValue(bucket.Name),
		c:        c,
	}
}
