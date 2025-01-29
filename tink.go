package cvault

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/tink-crypto/tink-go/v2/aead"
	"github.com/tink-crypto/tink-go/v2/core/registry"
	"github.com/tink-crypto/tink-go/v2/tink"
	"github.com/vmihailenco/msgpack/v5"
)

const (
	keyLength   = 32
	nonceLength = 24
)

type payload struct {
	Nonce   []byte
	Message []byte
}

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

	nonce := make([]byte, nonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted, err := primitive.Encrypt(data, nonce)
	if err != nil {
		return nil, err
	}

	return msgpack.Marshal(&payload{
		Message: encrypted,
		Nonce:   nonce,
	})
}

func decryptWithTink(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	var p payload
	if err := msgpack.Unmarshal(data, &p); err != nil {
		return nil, err
	}

	primitive, err := getPrimitive(ctx, keyUrl)
	if err != nil {
		return nil, err
	}

	return primitive.Decrypt(p.Message, p.Nonce[:])
}
