package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEnvironmentVariables(t *testing.T) {
	env := GetEnvironmentVariables()
	assert.NotEmpty(t, env, "GetEnvironmentVariables should return at least one environment variable")
}
