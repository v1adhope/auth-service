package serialization

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

func EncryptByGcm(text, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("serialization: encryption: EncryptByGcm: NewCipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("serialization: encryption: EncryptByGcm: NewGCM: %w", err)
	}

	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("serialization: encryption: EncryptByGcm: ReadFull: %w", err)
	}

	ciphertext := aesgcm.Seal(nonce, nonce, text, nil)

	return ciphertext, nil
}

func DecryptByGcm(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("serialization: encryption: DecryptByGcm: DecryptByGcm: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("serialization: encryption: DecryptByGcm: NewGCM: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	text, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("serialization: encryption: DecryptByGcm: Open: %w", err)
	}

	return text, nil
}
