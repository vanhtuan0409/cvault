package storage

import (
	"bytes"
	"context"
	"errors"
	"net/url"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/vanhtuan0409/cvault"
)

type s3Storage struct {
	bucket string
	client *s3.Client
}

func NewS3Storage(bucket string, client *s3.Client) *s3Storage {
	return &s3Storage{
		bucket: bucket,
		client: client,
	}
}

func (s *s3Storage) List(ctx context.Context) ([]*VaultItem, error) {
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})

	ret := []*VaultItem{}
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return ret, err
		}

		for _, obj := range page.Contents {
			if key := *obj.Key; cvault.IsEncryptedName(key) {
				ret = append(ret, &VaultItem{
					Key:          key,
					LastModified: *obj.LastModified,
				})
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

func parseS3StorageUrl(storeUrl string) (string, *s3.Client, error) {
	s3Url, err := url.Parse(storeUrl)
	if err != nil {
		return "", nil, err
	}
	if s3Url.Scheme != "s3" {
		return "", nil, errors.New("invalid schema")
	}

	ctx := context.TODO()
	bucket := s3Url.Hostname()
	opts := []func(*config.LoadOptions) error{}
	if user := s3Url.User; user != nil {
		accessKey := user.Username()
		secretKey, _ := user.Password()
		opts = append(opts, config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(accessKey, secretKey, ""),
		))
	}
	if region := s3Url.Query().Get("region"); region != "" {
		opts = append(opts, config.WithRegion(region))
	}
	if profile := s3Url.Query().Get("profile"); profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(profile))
	}

	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return "", nil, err
	}

	clientOpts := []func(*s3.Options){}
	if endpoint := s3Url.Query().Get("endpoint"); endpoint != "" {
		clientOpts = append(clientOpts, s3.WithEndpointResolver(customResolverFunc(endpoint)))
	}
	client := s3.NewFromConfig(cfg, clientOpts...)

	return bucket, client, nil
}

func customResolverFunc(endpoint string) s3.EndpointResolverFunc {
	return func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
		return aws.Endpoint{
			PartitionID:       "aws",
			URL:               endpoint,
			SigningRegion:     region,
			HostnameImmutable: true,
		}, nil
	}
}
