package crypto

import (
	"encoding/base64"

	"github.com/Vernacular-ai/vcore/errors"
)

// Decrypt a byte array
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption and returns the result in bytes
func DecryptBytesWithDataKey(cipherData []byte, dataKey string, clientId string) (data []byte, err error) {
	gcm, err := newCipherAESGCMObject(dataKey, clientId)
	if gcm == nil || err != nil {
		return
	}

	nonceSize := 12 // gcm.NonceSize() also defaults to 12
	if len(cipherData) < nonceSize {
		return
	}

	// Make sure auth(tag/mac) is always 16 bits(how?)
	nonce, cipherWithAuth := cipherData[:nonceSize], cipherData[nonceSize:]
	data, err = gcm.Open(nil, nonce, cipherWithAuth, nil)
	if err != nil {
		return
	}

	return
}

// Decrypt a byte array
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption,
// converts the result into a string and returns the string
func DecryptStringWithDataKey(data []byte, dataKey string, clientId string) (stringData string, err error) {
	// Decrypt bytes, convert to string and return
	byteData, err := DecryptBytesWithDataKey(data, dataKey, clientId)
	return string(byteData), err
}

// Decrypt a base64-encoded encrypted string to unencrypted string
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func, converts the result into a string and returns resultant string.
// Note: Only use when string data was encrypted.
func DecryptB64ToStringWithDataKey(data string, dataKey string, clientId string) (decryptedString string, err error) {
	// Convert incoming string to bytes
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decryptedData, err := DecryptBytesWithDataKey(byteData, dataKey, clientId)

	// Convert to string
	decryptedString = string(decryptedData)

	return
}

// Decrypt a base64-encoded encrypted string to unencrypted bytes
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func and returns resultant byte array.
func DecryptB64ToBytesWithDataKey(data string, dataKey string, clientId string) (decryptedData []byte, err error) {

	// Convert incoming string to bytes
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decryptedData, err = DecryptBytesWithDataKey(byteData, dataKey, clientId)

	return
}

/**
Decryption functions without data key
*/

// Decrypt a byte array
// Deprecated - to be removed in future releases - use DecryptBytesWithDataKey instead
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption and returns the result in bytes
func DecryptBytes(cipherData []byte) (data []byte, err error) {
	gcm, err := newCipherAESGCMObject("", "")
	if gcm == nil || err != nil {
		return
	}

	nonceSize := 12 // gcm.NonceSize() also defaults to 12
	if len(cipherData) < nonceSize {
		return
	}

	// Make sure auth(tag/mac) is always 16 bits(how?)
	nonce, cipherWithAuth := cipherData[:nonceSize], cipherData[nonceSize:]
	data, err = gcm.Open(nil, nonce, cipherWithAuth, nil)
	if err != nil {
		return
	}

	return
}

// Decrypt a byte array
//
// This function accepts an incoming byte array, decrypts it using AES-256 decryption,
// converts the result into a string and returns the string
func DecryptString(data []byte) (stringData string, err error) {
	// Decrypt bytes, convert to string and return
	byteData, err := DecryptBytesWithDataKey(data, "", "")
	return string(byteData), err
}

// Decrypt a base64-encoded encrypted string to unencrypted string
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func, converts the result into a string and returns resultant string.
// Note: Only use when string data was encrypted.
func DecryptB64ToString(data string) (decryptedString string, err error) {
	// Convert incoming string to bytes
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decryptedData, err := DecryptBytesWithDataKey(byteData, "", "")

	// Convert to string
	decryptedString = string(decryptedData)

	return
}

// Decrypt a base64-encoded encrypted string to unencrypted bytes
//
// This function accepts an incoming base64 encoded string, base64 decodes it,
// decrypts it using EncryptBytes func and returns resultant byte array.
func DecryptB64ToBytes(data string) (decryptedData []byte, err error) {

	// Convert incoming string to bytes
	byteData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		err = errors.NewError("Failed to base64-decode incoming string, please check if "+
			"base64 encoded string is supplied", err, false)
		return
	}

	// Decrypt bytes
	decryptedData, err = DecryptBytesWithDataKey(byteData, "", "")

	return
}
