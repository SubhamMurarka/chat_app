package S3

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SubhamMurarka/chat_app/Image/Config"
	"github.com/SubhamMurarka/chat_app/Image/models"
	"github.com/SubhamMurarka/chat_app/Image/redis"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConnectS3(c *gin.Context) {
	id := c.GetString("userid")

	roomid := c.DefaultQuery("roomid", "")
	if roomid == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "missing roomid"})
		return
	}

	key := fmt.Sprintf("%s:%s", id, roomid)

	is := redis.IsUserActive(key)
	if !is {
		c.JSON(http.StatusForbidden, gin.H{"error": "NOT ALLOWED"})
		return
	}

	RandomID := uuid.New().String()

	bucketName := Config.Conf.BucketName
	objectKey := Config.Conf.ObjectKey + RandomID

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Unable to load AWS SDK config: %v", err)
	}

	s3client := s3.NewFromConfig(cfg)
	presignedClient := s3.NewPresignClient(s3client)
	presignedurl, err := presignedClient.PresignPutObject(c.Request.Context(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		ACL:    "public-read",
	},
		s3.WithPresignExpires(30*time.Second),
	)

	if err != nil {
		log.Println("Error generating presigned URL:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL"})
		return
	}

	res := &models.Presigned{
		PresignedURL: presignedurl,
		Key:          RandomID,
	}

	c.JSON(http.StatusOK, gin.H{"response": res})
}
