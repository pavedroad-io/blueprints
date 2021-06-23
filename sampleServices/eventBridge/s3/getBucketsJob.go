package main

import (
	"context"

	"github.com/minio/minio-go/v7"
)

func listBuckets(c *minio.Client) ([]minio.BucketInfo, error) {

	buckets, err := c.ListBuckets(context.Background())
	if err != nil {
		return nil, err
	}
	return buckets, nil
}
