package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func SaveFile(src io.Reader, fileHeader *multipart.FileHeader) (string, error) {
	fileName := fileHeader.Filename
	fmt.Println(fileName)
	fileName = strings.ReplaceAll(fileName, " ", "")

	path := filepath.Join("/home/murarka/chat_app/server/uploads", fileName)
	save, err := os.Create(path)
	if err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}
	defer save.Close()

	if _, err := io.Copy(save, src); err != nil {
		return "", fmt.Errorf("error saving image: %v", err)
	}

	url, err := UploadtoS3(fileName, path)
	if err != nil {
		return "", fmt.Errorf("error uploading to S3: %v", err)
	}

	return url, nil
}

var (
	awsRegion   = "us-east-1"
	awsEndpoint = "http://localhost:4566"
	bucketName  = "test-bucket"
	s3svc       *s3.Client
)

func init() {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if awsEndpoint != "" {
			return aws.Endpoint{
				PartitionID:   "aws",
				URL:           awsEndpoint,
				SigningRegion: awsRegion,
			}, nil
		}

		// returning EndpointNotFoundError will allow the service to fallback to it's default resolution
		return aws.Endpoint{}, &aws.EndpointNotFoundError{}
	})

	awsCfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(awsRegion),
		config.WithEndpointResolverWithOptions(customResolver),
	)
	if err != nil {
		log.Fatalf("Cannot load the AWS configs: %s", err)
	}

	s3svc = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
}

func UploadtoS3(filename string, path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Printf("unable to upload file %v", err)
		return "", err
	}
	defer file.Close()

	_, err = s3svc.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(filename),
		Body:   file,
	})

	if err != nil {
		log.Printf("unable to upload file %v", err)
		return "", err
	}
	url := fmt.Sprintf("http://localhost:4566/%s/%s", bucketName, filename)

	return url, nil
}

//TODO DELETE FILE FROM DEVICE AFTER SUCCESSFUL ENTRY IN LOCALSTACK
