package models

import (
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type Presigned struct {
	PresignedURL *v4.PresignedHTTPRequest `json:"presignedURL"`
	Key          string                   `json:"key"`
	Url          string                   `json:"url"`
}
