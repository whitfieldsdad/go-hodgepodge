package main

import (
	"crypto/rand"

	"github.com/pkg/errors"
)

func GenerateSalt(saltLengthBytes int) ([]byte, error) {
	salt := make([]byte, saltLengthBytes)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate salt")
	}
	return salt, nil
}
