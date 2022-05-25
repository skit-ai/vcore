package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
)

// Read Env Vars
var VAULT_URI string = os.Getenv("VAULT_URI")
var VAULT_ROLE_ID string = os.Getenv("VAULT_ROLE_ID")
var VAULT_SECRET_ID string = os.Getenv("VAULT_SECRET_ID")
var VAULT_APPROLE_MOUNTPATH string = os.Getenv("VAULT_APPROLE_MOUNTPATH")
var VAULT_DATA_KEY_NAME string = os.Getenv("VAULT_DATA_KEY_NAME")
var ENCRYPTED_DATA_KEY string = os.Getenv("ENCRYPTED_DATA_KEY")

// Other Global Variables
var DATA_KEY []byte

// Vault functions
func getApproleAuth() *auth.AppRoleAuth {
	// Check if VAULT_APPROLE_MOUNTPATH has a value
	if len(VAULT_APPROLE_MOUNTPATH) == 0 {
		VAULT_APPROLE_MOUNTPATH = "approle-batch"
	}

	secretID := &auth.SecretID{
		FromString: VAULT_SECRET_ID,
	}
	appRoleAuth, err := auth.NewAppRoleAuth(VAULT_ROLE_ID, secretID, auth.WithMountPath(VAULT_APPROLE_MOUNTPATH))
	if err != nil {
		return nil
	}

	return appRoleAuth
}

func getDataKey() []byte {
	if len(DATA_KEY) != 0 {
		return DATA_KEY
	}

	// Initialize vault client
	config := &api.Config{
		Address: VAULT_URI,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil
	}

	appRoleAuth := getApproleAuth()
	secret_, err := client.Auth().Login(context.TODO(), appRoleAuth)
	if err != nil {
		return nil
	}

	// TODO: make client object a singleton
	// Initialize a lifetimewatcher to renew the token when it expires;
	// this would make more sense if the above code is reused and not re-initialized everytime
	client.NewLifetimeWatcher(&api.LifetimeWatcherInput{
		Secret: secret_,
	})

	// Decrypt the encrypted data key
	data := map[string]interface{}{
		"ciphertext": ENCRYPTED_DATA_KEY,
	}
	secret, err := client.Logical().Write("/transit/decrypt/"+VAULT_DATA_KEY_NAME, data)
	if err != nil {
		return nil
	}

	// Set DATA_KEY to plaintext value
	DATA_KEY, err := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return nil
	}

	return DATA_KEY
}

// Crypto functions
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
