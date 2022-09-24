package cvault

import (
	"context"
	"crypto/rand"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	keyLength   = 32
	nonceLength = 24
)

type payload struct {
	Key     []byte
	Nonce   [nonceLength]byte
	Message []byte
}

func Encrypt(ctx context.Context, client *kms.Client, keyId string, input []byte) ([]byte, error) {
	// gen data key
	resp, err := client.GenerateDataKey(ctx, &kms.GenerateDataKeyInput{
		KeyId:         aws.String(keyId),
		NumberOfBytes: aws.Int32(keyLength),
	})
	if err != nil {
		return nil, err
	}
	dataKey := [keyLength]byte{}
	copy(dataKey[:], resp.Plaintext)

	// gen nonce
	nonce := [nonceLength]byte{}
	if _, err = rand.Read(nonce[:]); err != nil {
		return nil, err
	}

	// encrypt plaintext
	var encrypted []byte
	encrypted = secretbox.Seal(encrypted, input, &nonce, &dataKey)

	// bundle payload
	p := payload{
		Key:     resp.CiphertextBlob,
		Nonce:   nonce,
		Message: encrypted,
	}
	return msgpack.Marshal(&p)
}

func Decrypt(ctx context.Context, client *kms.Client, keyId string, input []byte) ([]byte, error) {
	// decode ciphertext with gob
	var p payload
	if err := msgpack.Unmarshal(input, &p); err != nil {
		return nil, err
	}

	// decrypt data key
	resp, err := client.Decrypt(ctx, &kms.DecryptInput{
		KeyId:          aws.String(keyId),
		CiphertextBlob: p.Key,
	})
	if err != nil {
		return nil, err
	}
	dataKey := [keyLength]byte{}
	copy(dataKey[:], resp.Plaintext)

	// decrypt ciphertext
	var (
		decrypted []byte
		ok        bool
	)
	decrypted, ok = secretbox.Open(decrypted, p.Message, &p.Nonce, &dataKey)
	if !ok {
		return nil, fmt.Errorf("failed to open secretbox")
	}

	return decrypted, nil
}
