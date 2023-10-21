package hodgepodge

import (
	"crypto/sha256"

	"golang.org/x/crypto/pbkdf2"
)

const (
	PBKDF2Rounds     = 100000
	PBKDF2           = "pbkdf2"
	PBKDF2SaltLength = 8
)

func PBKDF2_HMAC_SHA256(password string, salt []byte, rounds int) []byte {
	return pbkdf2.Key([]byte(password), salt, rounds, 32, sha256.New)
}
