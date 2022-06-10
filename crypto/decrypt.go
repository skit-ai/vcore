package crypto

import (
	"encoding/base64"

	"github.com/Vernacular-ai/vcore/errors"
)

// Decrypt a byte array
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption and returns the result in bytes
func DecryptBytes(cipherData []byte) []byte {
	gcm := newCipherAESGCMObject()
	if gcm == nil {
		return nil
	}

	nonceSize := 12 // gcm.NonceSize() also defaults to 12
	if len(cipherData) < nonceSize {
		return nil
	}

	// Make sure auth(tag/mac) is always 16 bits(how?)
	nonce, cipher_withauth := cipherData[:nonceSize], cipherData[nonceSize:]
	data, err := gcm.Open(nil, nonce, cipher_withauth, nil)
	if err != nil {
		return nil
	}

	return data
}

// Decrypt a byte array
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption,
// converts the result into a string and returns the string
func DecryptString(data []byte) string {
	// Decrypt bytes, convert to string and return
	return string(DecryptBytes(data))
}

// Decrypt a base64-encoded encrypted string to unencrypted string
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func, converts the result into a string and returns resultant string.
// Note: Only use when string data was encrypted.
func DecryptB64ToString(data string) (decrypted_string string, err error) {
	// Convert incoming string to bytes
	byte_data, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decrypted_data := DecryptBytes(byte_data)

	// Convert to string
	decrypted_string = string(decrypted_data)

	return
}

// Decrypt a base64-encoded encrypted string to unencrypted bytes
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func and returns resultant byte array.
func DecryptB64ToBytes(data string) (decrypted_data []byte, err error) {

	// Convert incoming string to bytes
	byte_data, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decrypted_data = DecryptBytes(byte_data)

	return
}
