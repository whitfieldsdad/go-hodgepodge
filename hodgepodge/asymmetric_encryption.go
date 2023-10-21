package hodgepodge

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func RSAEncryptFile(inputPath, outputPath string, publicKey []byte) error {
	plaintext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	ciphertext, err := RSAEncryptBytes(plaintext, publicKey)
	if err != nil {
		return err
	}
	return os.WriteFile(outputPath, ciphertext, 0644)
}

func RSADecryptFile(inputPath, outputPath string, privateKey []byte) error {
	ciphertext, err := os.ReadFile(inputPath)
	if err != nil {
		return err
	}
	plaintext, err := RSADecryptBytes(ciphertext, privateKey)
	if err != nil {
		return nil
	}
	err = os.WriteFile(outputPath, plaintext, 0644)
	if err != nil {
		return err
	}
	return nil
}

func RSAEncryptBytes(data []byte, publicKey []byte) ([]byte, error) {
	k, err := decodeRSAPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	return rsa.EncryptOAEP(sha512.New(), rand.Reader, k, data, nil)
}

func RSADecryptBytes(data []byte, privateKey []byte) ([]byte, error) {
	k, err := decodeRSAPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	plaintext, err := rsa.DecryptOAEP(sha512.New(), rand.Reader, k, data, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func GenerateRSAKeyPair(bits int) ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}
	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	publicKeyPEM, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	publicKeyPEMBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyPEM,
	}
	return pem.EncodeToMemory(publicKeyPEMBlock), pem.EncodeToMemory(privateKeyPEM), nil
}

func decodePEM(key []byte) (*pem.Block, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}
	return block, nil
}

func decodeRSAPublicKey(key []byte) (*rsa.PublicKey, error) {
	block, err := decodePEM(key)
	if err != nil {
		return nil, err
	}
	derDecodedKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return derDecodedKey, nil
	}
	pemDecodedKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := pemDecodedKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to decode RSA public key")
	}
	return rsaKey, nil
}

func decodeRSAPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, err := decodePEM(key)
	if err != nil {
		return nil, err
	}
	derDecodedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return derDecodedKey, nil
	}
	pemDecodedKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	rsaKey, ok := pemDecodedKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to decode RSA private key")
	}
	return rsaKey, nil
}
