package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetFileMD5(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	result, err := GetFileMD5(tmp)
	if err != nil {
		t.Fatalf("Failed to hash file: %s", err)
	}
	assertEqual(t, TestStringMD5, result)
}

func TestGetFileSHA1(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	result, err := GetFileSHA1(tmp)
	if err != nil {
		t.Fatalf("Failed to hash file: %s", err)
	}
	assertEqual(t, TestStringSHA1, result)
}

func TestGetFileSHA256(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	result, err := GetFileSHA256(tmp)
	if err != nil {
		t.Fatalf("Failed to hash file: %s", err)
	}
	assertEqual(t, TestStringSHA256, result)
}

func TestGetFileSHA512(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	result, err := GetFileSHA512(tmp)
	if err != nil {
		t.Fatalf("Failed to hash file: %s", err)
	}
	assertEqual(t, TestStringSHA512, result)
}

func TestGetFileXXH64(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	result, err := GetFileXXH64(tmp)
	if err != nil {
		t.Fatalf("Failed to hash file: %s", err)
	}
	assertEqual(t, TestStringXXH64, result)
}

func TestGetFileHashes(t *testing.T) {
	tmp, err := NewTempFile([]byte(TestString))
	if err != nil {
		t.Fatalf("Failed to create temporary file: %s", err)
	}
	defer DeleteFile(tmp)

	expected := &Hashes{
		MD5:    TestStringMD5,
		SHA1:   TestStringSHA1,
		SHA256: TestStringSHA256,
		SHA512: TestStringSHA512,
		XXH64:  TestStringXXH64,
	}
	result, err := GetFileHashes(tmp)
	if err != nil {
		t.Fatalf("Failed to calculate file hashes: %s", err)
	}
	assertEqual(t, expected, result)
}

func TestGetMD5(t *testing.T) {
	result := GetMD5([]byte(TestString))
	assertEqual(t, TestStringMD5, result)
}

func TestGetSHA1(t *testing.T) {
	result := GetSHA1([]byte(TestString))
	assertEqual(t, TestStringSHA1, result)
}

func TestGetSHA256(t *testing.T) {
	result := GetSHA256([]byte(TestString))
	assertEqual(t, TestStringSHA256, result)
}

func TestGetSHA512(t *testing.T) {
	result := GetSHA256([]byte(TestString))
	assertEqual(t, TestStringSHA512, result)
}

func TestGetXXH64(t *testing.T) {
	result := GetXXH64([]byte(TestString))
	assertEqual(t, TestStringXXH64, result)
}

func assertEqual(t *testing.T, expected, result interface{}) {
	if !cmp.Equal(expected, result) {
		t.Fatalf("TestString %s got %s", expected, result)
	}
}
