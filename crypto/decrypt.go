package crypto

func DecryptBytes(cipherData []byte) []byte {
	/**
	Decrypt a byte array
	*/

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

func DecryptString(data []byte) string {
	/**
	**validate logic**
	Decrypt incoming string
	*/

	// Convert incoming string to bytes
	var byte_data = []byte(data)

	// Encrypt bytes
	return string(DecryptBytes(byte_data))
}
