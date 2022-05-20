package crypto

import (
	"crypto/rand"
	"io"
)

func EncryptBytes(data []byte) []byte {
	/**
	Encrypt a byte array
	*/

	gcm := newCipherAESGCMObject()
	if gcm == nil {
		return nil
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, 12) // gcm.NonceSize() also defaults to 12

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil
	}

	encrypted_bytes := gcm.Seal(nonce, nonce, data, nil)

	return encrypted_bytes
}

func EncryptString(data string) []byte {
	/**
	Decrypt a string
	*/

	// Convert incoming string to bytes
	var byte_data = []byte(data)

	// Encrypt bytes
	return EncryptBytes(byte_data)
}
