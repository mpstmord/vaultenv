// Package envelope provides encryption and decryption of secret values
// using a data encryption key (DEK) wrapped by a key encryption key (KEK).
package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidKeySize is returned when the provided key is not 32 bytes.
var ErrInvalidKeySize = errors.New("envelope: key must be 32 bytes (AES-256)")

// ErrDecryptFailed is returned when authenticated decryption fails.
var ErrDecryptFailed = errors.New("envelope: decryption failed")

// Cipher wraps and unwraps secret values using AES-256-GCM.
type Cipher struct {
	key []byte
}

// New creates a new Cipher using the provided 32-byte key.
func New(key []byte) (*Cipher, error) {
	if len(key) != 32 {
		return nil, ErrInvalidKeySize
	}
	copied := make([]byte, 32)
	copy(copied, key)
	return &Cipher{key: copied}, nil
}

// Seal encrypts plaintext using AES-256-GCM and returns the nonce+ciphertext.
func (c *Cipher) Seal(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("envelope: create GCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("envelope: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// Open decrypts a nonce+ciphertext blob produced by Seal.
func (c *Cipher) Open(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(c.key)
	if err != nil {
		return nil, fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("envelope: create GCM: %w", err)
	}
	ns := gcm.NonceSize()
	if len(ciphertext) < ns {
		return nil, ErrDecryptFailed
	}
	nonce, data := ciphertext[:ns], ciphertext[ns:]
	plaintext, err := gcm.Open(nil, nonce, data, nil)
	if err != nil {
		return nil, ErrDecryptFailed
	}
	return plaintext, nil
}
