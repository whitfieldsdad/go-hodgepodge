package main

import "errors"

const (
	AES256GCM                           = "aes-256-gcm"
	ChaCha20Poly1305                    = "chacha20-poly1305"
	DefaultSymmetricEncryptionAlgorithm = AES256GCM
)

func SymmetricEncryptFile(path, password, algorithm string) error {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMEncryptFile(path, password)
	} else if algorithm == ChaCha20Poly1305 {
		return ChaCha20Poly1305EncryptFile(path, password)
	} else {
		return errors.New("unsupported algorithm")
	}
}

func SymmetricDecryptFile(path, password, algorithm string) error {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMDecryptFile(path, password)
	} else if algorithm == ChaCha20Poly1305 {
		return ChaCha20Poly1305DecryptFile(path, password)
	} else {
		return errors.New("unsupported algorithm")
	}
}

func SymmetricEncryptBytes(data []byte, password, algorithm string) ([]byte, error) {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMEncryptBytes(data, password)
	} else if algorithm == ChaCha20Poly1305 {
		return ChaCha20Poly1305EncryptBytes(data, password)
	} else {
		return nil, errors.New("unsupported algorithm")
	}
}

func SymmetricDecryptBytes(data []byte, password, algorithm string) ([]byte, error) {
	if algorithm == "" {
		algorithm = DefaultSymmetricEncryptionAlgorithm
	}
	if algorithm == AES256GCM {
		return AES256GCMDecryptBytes(data, password)
	} else if algorithm == ChaCha20Poly1305 {
		return ChaCha20Poly1305DecryptBytes(data, password)
	} else {
		return nil, errors.New("unsupported algorithm")
	}
}

// TODO
func AES256GCMEncryptFile(path, password string) error {
	return nil
}

// TODO
func AES256GCMDecryptFile(path, password string) error {
	return nil
}

// TODO
func AES256GCMEncryptBytes(data []byte, password string) ([]byte, error) {
	return nil, nil
}

// TODO
func AES256GCMDecryptBytes(data []byte, password string) ([]byte, error) {
	return nil, nil
}

// TODO
func ChaCha20Poly1305EncryptFile(path, password string) error {
	return nil
}

// TODO
func ChaCha20Poly1305DecryptFile(path, password string) error {
	return nil
}

// TODO
func ChaCha20Poly1305EncryptBytes(data []byte, password string) ([]byte, error) {
	return nil, nil
}

// TODO
func ChaCha20Poly1305DecryptBytes(data []byte, password string) ([]byte, error) {
	return nil, nil
}