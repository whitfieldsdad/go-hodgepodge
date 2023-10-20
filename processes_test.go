package main

import (
	"testing"

	"github.com/charmbracelet/log"
	"github.com/stretchr/testify/assert"
)

func TestListProcesses(t *testing.T) {
	opts := &FileMetadataOptions{
		IncludeFileSize:       true,
		IncludeFileHashes:     true,
		IncludeFileTimestamps: true,
		IncludeFileTraits:     true,
	}
	processes, err := ListProcesses(opts)
	assert.Nil(t, err, "ListProcesses() failed")
	assert.Greater(t, len(processes), 0, "Expected at least one process")
	log.Infof("Found %d processes", len(processes))
}
