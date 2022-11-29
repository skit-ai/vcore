// The crypto module is meant to help services implement various cryptographic
// functions with ease.
//
// Current features
//
// 1. Encryption of []byte and string.
//
// Supported techniques: AES-256-GCM
//
// 2. Decryption of []byte.
//
// Supported techniques: AES-256-GCM
//
// AES-256 is PCI DSS compliant, as it is a recognised industry standard
// encryption.
//
// Key management
//
// Vault is used to generate the encrypted data key when an environment/client is set up. The encrypted data key is passed to vcore as an environment variable.
// Vcore then calls Vault APIs to decrypt the data key and proceed with the encryption/decryption.
//
// Environment Variables needed
//
// The following environment variables are needed to utilize the crypto module:
//
//     export VAULT_URI="http://localhost:8200"
//     export VAULT_ROLE_ID="****"
//     export VAULT_SECRET_ID="****"
//     export VAULT_APPROLE_MOUNTPATH="approle"
//     export ENCRYPTED_DATA_KEY="****"
//     export VAULT_DATA_KEY_NAME="datakey-name"
//
// Note: the above environment variables are just examples, set up vault and replace the actual values above.
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

var vault_uri string = os.Getenv("VAULT_URI")
var vault_role_id string = os.Getenv("VAULT_ROLE_ID")
var vault_secret_id string = os.Getenv("VAULT_SECRET_ID")
var vault_approle_mountpath string = os.Getenv("VAULT_APPROLE_MOUNTPATH")
var vault_data_key_name string = os.Getenv("VAULT_DATA_KEY_NAME")
var encrypted_data_key string = os.Getenv("ENCRYPTED_DATA_KEY")

// Other Global Variables

var data_key []byte
var dataKeyCache map[string][]byte = map[string][]byte{}

// Vault functions
func getApproleAuth() *auth.AppRoleAuth {
	// Check if vault_approle_mountpath has a value
	if len(vault_approle_mountpath) == 0 {
		vault_approle_mountpath = "approle-batch"
	}

	secretID := &auth.SecretID{
		FromString: vault_secret_id,
	}
	appRoleAuth, err := auth.NewAppRoleAuth(vault_role_id, secretID, auth.WithMountPath(vault_approle_mountpath))
	if err != nil {
		return nil
	}

	return appRoleAuth
}

func getDataKey(encrypted_data_key_ string, clientId string) (data_key_ []byte) {

	var cachedKeyPresent bool

	if clientId != "" {
		data_key_, cachedKeyPresent = dataKeyCache[clientId]
	}

	if cachedKeyPresent {
		return
	}

	// If clientId is not provided, check if global data key is set (environment variable)
	if len(data_key) != 0 && clientId == "" {
		data_key_ = data_key
		return
	}

	// If no cache value found, retrieve unencrypted data key value from vault

	// Initialize vault client
	config := &api.Config{
		Address: vault_uri,
	}
	client, err := api.NewClient(config)
	if err != nil {
		return nil
	}

	// Initialize approle auth
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

	// Check if data key is passed as a parameter
	var ciphertext string
	var isGlobal bool = false
	var vaultDataKeyName_ string

	if encrypted_data_key_ != "" {
		ciphertext = encrypted_data_key_
	} else {
		ciphertext = encrypted_data_key
		isGlobal = true
	}

	if clientId != "" {
		vaultDataKeyName_ = clientId
	} else {
		vaultDataKeyName_ = vault_data_key_name
	}

	// Decrypt the encrypted data key
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}
	secret, err := client.Logical().Write("/transit/decrypt/"+vaultDataKeyName_, data)
	if err != nil {
		return nil
	}

	// Set data_key_ to plaintext value
	data_key_, err = base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return nil
	}

	// Set clientId based cache
	if clientId != "" {
		dataKeyCache[clientId] = data_key_
	}

	// Set global data_key
	if isGlobal {
		data_key = data_key_
	}

	return
}

// Crypto functions
func newCipherAESGCMObject(data_key_b64_str string, clientId string) (gcm cipher.AEAD, err error) {

	// Get data key
	data_key := getDataKey(data_key_b64_str, clientId)

	// Generate new aes cipher using our 32 byte key
	c, err := aes.NewCipher(data_key)
	if err != nil {
		return
	}

	// GCM or Galois/Counter Mode, is a mode of operation for symmetric key cryptographic block ciphers
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return
	}

	return
}
