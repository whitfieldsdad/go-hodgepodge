package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindPaths(t *testing.T) {
	tmpDir, err := GetTempDir()
	assert.Nil(t, err)

	results := make(chan string)
	search := FileSearch{
		Roots: []string{tmpDir},
	}
	go search.FindPaths(results)

	total := 0
	for range results {
		total++
	}
	assert.Greater(t, total, 0)
}

func TestFindFiles(t *testing.T) {
	tmpDir, err := GetTempDir()
	assert.Nil(t, err)

	results := make(chan File)
	search := FileSearch{
		Roots: []string{tmpDir},
	}
	go search.FindFiles(results, nil)

	total := 0
	for range results {
		total++
	}
	assert.Greater(t, total, 0)
}

func TestFindFilesWithOptions(t *testing.T) {
	tmpDir, err := GetTempDir()
	assert.Nil(t, err)

	results := make(chan File)
	search := FileSearch{
		Roots: []string{tmpDir},
	}
	opts := &FileMetadataOptions{
		IncludeFileTimestamps: false,
		IncludeFileTraits:     false,
		IncludeFileHashes:     false,
	}
	go search.FindFiles(results, opts)

	total := 0
	for range results {
		total++
	}
	assert.Greater(t, total, 0)
}
