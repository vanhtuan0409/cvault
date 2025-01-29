package cvault

import (
	"context"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/core/registry"
	"github.com/tink-crypto/tink-go/v2/tink"
)

// nonce fixed associated data
var nonce = []byte("cvault")

func getPrimitive(ctx context.Context, keyUrl string) (tink.AEAD, error) {
	client, err := registry.GetKMSClient(keyUrl)
	if err != nil {
		return nil, err
	}
	kekAEAD, err := client.GetAEAD(keyUrl)
	if err != nil {
		return nil, err
	}
	return aead.NewKMSEnvelopeAEAD2(aead.AES128CTRHMACSHA256KeyTemplate(), kekAEAD), nil
}

func encryptWithTink(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	primitive, err := getPrimitive(ctx, keyUrl)
	if err != nil {
		return nil, err
	}

	return primitive.Encrypt(data, nonce)
}

func decryptWithTink(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	primitive, err := getPrimitive(ctx, keyUrl)
	if err != nil {
		return nil, err
	}

	return primitive.Decrypt(data, nonce)
}
