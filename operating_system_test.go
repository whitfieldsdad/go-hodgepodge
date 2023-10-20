package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetOperatingSystem(t *testing.T) {
	os := GetOperatingSystem()
	assert.NotEmpty(t, os.Type)
	assert.NotEmpty(t, os.Architecture)
	assert.NotEmpty(t, os.Name)
	assert.NotEmpty(t, os.Version)
	assert.NotEmpty(t, os.KernelVersion)
}
