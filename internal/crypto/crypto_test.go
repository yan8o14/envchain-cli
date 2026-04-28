package crypto_test

import (
	"bytes"
	"testing"

	"github.com/envchain-cli/envchain-cli/internal/crypto"
)

func TestDeriveKey(t *testing.T) {
	key := crypto.DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Fatalf("expected key length 32, got %d", len(key))
	}

	key2 := crypto.DeriveKey("my-secret-passphrase")
	if !bytes.Equal(key, key2) {
		t.Fatal("same passphrase should produce same key")
	}

	key3 := crypto.DeriveKey("different-passphrase")
	if bytes.Equal(key, key3) {
		t.Fatal("different passphrases should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := crypto.DeriveKey("test-passphrase")
	plaintext := []byte("SECRET_KEY=super_secret_value")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if bytes.Equal(ciphertext, plaintext) {
		t.Fatal("ciphertext should not equal plaintext")
	}

	decrypted, err := crypto.Decrypt(key, ciphertext)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestDecryptWithWrongKey(t *testing.T) {
	key := crypto.DeriveKey("correct-passphrase")
	wrongKey := crypto.DeriveKey("wrong-passphrase")
	plaintext := []byte("API_TOKEN=abc123")

	ciphertext, err := crypto.Encrypt(key, plaintext)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(wrongKey, ciphertext)
	if err == nil {
		t.Fatal("expected decryption to fail with wrong key")
	}
}

func TestEncryptProducesUniqueNonces(t *testing.T) {
	key := crypto.DeriveKey("nonce-test")
	plaintext := []byte("same plaintext")

	c1, _ := crypto.Encrypt(key, plaintext)
	c2, _ := crypto.Encrypt(key, plaintext)

	if bytes.Equal(c1, c2) {
		t.Fatal("two encryptions of the same plaintext should produce different ciphertexts")
	}
}
