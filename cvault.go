package cvault

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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
	Key      []byte
	Nonce    [nonceLength]byte
	Message  []byte
	Checksum string
}

func (p *payload) doCheckSum() (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write(p.Key); err != nil {
		return "", err
	}
	if _, err := hasher.Write(p.Nonce[:]); err != nil {
		return "", err
	}
	if _, err := hasher.Write(p.Message); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (p *payload) setCheckSum() error {
	cs, err := p.doCheckSum()
	if err != nil {
		return err
	}
	p.Checksum = cs
	return nil
}

func (p *payload) verifyChecksum() error {
	cs, err := p.doCheckSum()
	if err != nil {
		return err
	}
	if p.Checksum != cs {
		return errors.New("checksum not matched")
	}
	return nil
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
	if err := p.setCheckSum(); err != nil {
		return nil, err
	}
	return msgpack.Marshal(&p)
}

func Decrypt(ctx context.Context, client *kms.Client, keyId string, input []byte) ([]byte, error) {
	// decode ciphertext with gob
	var p payload
	if err := msgpack.Unmarshal(input, &p); err != nil {
		return nil, err
	}
	if err := p.verifyChecksum(); err != nil {
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
