package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFile = "../../tools/windows/calc.exe"
)

func TestGetFileMetadata(t *testing.T) {
	exists, err := Exists(testFile)
	assert.True(t, exists)
	assert.Nil(t, err)

	meta, err := GetFileMetadata(testFile, nil)
	assert.Nil(t, err)
	assert.NotNil(t, meta)

	if IncludeFileTraits {
		assert.NotNil(t, meta.Traits)
	} else {
		assert.Nil(t, meta.Traits)
	}
	if IncludeFileTimestamps {
		assert.NotNil(t, meta.Timestamps)
	} else {
		assert.Nil(t, meta.Timestamps)
	}
}

func TestGetFileMetadataWithOptions(t *testing.T) {
	opts := GetDefaultFileMetadataOptions()
	opts.IncludeFileTimestamps = false
	opts.IncludeFileTraits = false

	// Verify that we can disable timestamps and traits.
	meta, err := GetFileMetadata(testFile, opts)
	assert.Nil(t, err)
	assert.NotNil(t, meta)
	assert.Nil(t, meta.Traits)
	assert.Nil(t, meta.Timestamps)
}

func TestGetFileTimestamps(t *testing.T) {
	timestamps, err := GetFileTimestamps(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, timestamps)
}

func TestGetFileTraits(t *testing.T) {
	expected := &FileTraits{
		IsDirectory:       false,
		IsRegularFile:     true,
		IsSymbolicLink:    false,
		IsSocket:          false,
		IsHardLink:        false,
		IsNamedPipe:       false,
		IsBlockDevice:     false,
		IsCharacterDevice: false,
		IsHidden:          false,
	}
	result, err := GetFileTraits(testFile)
	assert.Nil(t, err)
	assert.Equal(t, expected, result)
}
