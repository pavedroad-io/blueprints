package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type s3LogLines []s3LogLine

type LambdaEvent struct {
	Data s3LogLine `json:"data"`
}

type s3LogLine struct {
	BucketOwner    string    `yaml:"bucketOwner"`
	Bucket         string    `yaml:"bucket"`
	Time           time.Time `yaml:"time"`
	RemoteIP       string    `yaml:"remoteIP"`
	Requester      string    `yaml:"requester"`
	RequestId      string    `yaml:"requestId"`
	Operation      string    `yaml:"operation"`
	Key            string    `yaml:"key"`
	RequestURI     string    `yaml:"requestURI"`
	HttpStatusCode int       `yaml:"httpStatusCode"`
	BytesSent      int       `yaml:"bytesSent"`
	ObjectSize     int       `yaml:"objectSize"`
	TotalTime      int       `yaml:"totalTime"`
	TurnAroundTime int       `yaml:"turnAroundTime"`
	Referrer       string    `yaml:"referrer"`
	UserAgent      string    `yaml:"userAgent"`
	VersionId      int       `yaml:"versionId"`
}

// S3Operation only supporting Rest, not including SOAP
type S3Operation struct {
	API          string `yaml:"api"`
	HTTPMethod   string `yaml:"httpMethod"`
	ResourceType string `yaml:"resourceType"`
}

// S3RequestURI
type S3RequestURI struct {
	HTTPMethod string `yaml:"httpMethod"`
	Path       string `yaml:"path"`
	Protocol   string `yaml:"protocol"`
}

type S3Filter struct {
	MatchMethods     []S3RequestURI `yaml:"matchMethods"`
	NotMatchMethods  []S3RequestURI `yaml:"notMatchMethods"`
	MatchProtocol    []S3RequestURI `yaml:"matchProtocol"`
	NotMatchProtocol []S3RequestURI `yaml:"notMatchProtocol"`
}

type S3PatternFilter struct {
	Match   string
	Type    string
	ApplyTo string
}

func (li *s3LogLine) String() string {
	return fmt.Sprintf(
		"%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%d",
		li.BucketOwner,
		li.Bucket,
		li.Time,
		li.RemoteIP,
		li.Requester,
		li.RequestId,
		li.Operation,
		li.Key,
		li.RequestURI,
		li.HttpStatusCode,
		li.BytesSent,
		li.ObjectSize,
		li.TotalTime,
		li.TurnAroundTime,
		li.Referrer,
		li.UserAgent,
		li.VersionId,
	)
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func parseS3(file string) ([]s3LogLine, error) {
	var items []s3LogLine

	lines, err := readLines(file)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// skiplines the first two lines of a Wasabi file are headers
	//   skip them
	skiplines := 0
	fmt.Printf("Lines to process %d\n", len(lines))
	for _, line := range lines {
		if skiplines < 2 {
			skiplines += 1
			continue
		}

		fields := strings.Fields(line)
		//		fmt.Println(fields)
		//		os.Exit(0)

		lineItem := new(s3LogLine)

		lineItem.BucketOwner = fields[0]
		lineItem.Bucket = fields[1]

		layout := "02/Jan/2006:15:04:05 -0700"
		t, _ := time.Parse(layout, fields[2])
		lineItem.Time = t

		lineItem.RemoteIP = fields[3]
		lineItem.Requester = fields[4]
		lineItem.RequestId = fields[5]
		lineItem.Operation = fields[6]
		lineItem.Key = fields[7]
		lineItem.RequestURI = fields[8]

		status, err := strconv.Atoi(fields[9])
		if err != nil {
			status = 0
		}
		lineItem.HttpStatusCode = status

		bytes, err := strconv.Atoi(fields[10])
		if err != nil {
			status = 0
		}
		lineItem.BytesSent = bytes

		size, err := strconv.Atoi(fields[11])
		if err != nil {
			status = 0
		}
		lineItem.ObjectSize = size

		tt, err := strconv.Atoi(fields[12])
		if err != nil {
			status = 0
		}
		lineItem.TotalTime = tt

		tt, err = strconv.Atoi(fields[13])
		if err != nil {
			status = 0
		}
		lineItem.TurnAroundTime = tt

		lineItem.Referrer = fields[14]
		lineItem.UserAgent = fields[14]

		version, err := strconv.Atoi(fields[15])
		if err != nil {
			status = 0
		}
		lineItem.VersionId = version

		items = append(items, *lineItem)
	}
	return items, nil
}
