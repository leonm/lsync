package main

import "testing"
import "time"

func TestScanRegularFiles(t *testing.T) {
	rootDirectory := RootDirectory{"/tmp", newStubWalker([]string{"/tmp/test1", "/tmp/test2"}, true), newCalculateHash(33)}
	scannedFiles := make(chan *FileEntry, 100)
	rootDirectory.scan(scannedFiles)
	if (<-scannedFiles).Path != "test1" {
		t.Errorf("should contain test1")
	}
	if (<-scannedFiles).Path != "test2" {
		t.Errorf("should contain test2")
	}
}

func TestScanNonRegularFiles(t *testing.T) {
	rootDirectory := RootDirectory{"/tmp", newStubWalker([]string{"/tmp/test1", "/tmp/test2"}, false), newCalculateHash(33)}
	scannedFiles := make(chan *FileEntry)
	rootDirectory.scan(scannedFiles)
	if (<-scannedFiles) != nil {
		t.Errorf("should not contain non regular files")
	}
}

func TestInitialHashShouldBe0(t *testing.T) {
	rootDirectory := RootDirectory{"/tmp", newStubWalker([]string{"/tmp/test1"}, true), newCalculateHash(33)}
	scannedFiles := make(chan *FileEntry, 100)
	rootDirectory.scan(scannedFiles)
	if (<-scannedFiles).Hash != 0 {
		t.Errorf("hash should be 0")
	}
}

func TestHashNewFiles(t *testing.T) {
	rootDirectory := RootDirectory{"/tmp", newStubWalker([]string{"/tmp/test1"}, true), newCalculateHash(33)}
	scannedFiles := make(chan *FileEntry, 100)
	hashedFiles := make(chan *FileEntry, 100)
	scannedFiles <- &FileEntry{"test1", 1, time.Now(), 0}
	go rootDirectory.hash(scannedFiles, hashedFiles)
	if (<-hashedFiles).Hash != 33 {
		t.Errorf("hash should be 33")
	}
}

// Mocks

func newStubWalker(paths []string, regular bool) func(root string, walkFn WalkFunc) error {
	return func(root string, walkFn WalkFunc) error {
		for _, path := range paths {
			walkFn(path, &FileInfo{regular, time.Now(), 1234}, nil)
		}
		return nil
	}
}

func newCalculateHash(value uint64) func(filename string) uint64 {
	return func(filename string) uint64 {
		return value
	}
}
