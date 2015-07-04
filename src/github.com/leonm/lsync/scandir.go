package main

import "os"
import "io"
import "path/filepath"
import "time"
import "hash/fnv"


type FileInfo struct {
  regular bool
  modTime time.Time
  size int64
}

type RootDirectory struct {
  root string
  walk func (root string, walkFn WalkFunc) error
  calculateHash func(filename string) uint64
}

type WalkFunc func(path string, info *FileInfo, err error) error

func NewRootDirectory(root string) *RootDirectory {
  return &RootDirectory {root,walk,calculateHash}
}

func walk(root string, walkFn WalkFunc) error {
  return filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
    return walkFn(path, &FileInfo {f.Mode().IsRegular(), f.ModTime(), f.Size()}, err)
  })
}

func calculateHash (filename string) uint64 {
  hasher := fnv.New64a()
  hasher.Reset()
  file, err := os.Open(filename)
  check(err)
  io.Copy(hasher,file)
  file.Close()
  return hasher.Sum64()
}

func (rootDirectory *RootDirectory) scan(out chan *FileEntry) {
  err := rootDirectory.walk(rootDirectory.root, func(path string, f *FileInfo, err error) error {
    check(err)
    if (f.regular) {
      relativePath, err := filepath.Rel(rootDirectory.root,path)
      check(err)
      out <- &FileEntry{relativePath, f.size, f.modTime, 0}
    }
    return nil
  })
  close(out)
  check(err)

}

func (rootDirectory *RootDirectory) hash(in chan *FileEntry, out chan *FileEntry) {

  for f := range in {
    f.Hash = rootDirectory.calculateHash(filepath.Join(rootDirectory.root,f.Path))
    out <- f
  }
  close(out)
}
