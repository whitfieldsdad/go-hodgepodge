package hodgepodge

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

func FindFiles(search FileSearch, results chan<- File, opts *FileOptions) {
	defer close(results)

	paths := make(chan string)
	go FindPaths(search, paths)

	for path := range paths {
		file, err := GetFile(path, opts)
		if err != nil {
			continue
		}
		results <- *file
	}
}

func FindPaths(search FileSearch, results chan<- string) {
	defer close(results)

	roots := ReducePaths(search.Roots)
	for _, root := range roots {
		info, err := os.Stat(root)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			if len(search.FileFilters) > 0 {
				matches := false
				for _, filter := range search.FileFilters {
					if filter.Matches(root, &info) {
						matches = true
						break
					}
				}
				if !matches {
					continue
				}
			}
			results <- root
		} else {
			_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
				if len(search.FileFilters) > 0 {
					matches := false
					for _, filter := range search.FileFilters {
						if filter.Matches(path, &info) {
							matches = true
							break
						}
					}
					if !matches {
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
