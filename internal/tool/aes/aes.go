package aes

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

type Service struct {
	Key []byte
	Iv  []byte
}

// PKCS7

func NewService(key, iv string) *Service {
	keyBytes, _ := base64.StdEncoding.DecodeString(key)
	ivBytes, _ := base64.StdEncoding.DecodeString(iv)
	return &Service{
		Key: keyBytes,
		Iv:  ivBytes,
	}
}

func GetAesService(key []byte, iv []byte) *Service {
	return &Service{
		Key: key,
		Iv:  iv,
	}
}

func GenerateKeyIv() (string, string) {
	key := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		panic(err)
	}

	keyText := base64.StdEncoding.EncodeToString(key)

	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	ivText := base64.StdEncoding.EncodeToString(iv)
	return keyText, ivText
}

func (s *Service) Crypt(data string) (string, error) {
	plaintext := []byte(data)
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}

	paddedPlaintext := pad(plaintext, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, s.Iv)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Service) CryptBytes(plaintext []byte) (string, error) {
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}

	paddedPlaintext := pad(plaintext, aes.BlockSize)

	mode := cipher.NewCBCEncrypter(block, s.Iv)
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (s *Service) Decrypt(data string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, s.Iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return string(unpad(plaintext)), nil
}

func (s *Service) DecryptGetBytes(data string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	block, err := aes.NewCipher(s.Key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(block, s.Iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	return unpad(plaintext), nil
}

func pad(input []byte, blockSize int) []byte {
	padding := blockSize - (len(input) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(input, padText...)
}

func unpad(input []byte) []byte {
	padding := int(input[len(input)-1])
	return input[:len(input)-padding]
}
