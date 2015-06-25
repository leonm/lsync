package main

import "os"
import "io"
import "path/filepath"
import "time"
import "hash/fnv"

func check(e error) {
  if e != nil {
    panic(e)
  }
}

type FileEntry struct {
  Path string
  Size int64
  Updated time.Time
  Hash uint64
}

type FileInfo struct {
  regular bool
  modTime time.Time
  size int64
}

type RootDirectory struct {
  root string
  walk func (root string, walkFn WalkFunc) error
}

type WalkFunc func(path string, info *FileInfo, err error) error

func NewRootDirectory(root string) *RootDirectory {
  return &RootDirectory {root,walk}
}

func walk(root string, walkFn WalkFunc) error {
  return filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
    return walkFn(root, &FileInfo {f.Mode().IsRegular(), f.ModTime(), f.Size()}, err)
  })
}

func (rootDirectory *RootDirectory) scan (out chan *FileEntry) {
  err := walk(rootDirectory.root, func(path string, f *FileInfo, err error) error {
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

func hash(root string, in chan *FileEntry, out chan *FileEntry) {
  hasher := fnv.New64a()

  for f := range in {
    hasher.Reset()
    file, err := os.Open(filepath.Join(root,f.Path))
    check(err)
    io.Copy(hasher,file)
    file.Close()
    f.Hash = hasher.Sum64()
    out <- f
  }
  close(out)
}
