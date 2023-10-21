package hodgepodge

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestListNetworkInterfaces(t *testing.T) {
	ifaces, err := ListNetworkInterfaces()
	assert.Nil(t, err, "ListNetworkInterfaces() failed")
	assert.Greater(t, len(ifaces), 0, "Expected at least one network interface")
	log.Infof("Found %d network interfaces", len(ifaces))
}
