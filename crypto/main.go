// A module to help with cryptographic requirements like encryption and hashing
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

func getDataKey(encrypted_data_key_ string) (data_key []byte) {
	if len(data_key) != 0 {
		return data_key
	}

	// Initialize vault client
	config := &api.Config{
		Address: vault_uri,
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

	// Check if data key is passed as a parameter
	var ciphertext string

	if encrypted_data_key_ != "" {
		ciphertext = encrypted_data_key_
	} else {
		ciphertext = encrypted_data_key
	}

	// Decrypt the encrypted data key
	data := map[string]interface{}{
		"ciphertext": ciphertext,
	}
	secret, err := client.Logical().Write("/transit/decrypt/"+vault_data_key_name, data)
	if err != nil {
		return nil
	}

	// Set data_key to plaintext value
	data_key, err = base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return nil
	}

	return
}

// Crypto functions
func newCipherAESGCMObject(data_key_b64_str string) (gcm cipher.AEAD) {

	// Get data key
	data_key := getDataKey(data_key_b64_str)

	// Generate new aes cipher using our 32 byte key
	c, err := aes.NewCipher(data_key)
	if err != nil {
		return nil
	}

	// GCM or Galois/Counter Mode, is a mode of operation for symmetric key cryptographic block ciphers
	gcm, err = cipher.NewGCM(c)
	if err != nil {
		return nil
	}

	return
}
