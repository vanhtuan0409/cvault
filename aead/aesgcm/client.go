package aesgcm

import (
	"errors"
	"strings"

	"github.com/google/tink/go/core/registry"
	"github.com/google/tink/go/tink"
)

const (
	aesPrefix = "aesgcm://"
)

type aesClient struct {
	keyPrefix    string
	promptScript string
}

func NewClient(urlPrefix string, promptScript string) (registry.KMSClient, error) {
	if !strings.HasPrefix(strings.ToLower(urlPrefix), aesPrefix) {
		return nil, errors.New("invalid aes key url")
	}

	return &aesClient{
		keyPrefix:    urlPrefix,
		promptScript: promptScript,
	}, nil
}

func (c *aesClient) Supported(keyURI string) bool {
	return strings.HasPrefix(keyURI, c.keyPrefix)
}

func (c *aesClient) GetAEAD(keyURI string) (tink.AEAD, error) {
	if !c.Supported(keyURI) {
		return nil, errors.New("unsupported key uri")
	}

	return &aesAEAD{
		promptScript: c.promptScript,
	}, nil
}
