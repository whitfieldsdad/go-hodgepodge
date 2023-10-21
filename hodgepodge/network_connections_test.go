package hodgepodge

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestListNetworkConnections(t *testing.T) {
	connections, err := ListNetworkConnections()
	assert.Nil(t, err, "ListNetworkConnections() failed")
	assert.Greater(t, len(connections), 0, "Expected at least one network connection")
	log.Infof("Found %d network connections", len(connections))
}
