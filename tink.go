package cvault

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/google/tink/go/aead"
	"github.com/google/tink/go/keyset"
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

func encryptWithTink(ctx context.Context, keyUrl string, data []byte) ([]byte, error) {
	dek := aead.AES128CTRHMACSHA256KeyTemplate()
	kh, err := keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(keyUrl, dek))
	if err != nil {
		return nil, err
	}

	a, err := aead.New(kh)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted, err := a.Encrypt(data, nonce)
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

	dek := aead.AES128CTRHMACSHA256KeyTemplate()
	kh, err := keyset.NewHandle(aead.KMSEnvelopeAEADKeyTemplate(keyUrl, dek))
	if err != nil {
		return nil, err
	}

	a, err := aead.New(kh)
	if err != nil {
		return nil, err
	}

	return a.Decrypt(p.Message, p.Nonce[:])
}
