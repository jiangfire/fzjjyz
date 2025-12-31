// Package utils provides key loading utilities (DRY: eliminate key loading repetition).
package utils

import (
	"fmt"

	"codeberg.org/jiangfire/fzjjyz/internal/crypto"
	"codeberg.org/jiangfire/fzjjyz/internal/i18n"
)

// LoadHybridPrivateKey loads hybrid private key (eliminates 4 repetitions).
func LoadHybridPrivateKey(path string) (*crypto.HybridPrivateKey, error) {
	key, err := crypto.LoadPrivateKeyCached(path)
	if err != nil {
		return nil, fmt.Errorf("load private key failed: %w",
			i18n.TranslateError("error.load_private_key_failed", err, path))
	}
	return key, nil
}

// LoadDilithiumVerifyKey loads signature verification public key (eliminates 3 repetitions).
func LoadDilithiumVerifyKey(path string) (interface{}, error) {
	if path == "" {
		return nil, nil
	}
	key, err := crypto.LoadDilithiumPublicKeyCached(path)
	if err != nil {
		return nil, fmt.Errorf("load verify key failed: %w",
			i18n.TranslateError("error.load_verify_key_failed", err, path))
	}
	return key, nil
}

// LoadHybridPublicKey loads hybrid public key.
func LoadHybridPublicKey(path string) (*crypto.HybridPublicKey, error) {
	key, err := crypto.LoadPublicKeyCached(path)
	if err != nil {
		return nil, fmt.Errorf("load public key failed: %w",
			i18n.TranslateError("error.load_public_key_failed", err, path))
	}
	return key, nil
}

// LoadDilithiumPrivateKey loads signature private key.
func LoadDilithiumPrivateKey(path string) (interface{}, error) {
	key, err := crypto.LoadDilithiumPrivateKeyCached(path)
	if err != nil {
		return nil, fmt.Errorf("load dilithium private key failed: %w",
			i18n.TranslateError("error.load_dilithium_private_key_failed", err, path))
	}
	return key, nil
}
