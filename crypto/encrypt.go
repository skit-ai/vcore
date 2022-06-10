package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

// Encrypt a byte array
//
// This function accepts an incoming byte array, encrypts it using AES-256 decryption and returns the result in bytes
func EncryptBytes(data []byte) []byte {
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

// Encrypt a byte array
//
// This function accepts an incoming string, encrypts it using EncryptBytes func and returns the result in bytes
func EncryptString(data string) []byte {
	// Convert incoming string to bytes
	var byte_data = []byte(data)

	// Encrypt bytes
	return EncryptBytes(byte_data)
}

// Encrypt a string
//
// This function accepts an incoming string, encrypts it using EncryptBytes func,
// encodes the bytearray to base64 string and returns the resultant string
func EncryptToB64String(data string) (encrypted_data_b64_str string, err error) {

	// Convert incoming string to bytes
	var byte_data = []byte(data)

	// Encrypt bytes
	encrypted_data_bytes := EncryptBytes(byte_data)

	// Encode encrypted bytes to b64 string
	encrypted_data_b64_str = base64.StdEncoding.EncodeToString(encrypted_data_bytes)

	return
}
