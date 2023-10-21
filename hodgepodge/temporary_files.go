package hodgepodge

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

func NewTempFile(data []byte) (string, error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return "", err
	}
	path := f.Name()
	if data != nil {
		_, err = f.Write(data)
		if err != nil {
			return "", err
		}
		err = f.Close()
		if err != nil {
			return path, err
		}
	}
	return path, nil
}

func GetTempDir() (string, error) {
	return os.TempDir(), nil
}

func NewTempDir() (string, error) {
	return os.MkdirTemp("", "")
}

func DeleteFile(path string) error {
	dir, err := IsDirectory(path)
	if err != nil {
		return errors.Wrap(err, "failed to check if path is a directory")
	}
	if dir {
		return os.RemoveAll(path)
	}
	return os.Remove(path)
}

func MoveFile(src, dst string) error {
	return os.Rename(src, dst)
}

func CopyFile(src, dst string) error {
	dir, err := IsDirectory(src)
	if err != nil {
		return err
	}
	if dir {
		return copyDirectory(src, dst)
	} else {
		return copyFile(src, dst)
	}
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return errors.Wrap(err, "failed to open source file")
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return errors.Wrap(err, "failed to open destination file")
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return errors.Wrap(err, "failed to copy file")
	}
	return nil
}

func copyDirectory(src, dst string) error {
	srcFiles, err := os.ReadDir(src)
	if err != nil {
		return errors.Wrap(err, "failed to read source directory")
	}
	err = os.MkdirAll(dst, 0755)
	if err != nil {
		return errors.Wrap(err, "failed to create destination directory")
	}
	for _, file := range srcFiles {
		srcPath := src + "/" + file.Name()
		dstPath := dst + "/" + file.Name()
		if file.IsDir() {
			err = copyDirectory(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = copyFile(srcPath, dstPath)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
