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


func scan (root string, out chan *FileEntry) {
  err := filepath.Walk(root, func(path string, f os.FileInfo, err error) error {
    check(err)
    if (f.Mode().IsRegular()) {
      relativePath, err := filepath.Rel(root,path)
      check(err)
      out <- &FileEntry{relativePath, f.Size(), f.ModTime(), 0}
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
