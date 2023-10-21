package hodgepodge

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"

	"github.com/charmbracelet/log"
)

func CreateTarballFile(archivePath string, inputFiles []string) error {
	log.Infof("Creating archive: %s", archivePath)
	file, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	log.Infof("Adding %d files to archive: %s", len(inputFiles), archivePath)
	err = CreateTarball(file, inputFiles)
	if err != nil {
		return err
	}
	log.Infof("Added %d files to archive: %s", len(inputFiles), archivePath)
	return nil
}

func CreateTarball(buf io.Writer, inputFiles []string) error {
	gw := gzip.NewWriter(buf)
	defer gw.Close()

	tw := tar.NewWriter(gw)
	defer tw.Close()

	for _, path := range inputFiles {
		log.Debugf("Adding file to archive: %s", path)
		err := addFileToTarWriter(tw, path)
		if err != nil {
			return err
		}
	}
	return nil
}

func addFileToTarWriter(tw *tar.Writer, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Each entry contains information about the file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}
	header.Name = path
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy the file to the archive.
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}
	return nil
}
