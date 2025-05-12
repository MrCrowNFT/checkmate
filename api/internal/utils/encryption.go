package utils

import (
	"encoding/base64"
	"errors"
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
