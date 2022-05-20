package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"os"
)

// Read Env Vars
var VAULT_URI string = os.Getenv("VAULT_URI")
var VAULT_TOKEN string = os.Getenv("VAULT_TOKEN")
var ENCRYPTED_DATA_KEY string = os.Getenv("ENCRYPTED_DATA_KEY")

// Other Global Variables
var DATA_KEY []byte

// Some common functions
func newCipherAESGCMObject() cipher.AEAD {
	// Get data key
	var data_key = getDataKey()

	// Generate new aes cipher using our 32 byte key
	c, err := aes.NewCipher(data_key)
	if err != nil {
		return nil
	}

	// GCM or Galois/Counter Mode, is a mode of operation for symmetric key cryptographic block ciphers
	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil
	}

	return gcm
}
