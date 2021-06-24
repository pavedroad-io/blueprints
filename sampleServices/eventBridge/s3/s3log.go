package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type s3LogLines []s3LogLine

type LambdaEvent struct {
	Data s3LogLine `json:"data"`
}

type s3LogLine struct {
	BucketOwner    string `json:"bucketOwner"`    //0
	Bucket         string `json:"bucket"`         //1
	Time           string `json:"time"`           //2
	RemoteIP       string `json:"remoteIP"`       //3
	Requester      string `json:"requester"`      //4
	RequestId      string `json:"requestId"`      //5
	Operation      string `json:"operation"`      //6
	Key            string `json:"key"`            //7
	RequestURI     string `json:"requestURI"`     //8
	HttpStatusCode int    `json:"httpStatusCode"` //9
	ErrorCode      string `json:"errorCode"`      //10
	BytesSent      int    `json:"bytesSent"`      //11
	ObjectSize     int    `json:"objectSize"`     //12
	TotalTime      int    `json:"totalTime"`      //13
	TurnAroundTime int    `json:"turnAroundTime"` //14
	Referrer       string `json:"referrer"`       //15
	UserAgent      string `json:"userAgent"`      //16
	VersionId      string `json:"versionId"`      //17
}

// S3Operation only supporting Rest, not including SOAP
type S3Operation struct {
	API          string `json:"api"`
	HTTPMethod   string `json:"httpMethod"`
	ResourceType string `json:"resourceType"`
}

// S3RequestURI
type S3RequestURI struct {
	HTTPMethod string `json:"httpMethod"`
	Path       string `json:"path"`
	Protocol   string `json:"protocol"`
}

type S3Filter struct {
	MatchMethods     []S3RequestURI `json:"matchMethods"`
	NotMatchMethods  []S3RequestURI `json:"notMatchMethods"`
	MatchProtocol    []S3RequestURI `json:"matchProtocol"`
	NotMatchProtocol []S3RequestURI `json:"notMatchProtocol"`
}

type S3PatternFilter struct {
	Match   string
	Type    string
	ApplyTo string
}

func (li *s3LogLine) String() string {
	return fmt.Sprintf(
		"%s\t%s\t%d\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%d\t%d\t%d\t%d\t%d\t%s\t%s\t%d",
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
		li.ErrorCode,
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

		fields := strings.Split(line, " ")
		//		fmt.Println(fields)
		//		os.Exit(0)

		lineItem := new(s3LogLine)

		lineItem.BucketOwner = fields[0]
		lineItem.Bucket = fields[1]

		/*
			layout := "02/Jun/2021:01:28:20 +0000"
			t, _ := time.Parse(layout, fields[2])
		*/
		lineItem.Time = fields[2] + fields[3]

		lineItem.RemoteIP = fields[4]
		lineItem.Requester = fields[5]
		lineItem.RequestId = fields[6]
		lineItem.Operation = fields[7]
		lineItem.Key = fields[8]
		lineItem.RequestURI = fields[9]

		status, err := strconv.Atoi(fields[10])
		if err != nil {
			status = 0
		}
		lineItem.HttpStatusCode = status

		lineItem.ErrorCode = fields[11]

		bytes, err := strconv.Atoi(fields[12])
		if err != nil {
			status = 0
		}
		lineItem.BytesSent = bytes

		size, err := strconv.Atoi(fields[13])
		if err != nil {
			status = 0
		}
		lineItem.ObjectSize = size

		tt, err := strconv.Atoi(fields[14])
		if err != nil {
			status = 0
		}
		lineItem.TotalTime = tt

		tt, err = strconv.Atoi(fields[15])
		if err != nil {
			status = 0
		}
		lineItem.TurnAroundTime = tt

		lineItem.Referrer = fields[16]
		lineItem.UserAgent = fields[17]
		lineItem.VersionId = fields[18]

		items = append(items, *lineItem)
	}
	return items, nil
}
