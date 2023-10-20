package main

import (
	"os"
	"path/filepath"
	"strings"
)

// FileExists returns true if the file exists.
func FileExists(path string) (bool, error) {
	path, err := RealPath(path)
	if err != nil {
		return false, err
	}
	_, err = os.Stat(path)
	return !os.IsNotExist(err), nil
}

func SplitPath(path string) []string {
	dir, last := filepath.Split(path)
	if dir == "" {
		return []string{last}
	}
	dir = filepath.Clean(dir)
	return append(SplitPath(dir), last)
}

// Resolve resolves a list of paths by expanding environment variables, expanding ~, resolving symlinks.
func RealPaths(paths []string) ([]string, error) {
	var resolved []string
	for _, path := range paths {
		path, err := RealPath(path)
		if err != nil {
			return nil, err
		}
		resolved = append(resolved, path)
	}
	return resolved, nil
}

func RealPath(path string) (string, error) {

	// Expand environment variables.
	path = os.ExpandEnv(path)

	// Expand tilde (~).
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		path = strings.Replace(path, "~", home, 1)
	}

	// Resolve relative paths.
	path, _ = filepath.Abs(path)

	// Resolve symlinks.
	path, err := filepath.EvalSymlinks(path)
	if err != nil {
		return "", err
	}

	// Clean up the path.
	path = filepath.Clean(path)
	return path, nil
}
