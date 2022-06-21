package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"io"
)

/**
Encryption functions with data key
*/

// Encrypt a byte array
//
// This function accepts an incoming byte array, encrypts it using AES-256 decryption and returns the result in bytes
func EncryptBytesWithDataKey(data []byte, dataKey string, clientId string) (encryptedBytes []byte, err error) {
	gcm, err := newCipherAESGCMObject(dataKey, clientId)
	if gcm == nil || err != nil {
		return
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, 12) // gcm.NonceSize() also defaults to 12

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	encryptedBytes = gcm.Seal(nonce, nonce, data, nil)

	return
}

// Encrypt a byte array
//
// This function accepts an incoming string, encrypts it using EncryptBytes func and returns the result in bytes.
func EncryptStringWithDataKey(data string, dataKey string, clientId string) (byteData []byte, err error) {
	// Convert incoming string to bytes
	byteData = []byte(data)

	// Encrypt bytes
	return EncryptBytesWithDataKey(byteData, dataKey, clientId)
}

// Encrypt a string
//
// This function accepts an incoming string, encrypts it using EncryptBytes func,
// encodes the bytearray to base64 string and returns the resultant string.
func EncryptToB64StringWithDataKey(data string, dataKey string, clientId string) (encryptedDataB64Str string, err error) {

	// Convert incoming string to bytes
	var byteData = []byte(data)

	// Encrypt bytes
	encryptedDataBytes, err := EncryptBytesWithDataKey(byteData, dataKey, clientId)

	// Encode encrypted bytes to b64 string
	encryptedDataB64Str = base64.StdEncoding.EncodeToString(encryptedDataBytes)

	return
}

/**
Encryption functions without data key
*/

// Encrypt a byte array
// Deprecated - to be removed in future releases - use EncryptBytesWithDataKey instead
//
// This function accepts an incoming byte array, encrypts it using AES-256 decryption and returns the result in bytes
func EncryptBytes(data []byte) (encryptedBytes []byte, err error) {
	gcm, err := newCipherAESGCMObject("", "")
	if gcm == nil || err != nil {
		return
	}

	// creates a new byte array the size of the nonce
	// which must be passed to Seal
	nonce := make([]byte, 12) // gcm.NonceSize() also defaults to 12

	// populates our nonce with a cryptographically secure
	// random sequence
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return
	}

	encryptedBytes = gcm.Seal(nonce, nonce, data, nil)

	return
}

// Encrypt a byte array
//
// This function accepts an incoming string, encrypts it using EncryptBytes func and returns the result in bytes
func EncryptString(data string) (byteData []byte, err error) {
	// Convert incoming string to bytes
	byteData = []byte(data)

	// Encrypt bytes
	return EncryptBytesWithDataKey(byteData, "", "")
}

// Encrypt a string
//
// This function accepts an incoming string, encrypts it using EncryptBytes func,
// encodes the bytearray to base64 string and returns the resultant string
func EncryptToB64String(data string) (encryptedDataB64Str string, err error) {

	// Convert incoming string to bytes
	var byteData = []byte(data)

	// Encrypt bytes
	encryptedDataBytes, err := EncryptBytesWithDataKey(byteData, "", "")

	// Encode encrypted bytes to b64 string
	encryptedDataB64Str = base64.StdEncoding.EncodeToString(encryptedDataBytes)

	return
}
