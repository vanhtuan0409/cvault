package storage

import (
	"bytes"
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vanhtuan0409/cvault"
)

type s3Storage struct {
	bucket string
	client *s3.Client
}

func NewS3Storage(storeUrl string, client *s3.Client) *s3Storage {
	bucket := strings.TrimPrefix(storeUrl, "s3://")
	return &s3Storage{
		bucket: bucket,
		client: client,
	}
}

func (s *s3Storage) List(ctx context.Context) ([]string, error) {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})

	ret := []string{}
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return []string{}, err
		}

		for _, obj := range page.Contents {
			if key := *obj.Key; cvault.IsEncryptedName(key) {
				ret = append(ret, key)
			}
		}
	}
	return ret, nil
}

func (s *s3Storage) Get(ctx context.Context, key string) ([]byte, error) {
	downloader := manager.NewDownloader(s.client)
	wbuf := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(ctx, wbuf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return wbuf.Bytes(), nil
}

func (s *s3Storage) Put(ctx context.Context, key string, data []byte) error {
	uploader := manager.NewUploader(s.client)
	_, err := uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   bytes.NewBuffer(data),
	})
	return err
}

func (s *s3Storage) Remove(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	return err
}
