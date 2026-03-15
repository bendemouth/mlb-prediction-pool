package handlers

import (
	"context"
	"mime/multipart"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Handler struct {
	client *s3.Client
}

func NewS3Client(ctx context.Context, options s3.Options) (*S3Handler, error) {
	return &S3Handler{client: s3.New(options)}, nil
}

func (s *S3Handler) UploadFileToS3(fileInput multipart.File, s3Key string, ctx context.Context) (isSuccess bool, returnedKey string, err error) {
	client := s.client

	_, err = client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("AWS_S3_BUCKET_NAME")),
		Key:    aws.String(s3Key),
		Body:   fileInput,
	})

	return err == nil, s3Key, err
}
