package crypto

import (
	"encoding/base64"

	"github.com/hashicorp/vault/api"
)

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

	// Update Vault Token
	client.SetToken(VAULT_TOKEN)

	// Decrypt the encrypted data key
	data := map[string]interface{}{
		"ciphertext": ENCRYPTED_DATA_KEY,
	}
	secret, err := client.Logical().Write("/transit/decrypt/trialkey-0", data)

	// Set DATA_KEY to plaintext value
	DATA_KEY, err := base64.StdEncoding.DecodeString(secret.Data["plaintext"].(string))
	if err != nil {
		return nil
	}

	return DATA_KEY
}
