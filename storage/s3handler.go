package storage

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"os"
)

type S3Uploader struct {
	Client *s3.Client
	Bucket string
}

func NewS3Uploader(bucket string) (*S3Uploader, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-south-1"))
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg)

	return &S3Uploader{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (u *S3Uploader) Delete(s3Key string) error {
	_, err := u.Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(s3Key),
	})
	return err
}

func (u *S3Uploader) Upload(filePath, s3Key string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = u.Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(s3Key),
		Body:   file,
	})

	return err
}

func (u *S3Uploader) Download(s3Key, downloadPath string) error {
	out, err := u.Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(u.Bucket),
		Key:    aws.String(s3Key),
	})
	if err != nil {
		return err
	}
	defer out.Body.Close()

	file, err := os.Create(downloadPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, out.Body)
	return err
}
