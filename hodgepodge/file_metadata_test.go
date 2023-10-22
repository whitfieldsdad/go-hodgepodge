package hodgepodge

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testFile = "../tools/windows/calc.exe"
)

func TestGetFile(t *testing.T) {
	exists, err := FileExists(testFile)
	assert.Nil(t, err)
	assert.True(t, exists)

	meta, err := GetFile(testFile, nil)
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

func TestGetFileWithOptions(t *testing.T) {
	opts := GetDefaultFileOptions()
	opts.IncludeFileTimestamps = false
	opts.IncludeFileTraits = false

	// Verify that we can disable timestamps and traits.
	meta, err := GetFile(testFile, opts)
	assert.Nil(t, err)
	assert.NotNil(t, meta)
}

func TestGetFileTimestamps(t *testing.T) {
	timestamps, err := GetFileTimestamps(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, timestamps)
}

func TestGetFileTraits(t *testing.T) {
	expected := &FileTraits{
		IsBlockDevice:     false,
		IsCharacterDevice: false,
		IsDirectory:       false,
		IsHardLink:        false,
		IsHidden:          false,
		IsNamedPipe:       false,
		IsRegularFile:     true,
		IsSocket:          false,
		IsSymbolicLink:    false,
	}
	result, err := GetFileTraits(testFile)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expected, result)
}
