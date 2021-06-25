package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type s3LogLines []s3LogLine

type LambdaEvent struct {
	Data s3LogLine `json:"data"`
}

// s3LogLine format of S3 buckets logs
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

// Constants for indexing into a regex
const (
	BUCKETOWNER = iota + 1
	BUCKET
	TIME
	REMOTEIP
	REQUESTER
	REQUESTID
	OPERATION
	KEY
	REQUESTURI
	HTTPSTATUSCODE
	ERRORCODE
	BYTESSENT
	OBJECTSIZE
	TOTALTIME
	TURNAROUNDTIME
	REFERRER
	USERAGENT
	VERSIONID
)

// Regex match patterns for each field
const (
	RXBucketOwner    string = `^(\w*)\s`
	RXBucket         string = `([a-zA-Z0-9\-]*)\s`
	RXTime           string = `(\[.*\])\s`
	RXRemoteIP       string = `([0-9\.]*)\s`
	RXRequester      string = `(\w*)\s`
	RXRequestId      string = `(\w*)\s`
	RXOperation      string = `([a-zA-Z0-9\.]*)\s`
	RXKey            string = `([a-zA-Z0-9+\%\.-]*)\s`
	RXRequestURI     string = `"(.*?)"\s`
	RXHttpStatusCode string = `(\d*)\s`
	RXErrorCode      string = `([\w-]*)\s`
	RXBytesSent      string = `([0-9-]*)\s`
	RXObjectSize     string = `([0-9-]*)\s`
	RXTotalTime      string = `([0-9-]*)\s`
	RXTurnAroundTime string = `([0-9-]*)\s`
	RXReferrer       string = `"(.*?)"\s`
	RXUserAgent      string = `"(.*?)"\s`
	RXVersionId      string = `(.*)$`
)

// Regex for match an entire log line
const S3Regex string = RXBucketOwner + RXBucket + RXTime +
	RXRemoteIP + RXRequester + RXRequestId + RXOperation + RXKey +
	RXRequestURI + RXHttpStatusCode + RXErrorCode + RXBytesSent +
	RXObjectSize + RXTotalTime + RXTurnAroundTime + RXReferrer +
	RXUserAgent + RXVersionId

// S3Operation only supporting Rest, not including SOAP
//  Future use for parsing sub-objects
type S3Operation struct {
	API          string `json:"api"`
	HTTPMethod   string `json:"httpMethod"`
	ResourceType string `json:"resourceType"`
}

// S3RequestURI
//  Future use for parsing sub-objects
type S3RequestURI struct {
	HTTPMethod string `json:"httpMethod"`
	Path       string `json:"path"`
	Protocol   string `json:"protocol"`
}

//  Future use for filtering logic
type S3Filter struct {
	MatchMethods     []S3RequestURI `json:"matchMethods"`
	NotMatchMethods  []S3RequestURI `json:"notMatchMethods"`
	MatchProtocol    []S3RequestURI `json:"matchProtocol"`
	NotMatchProtocol []S3RequestURI `json:"notMatchProtocol"`
}

//  Future use for filtering logic
type S3PatternFilter struct {
	Match   string
	Type    string
	ApplyTo string
}

// String print a log line as a string
func (li *s3LogLine) String() string {
	return fmt.Sprintf(
		"%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v\n%v",
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

// readLines from the file specified
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

// parseS3 S3 log file
func parseS3(file string) ([]s3LogLine, error) {
	var items []s3LogLine

	lines, err := readLines(file)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// skiplines the first two lines of a Wasabi file are headers
	//   skip them not sure if that is true for S3 proper
	skiplines := 0

	// Compile regex once
	regex := *regexp.MustCompile(S3Regex)

	// Log to metrics in JOb
	fmt.Printf("Lines to process %d\n", len(lines))
	for _, line := range lines {
		if skiplines < 2 {
			skiplines += 1
			continue
		}

		match := regex.FindStringSubmatch(line)

		// If the match fails we have a log line we don't
		// know how to process check pattern
		if len(match) == 0 {
			fmt.Println(S3Regex)
			fmt.Println(line)
			os.Exit(0)
		}

		lineItem := new(s3LogLine)

		lineItem.BucketOwner = match[BUCKETOWNER]
		lineItem.Bucket = match[BUCKET]

		/*
			layout := "02/Jun/2021:01:28:20 +0000"
			t, _ := time.Parse(layout, match[TIME])
		*/

		lineItem.Time = match[TIME]

		lineItem.RemoteIP = match[REMOTEIP]
		lineItem.Requester = match[REQUESTER]
		lineItem.RequestId = match[REQUESTID]
		lineItem.Operation = match[OPERATION]
		lineItem.Key = match[KEY]
		lineItem.RequestURI = match[REQUESTURI]

		status, err := strconv.Atoi(match[HTTPSTATUSCODE])
		if err != nil {
			status = 0
		}
		lineItem.HttpStatusCode = status

		lineItem.ErrorCode = match[ERRORCODE]

		bytes, err := strconv.Atoi(match[BYTESSENT])
		if err != nil {
			status = 0
		}
		lineItem.BytesSent = bytes

		size, err := strconv.Atoi(match[OBJECTSIZE])
		if err != nil {
			status = 0
		}
		lineItem.ObjectSize = size

		tt, err := strconv.Atoi(match[TOTALTIME])
		if err != nil {
			status = 0
		}
		lineItem.TotalTime = tt

		tt, err = strconv.Atoi(match[TURNAROUNDTIME])
		if err != nil {
			status = 0
		}
		lineItem.TurnAroundTime = tt

		lineItem.Referrer = match[REFERRER]
		lineItem.UserAgent = match[USERAGENT]
		lineItem.VersionId = match[VERSIONID]

		items = append(items, *lineItem)
	}
	return items, nil
}
