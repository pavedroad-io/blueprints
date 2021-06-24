package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

func main() {
	c := Customer{}
	customers, err := c.LoadFromDisk("customer.yaml")
	if err != nil {
		log.Fatalf("fail loading customer.yaml: %v\n", err)

	}
	opts := minio.ListObjectsOptions{
		Recursive: true,
		Prefix:    "",
	}

	// This will be a Job
	var logQueue []LogQueueItem
	var plogs ProcessedLogs
	for _, c := range customers {
		// Load a list of previously processed logs
		// For now ignore error if not found
		plogs.LoadFromDisk(c.ID.String())

		// Build a list of providers the customer
		// uses
		plist := c.Providers

		for i, l := range c.Logs {
			p, err := plist.Lookup(l.Provider)
			if err != nil {
				log.Printf("Provider not found: %v\n", err)
			}
			s3Client, err := newClient(p)
			if err != nil {
				log.Fatalln(err)
			}
			objects, err := listBucketObjects(s3Client, l.Name, opts)
			if err != nil {
				log.Fatalln(err)
			}

			for _, o := range objects {

				// fmt.Println(o.Key)
				f, err := getObject(s3Client, l.Name, o.Key, minio.GetObjectOptions{})
				if err != nil {
					log.Fatalln(err)
				}

				if plogs.Processed(l.Name, o.Key) {
					//					fmt.Printf("Skipping %s bucket %s logs\n", l.Name, o.Key)
					continue
				}

				item := LogQueueItem{
					ID:        c.ID.String(),
					Bucket:    l.Name,
					Name:      o.Key,
					Created:   time.Now(),
					Location:  f,
					LogFormat: c.Logs[i].LogFormat,
					Processed: false,
					Prune:     c.Logs[i].PruneAfterProcessing,
				}
				logQueue = append(logQueue, item)
			}
		}
	}

	for _, l := range logQueue {
		switch l.LogFormat {
		case S3:
			po, err := parseS3(l.Location)
			if err != nil {
				fmt.Printf("Parse failed with error: %w\n", err)
			}
			for _, eventData := range po {
				le := LambdaEvent{
					Data: eventData,
				}
				j, _ := json.Marshal(le)
				//_, _ = json.Marshal(le)

				// fmt.Println(string(j))
				/*
					postBody := bytes.NewBuffer(j)

					resp, err := http.Post("http://localhost:12001/eventbridge", "application/json", postBody)

					if err != nil {
						log.Printf("HTTP POST failed error %w\n", err)
					}
					if resp.StatusCode != 200 {
						log.Printf("HTTP POST failed non 200 status code %v\n", resp.StatusCode)
					}
				*/

				fmt.Println(string(j))
			}

			l.Processed = true
			if l.Prune {
				if err := os.Remove(l.Location); err != nil {
					log.Printf("Failed to prune %s error %w\n", l.Location, err)
				}

			}

			nid, err := uuid.Parse(l.ID)
			if err != nil {
				fmt.Printf("Fail converting ID %s to UUID err %w\n", l.ID, err)
			}
			i := ProcessedLogItem{
				Date:     time.Now(),
				Bucket:   l.Bucket,
				Name:     l.Name,
				FileName: l.Location,
				Pruned:   l.Prune,
			}
			plogs.ID = nid
			plogs.AddProcessLog(l.ID, i)
		}
	}
}
