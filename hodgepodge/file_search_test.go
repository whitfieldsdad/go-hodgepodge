package hodgepodge

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
	go FindPaths(search, results)

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
	go FindFiles(search, results, nil)

	total := 0
	for range results {
		total++
	}
	assert.Greater(t, total, 0)
}
