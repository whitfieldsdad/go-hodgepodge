package hodgepodge

const (
	IncludeFileTraits     = false
	IncludeFileHashes     = false
	IncludeFileTimestamps = false
)

// FileOptions contains options for collecting file metadata.
type FileOptions struct {
	IncludeFileHashes     bool `json:"include_file_hashes"`     // If true, collect file hashes (i.e. MD5, SHA1, etc.).
	IncludeFileTimestamps bool `json:"include_file_timestamps"` // If true, collect file timestamps (MACb).
	IncludeFileTraits     bool `json:"include_file_traits"`     // If true, collect file traits (i.e. file type, permissions, etc.)
}

func GetDefaultFileOptions() *FileOptions {
	return &FileOptions{
		IncludeFileTimestamps: IncludeFileTimestamps,
		IncludeFileTraits:     IncludeFileTraits,
		IncludeFileHashes:     IncludeFileHashes,
	}
}
