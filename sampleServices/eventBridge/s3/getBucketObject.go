package main

import (
	"context"
	"io"
	"io/ioutil"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
)

func getObject(client *minio.Client, bucket string, object string, opts minio.GetObjectOptions) (file string, er error) {

	tmpfile, err := ioutil.TempFile("/tmp/", bucket+"-"+object+"-")
	if err != nil {
		log.Fatal(err)
	}

	reader, err := client.GetObject(context.Background(), bucket, object, minio.GetObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	defer reader.Close()

	localFile, err := os.Create(tmpfile.Name())
	if err != nil {
		log.Fatalln(err)
	}
	defer localFile.Close()

	stat, err := reader.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := io.CopyN(localFile, reader, stat.Size); err != nil {
		log.Fatalln(err)
	}

	return tmpfile.Name(), nil
}
