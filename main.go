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
		output string
	)
	flag.StringVar(&output, "output", "/dev/stdout", "Place to send the output")
	flag.Parse()

	if flag.NArg() != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [-output=FILEPATH] <bucket> <key>\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}
	bucket, key := flag.Arg(0), flag.Arg(1)

	s := s3.New(&aws.Config{Region: "eu-west-1"})

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
