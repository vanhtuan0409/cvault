package storage

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Storage interface {
	List(context.Context) ([]string, error)
	Get(context.Context, string) ([]byte, error)
	Put(context.Context, string, []byte) error
	Remove(context.Context, string) error
}

func GetStorage(storeUrl string, s3Client *s3.Client) (Storage, error) {
	switch {
	case strings.HasPrefix(storeUrl, "local://"):
		return NewLocalStorage(storeUrl), nil
	case strings.HasPrefix(storeUrl, "s3://"):
		return NewS3Storage(storeUrl, s3Client), nil
	default:
		return nil, errors.New("invalid storage source")
	}
}
