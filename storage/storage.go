package storage

import (
	"context"
	"errors"
	"strings"
	"time"
)

type VaultItem struct {
	Key          string
	LastModified time.Time
}

type Storage interface {
	List(context.Context) ([]*VaultItem, error)
	Get(context.Context, string) ([]byte, error)
	Put(context.Context, string, []byte) error
	Remove(context.Context, string) error
}

func GetStorage(storeUrl string) (Storage, error) {
	switch {
	case strings.HasPrefix(storeUrl, "local://"):
		return NewLocalStorage(storeUrl), nil
	case strings.HasPrefix(storeUrl, "s3://"):
		bucket, client, err := parseS3StorageUrl(storeUrl)
		if err != nil {
			return nil, err
		}
		return NewS3Storage(bucket, client), nil
	default:
		return nil, errors.New("invalid storage source")
	}
}
