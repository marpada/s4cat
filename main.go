package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
	"log"
	"os"
)

func main() {
	var (
		output, region string
	)
	flag.StringVar(&output, "output", "/dev/stdout", "Place to send the output")
	flag.StringVar(&region, "region", os.Getenv("AWS_DEFAULT_REGION"), "Region where the bucket is located")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-output=FILEPATH] [-region=REGION] <bucket> <key>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	bucket, key := flag.Arg(0), flag.Arg(1)

	s := s3.New(&aws.Config{Region: region})
	resp, err := s.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		log.Fatalf("GetObject Failed: %#+v", err)
	}
	fd, err := os.Create(output)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()
	n, err := io.Copy(fd, resp.Body)
	if err != nil {
		log.Fatal("Copy failed after", n, "bytes:", err)
	}
}
