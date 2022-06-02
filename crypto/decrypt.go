package crypto

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
// This function accepts an incoming byte array, decrypts it using AES-256 decryption, converts the result into a string and returns the string
func DecryptString(data []byte) string {
	/**
	**validate logic**
	Decrypt incoming string
	*/

	// Convert incoming string to bytes
	var byte_data = []byte(data)

	// Decrypt bytes, convert to string and return
	return string(DecryptBytes(byte_data))
}
