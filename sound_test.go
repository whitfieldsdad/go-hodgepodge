package main

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/whitfieldsdad/hodgepodge/pkg/crypto"
)

func TestPlayLocalMP3(t *testing.T) {
	f, err := os.CreateTemp("", "test.mp3")
	assert.Nil(t, err)

	path := f.Name()
	defer os.Remove(path)

	// Close the file.
	err = f.Close()
	assert.Nil(t, err)

	// Download the example MP3 file.
	err = downloadExampleMP3(path)
	assert.Nil(t, err)

	// Confirm that the example MP3 file has the expected SHA-256 file hash.
	h, err := crypto.GetFileSHA256(path)
	assert.Nil(t, err)
	assert.Equal(t, "4216d76339f030f3e0e7311f23da06e4bd105136ee25f0806fe34e5af5668301", h)

	// Let the MP3 play for 1 second before stopping it.
	ctx, cancel := context.WithCancel(context.Background())
	playing := make(chan interface{})
	go PlayMP3FileWithCallback("https://bit.ly/3RE8BGk", ctx, playing, nil)
	<-playing

	select {
	case <-time.After(1 * time.Second):
		cancel()
	}
}

func TestPlayLocalMP3WithStartCallback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	playing := make(chan interface{})
	go PlayMP3FileWithCallback("https://bit.ly/3RE8BGk", ctx, playing, nil)
	<-playing

	// Let the MP3 play for 1 second before stopping it.
	select {
	case <-time.After(1 * time.Second):
		cancel()
	}
}

func TestPlayRemoteMP3WithStartCallback(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	playing := make(chan interface{})
	go PlayMP3FileWithCallback("https://bit.ly/3RE8BGk", ctx, playing, nil)
	<-playing

	// Let the MP3 play for 1 second before stopping it.
	select {
	case <-time.After(1 * time.Second):
		cancel()
	}
}

func downloadExampleMP3(path string) error {
	return DownloadFile("https://bit.ly/3RE8BGk", path)
}
