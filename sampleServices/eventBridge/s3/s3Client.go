package main

import (
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func newClient(p Provider) (*minio.Client, error) {
	s3Client, err := minio.New(p.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(p.Credentials, p.Key, p.Region),
		Secure: true,
	})
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return s3Client, nil
}
