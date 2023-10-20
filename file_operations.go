package main

import (
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/log"
	"github.com/pkg/errors"
)

func DownloadFile(url string, path string) error {
	log.Infof("Downloading %s to %s", url, path)
	resp, err := http.Get(url)
	if err != nil {
		return errors.Wrap(err, "failed to open remote file")
	}
	defer resp.Body.Close()

	fileHandle, err := os.Create(path)
	if err != nil {
		return errors.Wrap(err, "failed to open local file")
	}
	defer fileHandle.Close()

	_, err = io.Copy(fileHandle, resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to download file")
	}
	log.Infof("Downloaded %s to %s", url, path)
	return nil
}

func DownloadFileIntoMemory(url string) ([]byte, error) {
	log.Infof("Downloading %s into memory", url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to download file")
	}
	log.Infof("Downloaded %s into memory", url)
	return data, err
}
