package hodgepodge

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateRSAKeyPair(t *testing.T) {
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, publicKey, "Public key should not be nil")
	assert.NotNil(t, privateKey, "Private key should not be nil")
}

func TestRSAEncryptDecrypt(t *testing.T) {
	plaintext := []byte("Hello, world!")
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	assert.Nil(t, err, "Error should be nil")

	// Encrypt.
	ciphertext, err := RSAEncryptBytes(plaintext, publicKey)
	assert.Nil(t, err, "Error should be nil")
	assert.NotNil(t, ciphertext, "Ciphertext should not be nil")

	// Decrypt.
	decrypted, err := RSADecryptBytes(ciphertext, privateKey)
	assert.Nil(t, err, "Error should be nil")
	assert.Equal(t, plaintext, decrypted, "Plaintext and decrypted should be equal")
}

func TestRSAEncryptDecryptFile(t *testing.T) {
	tempDir, err := GetTempDir()
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tempDir)

	plaintextFile := filepath.Join(tempDir, "plaintext.txt")
	ciphertextFile := filepath.Join(tempDir, "ciphertext.txt")
	decryptedFile := filepath.Join(tempDir, "decrypted.txt")

	// Create a temporary file with some content.
	err = os.WriteFile(plaintextFile, []byte("Hello, world!"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}

	// Encrypt the file.
	publicKey, privateKey, err := GenerateRSAKeyPair(2048)
	assert.Nil(t, err, "Error should be nil")

	// Encrypt.
	err = RSAEncryptFile(plaintextFile, ciphertextFile, publicKey)
	assert.Nil(t, err, "Error should be nil")

	// Decrypt.
	err = RSADecryptFile(ciphertextFile, decryptedFile, privateKey)
	assert.Nil(t, err, "Error should be nil")

	// Compare the files.
	m, err := os.ReadFile(plaintextFile)
	assert.Nil(t, err, "Error should be nil")

	n, err := os.ReadFile(decryptedFile)
	assert.Nil(t, err, "Error should be nil")

	assert.Equal(t, m, n, "Plaintext and decrypted should be equal")
}
