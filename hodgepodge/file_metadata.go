package hodgepodge

import (
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
	"github.com/djherbis/times"
)

// FileTimestamps contains the MACb timestamps of a file.
type FileTimestamps struct {
	ModifyTime *time.Time `json:"modify_time"`
	AccessTime *time.Time `json:"access_time"`
	ChangeTime *time.Time `json:"change_time"`
	BirthTime  *time.Time `json:"birth_time"`
}

// File contains basic information about a file.
type File struct {
	Path       string          `json:"path"`
	Filename   string          `json:"filename"`
	Directory  string          `json:"directory"`
	Extension  string          `json:"extension"`
	Size       *int64          `json:"size"`
	Timestamps *FileTimestamps `json:"timestamps,omitempty"`
	Traits     *FileTraits     `json:"traits,omitempty"`
	Hashes     *Hashes         `json:"hashes,omitempty"`
}

// GetFile returns basic information about a file.
func GetFile(path string, opts *FileOptions) (*File, error) {
	m := &File{
		Path:      path,
		Filename:  filepath.Base(path),
		Directory: filepath.Dir(path),
		Extension: filepath.Ext(path),
	}
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	size := info.Size()
	m.Size = &size

	// Optional data collection.
	if opts == nil {
		opts = GetDefaultFileOptions()
	}
	if opts.IncludeFileTimestamps {
		timestamps, err := GetFileTimestamps(path)
		if err != nil {
			log.Errorf("Encountered error while identifying timestamps of %s - %s", path, err.Error())
		}
		m.Timestamps = timestamps
	}
	if opts.IncludeFileTraits {
		traits, err := GetFileTraits(path)
		if err != nil {
			log.Errorf("Encountered error while identifying traits of %s - %s", path, err.Error())
		}
		m.Traits = traits
	}
	if opts.IncludeFileHashes {
		hashes, err := GetFileHashes(path)
		if err != nil {
			log.Errorf("Encountered error while identifying hashes of %s - %s", path, err.Error())
		}
		m.Hashes = hashes
	}
	return m, nil
}

// GetFileTimestamps returns the MACb timestamps of a file.
func GetFileTimestamps(path string) (*FileTimestamps, error) {
	st, err := times.Stat(path)
	if err != nil {
		return nil, err
	}
	m := st.ModTime()
	a := st.AccessTime()
	timestamps := &FileTimestamps{
		ModifyTime: &m,
		AccessTime: &a,
	}
	if st.HasChangeTime() {
		c := st.ChangeTime()
		timestamps.ChangeTime = &c
	}
	if st.HasBirthTime() {
		b := st.BirthTime()
		timestamps.BirthTime = &b
	}
	return timestamps, nil
}
