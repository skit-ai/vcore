// A module to help with cryptographic requirements like encryption and hashing
package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
	auth "github.com/hashicorp/vault/api/auth/approle"
	"github.com/skit-ai/vcore/env"
)

// Read Env Vars

var vault_uri string = os.Getenv("VAULT_URI")
var vault_role_id string = os.Getenv("VAULT_ROLE_ID")
var vault_secret_id string = os.Getenv("VAULT_SECRET_ID")
var vault_approle_mountpath string = os.Getenv("VAULT_APPROLE_MOUNTPATH")
var vault_data_key_name string = os.Getenv("VAULT_DATA_KEY_NAME")
var encrypted_data_key string = os.Getenv("ENCRYPTED_DATA_KEY")
var use_static_data_key bool = env.Bool("USE_STATIC_DATA_KEY", false)
var static_data_key string = env.String("STATIC_DATA_KEY", "")
var log_crypto_internal_info bool = env.Bool("LOG_CRYPTO_INTERNAL_INFO", false)

// Other Global Variables

var data_key []byte
var dataKeyCache map[string][]byte = map[string][]byte{}

func isValidBase64(static_data_key string) bool {
	_, err := base64.StdEncoding.DecodeString(static_data_key)
	return err == nil
}

func getByteString(static_data_key string) []byte {
	return []byte(static_data_key)
}

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

	var data_key []byte
	// Get data key
	if use_static_data_key && isValidBase64(static_data_key) {
		data_key = getByteString(static_data_key)
	} else {
		data_key = getDataKey(data_key_b64_str, clientId)
	}

	if log_crypto_internal_info {
		fmt.Printf("Data Key obtained - %s", data_key)
	}

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
