package hodgepodge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHost(t *testing.T) {
	host, err := GetHost()
	assert.Nil(t, err, "GetHost() should not return an error")
	assert.NotNil(t, host, "GetHost() should return a host")
}
