package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/whitfieldsdad/hodgepodge/pkg/crypto"
)

func TestDownloadFile(t *testing.T) {
	f, err := os.CreateTemp("", "sample")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	path := f.Name()
	err = DownloadFile("https://upload.wikimedia.org/wikipedia/commons/8/8f/Example_image.svg", path)
	assert.Nil(t, err)

	h, err := crypto.GetFileSHA256(path)
	assert.Nil(t, err)
	assert.Equal(t, "63f33c7a2a11d0993c3f4150b8dc7e0335cc00ede829197de93a32cea494b2ef", h)
}

func TestDownloadFileIntoMemory(t *testing.T) {
	data, err := DownloadFileIntoMemory("https://upload.wikimedia.org/wikipedia/commons/8/8f/Example_image.svg")
	if err != nil {
		t.Fatal(err.Error())
	}
	assert.Equal(t, "63f33c7a2a11d0993c3f4150b8dc7e0335cc00ede829197de93a32cea494b2ef", GetSHA256(data))
}
