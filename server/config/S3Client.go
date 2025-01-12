package config

import (
	"bytes"
	"context"
	"io"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var FileContent string

func ConnectS3() {

	bucketName := "shubhamgochat"
	objectKey := "AbuseMasker/Badwords.txt"

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	output, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		log.Fatalf("Unable to get object: %v", err)
	}

	defer output.Body.Close()

	var buf bytes.Buffer
	if _, err := io.Copy(&buf, output.Body); err != nil {
		log.Fatalf("Unable to read object content: %v", err)
	}

	FileContent = buf.String()
}
