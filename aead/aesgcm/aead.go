package aesgcm

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"

	"golang.org/x/crypto/scrypt"
)

const (
	saltLength  = 8
	nonceLength = 12
	keyLength   = 32
)

type aesAEAD struct {
	promptScript string
}

func (a *aesAEAD) Encrypt(plaintext, associatedData []byte) ([]byte, error) {
	salt := make([]byte, saltLength)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, err
	}

	aesgcm, err := getPrimitive(a.promptScript, salt)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	encrypted := aesgcm.Seal(nil, nonce, plaintext, associatedData)
	return pack(encrypted, salt, nonce), nil
}

func (a *aesAEAD) Decrypt(tuple, associatedData []byte) ([]byte, error) {
	ciphertext, salt, nonce, err := unpack(tuple)
	if err != nil {
		return nil, err
	}

	aesgcm, err := getPrimitive(a.promptScript, salt)
	if err != nil {
		return nil, err
	}

	return aesgcm.Open(nil, nonce, ciphertext, associatedData)
}

func getPrimitive(promptScript string, salt []byte) (cipher.AEAD, error) {
	passphrase, err := promptPassPhrase(promptScript)
	if err != nil {
		return nil, err
	}
	masterKey, err := deriveKey(passphrase, salt)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}
	return cipher.NewGCM(block)
}

func pack(ciphertext, salt, nonce []byte) []byte {
	ret := salt
	ret = append(ret, nonce...)
	ret = append(ret, ciphertext...)
	return ret
}

func unpack(tuple []byte) (ciphertext []byte, salt []byte, nonce []byte, err error) {
	if len(tuple) <= (saltLength + nonceLength) {
		return nil, nil, nil, errors.New("cannot unpack ciphertext and nonce")
	}
	salt, nonce, ciphertext = tuple[0:saltLength], tuple[saltLength:nonceLength+saltLength], tuple[nonceLength+saltLength:]
	return
}

func deriveKey(passphrase string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(passphrase), salt, 1<<15, 8, 1, keyLength)
}
