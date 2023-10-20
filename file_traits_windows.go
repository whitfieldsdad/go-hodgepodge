//go:build windows

package main

import (
	"path/filepath"
	"syscall"
)

func isHidden(path string) (bool, error) {
	path, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return false, err
	}
	pointer, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, err
	}
	attributes, err := syscall.GetFileAttributes(pointer)
	if err != nil {
		return false, err
	}
	return attributes&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}
