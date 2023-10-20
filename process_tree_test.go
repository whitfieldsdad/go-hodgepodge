package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProcessTree(t *testing.T) {
	_, err := GetProcessTree()
	assert.Nil(t, err, "GetProcessTree should never return an error")
}
