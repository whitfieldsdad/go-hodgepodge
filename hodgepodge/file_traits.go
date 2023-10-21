package hodgepodge

import (
	"os"

	"github.com/charmbracelet/log"
	"golang.org/x/exp/slices"
)

type FileTrait string

const (
	Directory       FileTrait = "directory"
	RegularFile     FileTrait = "regular_file"
	SymbolicLink    FileTrait = "symbolic_link"
	Socket          FileTrait = "socket"
	HardLink        FileTrait = "hard_link"
	NamedPipe       FileTrait = "named_pipe"
	BlockDevice     FileTrait = "block_device"
	CharacterDevice FileTrait = "character_device"
	Hidden          FileTrait = "hidden"
)

type FileTraits struct {
	IsDirectory       bool `json:"is_directory"`
	IsRegularFile     bool `json:"is_regular_file"`
	IsSymbolicLink    bool `json:"is_symbolic_link"`
	IsSocket          bool `json:"is_socket"`
	IsHardLink        bool `json:"is_hard_link"`
	IsNamedPipe       bool `json:"is_named_pipe"`
	IsBlockDevice     bool `json:"is_block_device"`
	IsCharacterDevice bool `json:"is_character_device"`
	IsHidden          bool `json:"is_hidden"`
}

func (traits FileTraits) ToList() []FileTrait {
	checks := map[FileTrait]bool{
		BlockDevice:     traits.IsBlockDevice,
		CharacterDevice: traits.IsCharacterDevice,
		Directory:       traits.IsDirectory,
		HardLink:        traits.IsHardLink,
		Hidden:          traits.IsHidden,
		NamedPipe:       traits.IsNamedPipe,
		RegularFile:     traits.IsRegularFile,
		Socket:          traits.IsSocket,
		SymbolicLink:    traits.IsSymbolicLink,
	}
	list := make([]FileTrait, 0)
	for trait, check := range checks {
		if check {
			list = append(list, trait)
		}
	}
	slices.Sort(list)
	return list
}

func GetFileTraits(path string) (*FileTraits, error) {
	st, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	traits := &FileTraits{
		IsBlockDevice:     isBlockDevice(st),
		IsCharacterDevice: isCharacterDevice(st),
		IsDirectory:       isDirectory(st),
		IsHardLink:        isHardLink(st),
		IsNamedPipe:       isNamedPipe(st),
		IsRegularFile:     isRegularFile(st),
		IsSocket:          isSocket(st),
		IsSymbolicLink:    isSymbolicLink(st),
	}
	hidden, err := IsHidden(path)
	if err != nil {
		log.Warnf("Failed to check if file is hidden: %s - %v", path, err)
	} else {
		traits.IsHidden = hidden
	}
	return traits, nil
}

// IsHidden returns true if the file is hidden.
func IsHidden(path string) (bool, error) {
	return isHidden(path)
}

func hasTrait(path string, oracle func(os.FileInfo) bool) (bool, error) {
	st, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return oracle(st), nil
}

// IsDirectory returns true if the file is a directory.
func IsDirectory(path string) (bool, error) {
	return hasTrait(path, isDirectory)
}

// IsRegularFile returns true if the file is a regular file.
func IsRegularFile(path string) (bool, error) {
	return hasTrait(path, isRegularFile)
}

// IsSymbolicLink returns true if the file is a symbolic link.
func IsSymbolicLink(path string) (bool, error) {
	return hasTrait(path, isSymbolicLink)
}

// IsHardLink returns true if the file is a hard link.
func IsHardLink(path string) (bool, error) {
	return hasTrait(path, isHardLink)
}

// IsSocket returns true if the file is a socket.
func IsSocket(path string) (bool, error) {
	return hasTrait(path, isSocket)
}

// IsNamedPipe returns true if the file is a named pipe.
func IsNamedPipe(path string) (bool, error) {
	return hasTrait(path, isNamedPipe)
}

// IsBlockDevice returns true if the file is a block device.
func IsBlockDevice(path string) (bool, error) {
	return hasTrait(path, isBlockDevice)
}

// IsCharacterDevice returns true if the file is a character device.
func IsCharacterDevice(path string) (bool, error) {
	return hasTrait(path, isCharacterDevice)
}

func isDirectory(st os.FileInfo) bool {
	return st.Mode().IsDir()
}

func isRegularFile(st os.FileInfo) bool {
	return st.Mode().IsRegular()
}

func isSymbolicLink(st os.FileInfo) bool {
	return st.Mode()&os.ModeSymlink != 0
}

func isHardLink(st os.FileInfo) bool {
	return st.Mode()&os.ModeDevice != 0
}

func isSocket(st os.FileInfo) bool {
	return st.Mode()&os.ModeSocket != 0
}

func isNamedPipe(st os.FileInfo) bool {
	return st.Mode()&os.ModeNamedPipe != 0
}

func isBlockDevice(st os.FileInfo) bool {
	return st.Mode()&os.ModeDevice != 0
}

func isCharacterDevice(st os.FileInfo) bool {
	return st.Mode()&os.ModeCharDevice != 0
}
