package main

import "time"

type FileEntry struct {
  Path string
  Size int64
  Updated time.Time
  Hash uint64
}

type FileList struct {
  Current []FileEntry
  New []FileEntry
}

func NewFileList() *FileList {
  return &FileList{[]FileEntry{},[]FileEntry{}}
}

func (list *FileList) update (in chan *FileEntry) {

  list.New = []FileEntry{}
  working := []FileEntry{}

  for f := range in {
    list.New = append(list.New, *f)
    working = append(working, *f)
  }

  list.Current = working
}
