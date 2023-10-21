package hodgepodge

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"
	"strconv"

	"github.com/cespare/xxhash/v2"
	"github.com/charmbracelet/log"
	"github.com/dustin/go-humanize"
)

const (
	FileReadBufferSize int32 = 1000000
)

type Hashes struct {
	MD5    string `json:"md5"`
	SHA1   string `json:"sha1"`
	SHA256 string `json:"sha256"`
	SHA512 string `json:"sha512,omitempty"`
	XXH64  string `json:"xxh64"`
}

// GetFileHashes returns the MD5, SHA1, and SHA256 hashes of a file.
func GetFileHashes(path string) (*Hashes, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	info, err := f.Stat()
	if err != nil {
		return nil, err
	}
	sz := info.Size()
	log.Debugf("Hashing %s (size: %s)", path, humanize.Bytes(uint64(sz)))

	defer f.Close()
	return GetReaderHashes(f)
}

// GetReaderHashes returns the MD5, SHA1, and SHA256 hashes of a reader.
func GetReaderHashes(rd io.Reader) (*Hashes, error) {
	md5 := md5.New()
	sha1 := sha1.New()
	sha256 := sha256.New()
	sha512 := sha512.New()
	xxh64 := xxhash.New()

	pagesize := os.Getpagesize()
	reader := bufio.NewReaderSize(rd, pagesize)
	multiWriter := io.MultiWriter(md5, sha1, sha256, sha512, xxh64)
	_, err := io.Copy(multiWriter, reader)
	if err != nil {
		return nil, err
	}
	hashes := &Hashes{
		MD5:    fmt.Sprintf("%x", md5.Sum(nil)),
		SHA1:   fmt.Sprintf("%x", sha1.Sum(nil)),
		SHA256: fmt.Sprintf("%x", sha256.Sum(nil)),
		SHA512: fmt.Sprintf("%x", sha512.Sum(nil)),
		XXH64:  fmt.Sprintf("%x", xxh64.Sum(nil)),
	}
	return hashes, nil
}

func GetMD5(b []byte) string {
	h := md5.Sum(b)
	return fmt.Sprintf("%x", h)
}

func GetSHA1(b []byte) string {
	h := sha1.Sum(b)
	return fmt.Sprintf("%x", h)
}

func GetSHA256(b []byte) string {
	h := sha256.Sum256(b)
	return fmt.Sprintf("%x", h)
}

// GetSha256 returns the SHA256 hash of a byte slice.
func GetSHA512(b []byte) string {
	h := sha512.Sum512(b)
	return fmt.Sprintf("%x", h)
}

func GetXXH64(b []byte) string {
	h := xxhash.Sum64(b)
	return strconv.FormatUint(h, 16)
}

func GetFileMD5(path string) (string, error) {
	return getFileHash(path, md5.New())
}

func GetFileSHA1(path string) (string, error) {
	return getFileHash(path, sha1.New())
}

func GetFileSHA256(path string) (string, error) {
	return getFileHash(path, sha256.New())
}

func GetFileSHA512(path string) (string, error) {
	return getFileHash(path, sha512.New())
}

func GetFileXXH64(path string) (string, error) {
	return getFileHash(path, xxhash.New())
}

func getFileHash(path string, h hash.Hash) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}
