package cvault

import (
	"context"
)

func Encrypt(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	switch {
	default:
		return encryptWithTink(ctx, keyUrl, data)
	}
}

func Decrypt(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	switch {
	default:
		return decryptWithTink(ctx, keyUrl, data)
	}
}
