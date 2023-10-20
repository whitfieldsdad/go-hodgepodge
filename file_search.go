package main

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/exp/slices"
)

type FileFilter struct {
	Filenames []string `json:"filenames"`
	Paths     []string `json:"paths"`
}

// Matches is a boolean function that returns true if the provided file matches the filter. If any errors arise, they will be silently ignored.
func (f FileFilter) Matches(path string, info *os.FileInfo) bool {
	if len(f.Filenames) > 0 {
		if !hasMatchingFilename(path, f.Filenames) {
			return false
		}
	}
	if len(f.Paths) > 0 {
		if !hasMatchingPath(path, f.Paths) {
			return false
		}
	}
	return true
}

func hasMatchingFilename(path string, filenames []string) bool {
	filename := filepath.Base(path)
	return slices.Contains(filenames, filename)
}

func hasMatchingPath(path string, paths []string) bool {
	for _, p := range paths {
		if strings.HasPrefix(path, p) {
			return true
		}
	}
	return false
}

type FileSearch struct {
	Roots       []string
	FileFilters []FileFilter `json:"file_filters"`
}

func (s FileSearch) matchesAnyFilter(root string, info os.FileInfo) bool {
	for _, filter := range s.FileFilters {
		if filter.Matches(root, &info) {
			return true
		}
	}
	return false
}

func (s FileSearch) FindFiles(results chan<- File, opts *FileOptions) {
	defer close(results)

	paths := make(chan string)
	go s.FindPaths(paths)

	for path := range paths {
		file, err := GetFile(path, opts)
		if err != nil {
			continue
		}
		results <- *file
	}
}

func (s FileSearch) FindPaths(results chan<- string) {
	defer close(results)

	roots := ReducePaths(s.Roots)
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if len(s.FileFilters) > 0 {
				if !s.matchesAnyFilter(root, info) {
					continue
				}
			}
			results <- root
		} else {
			_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if len(s.FileFilters) > 0 {
					if !s.matchesAnyFilter(root, info) {
						return nil
					}
				}
				results <- path
				return nil
			})
		}
	}
}

// ReducePaths removes any paths that are subpaths of other paths (e.g. /usr, /usr/local -> /usr)
func ReducePaths(paths []string) []string {
	result := []string{}
	for i := 0; i < len(paths); i++ {
		includePath := true
		for j := 0; j < len(paths); j++ {
			if i != j && strings.HasPrefix(paths[i], paths[j]) {
				includePath = false
				break
			}
		}
		if includePath {
			result = append(result, paths[i])
		}
	}
	return result
}
