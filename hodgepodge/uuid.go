package hodgepodge

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/oklog/ulid"
)

// NewUUID4 returns a type 4 UUID.
func NewUUID4() string {
	return uuid.New().String()
}

// NewUUID5 returns a type 5 UUID based using the provided namespace and name.
func NewUUID5(namespace uuid.UUID, name string) string {
	return uuid.NewSHA1(namespace, []byte(name)).String()
}

// NewULID returns a Universally Unique Lexicographically Sortable Identifier (ULID).
func NewULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	return ulid.MustNew(ms, entropy).String()
}
