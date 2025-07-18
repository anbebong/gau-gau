package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

var Key = []byte("1234567897654321") // 16 bytes for AES-128
var block cipher.Block

func InitCipher() error {
	var err error
	block, err = aes.NewCipher(Key)
	return err
}

func Encrypt(plaintext string) (string, error) {
	if block == nil {
		return "", fmt.Errorf("cipher not initialized")
	}
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))
	return base64.StdEncoding.EncodeToString(ciphertext), nil
	// return plaintext, nil // Placeholder, implement actual encryption
}

func Decrypt(cryptoText string) (string, error) {
	if block == nil {
		return "", fmt.Errorf("cipher not initialized")
	}
	ciphertext, err := base64.StdEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
	// return cryptoText, nil // Placeholder, implement actual decryption
}
