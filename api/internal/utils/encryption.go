package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

var (
	//32 bytes for AES-256
	encryptionKey []byte
)

// initializes the encryption system
func InitEncryption() error {
	// load key from environment variable
	keyStr := os.Getenv("ENCRYPTION_KEY")
	if keyStr == "" {
		return errors.New("ENCRYPTION_KEY environment variable not set")
	}

	// decode the base64 encoded key
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return err
	}

	// ensure key length (32 bytes for AES-256)
	if len(key) != 32 {
		return errors.New("encryption key must be 32 bytes (256 bits)")
	}

	encryptionKey = key
	return nil
}

func EncryptString(plaintext string) (string, error) {
	if encryptionKey == nil {
		return "", errors.New("encryption not initialized")
	}
	
	// Create cipher block
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	
	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// Create nonce 
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	
	// Encrypt data
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	
	// Encode as base64 string for storage
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptString(encryptedStr string) (string, error) {
	if encryptionKey == nil {
		return "", errors.New("encryption not initialized")
	}
	
	// Decode from base64
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", err
	}
	
	// Create cipher block
	block, err := aes.NewCipher(encryptionKey)
	if err != nil {
		return "", err
	}
	
	// Create GCM cipher mode
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	
	// Ensure ciphertext is large enough
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("ciphertext too short")
	}
	
	// Extract nonce and ciphertext
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	
	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	
	return string(plaintext), nil
}