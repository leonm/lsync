package main

import "time"
import "path/filepath"

type FileEntry struct {
	Path    string
	Size    int64
	Updated time.Time
	Hash    uint64
}

type FileList struct {
	Current []FileEntry
	New     []FileEntry
}

func NewFileList() *FileList {
	return &FileList{[]FileEntry{}, []FileEntry{}}
}

func (list *FileList) update(in chan *FileEntry) {

	list.New = []FileEntry{}
	working := []FileEntry{}

	for f := range in {
		list.New = append(list.New, *f)
		working = append(working, *f)
	}

	list.Current = working
}

func (fileEntry *FileEntry) IsUptoDate(t time.Time, size int64) bool {
	return fileEntry.Updated.Equal(t) && fileEntry.Size == size
}

func (fileEntry *FileEntry) Location(rootPath string) string {
	return filepath.Join(rootPath, fileEntry.Path)
}

func (fileEntry *FileEntry) TempLocation(rootPath string) string {
	return fileEntry.Location(rootPath) + ".part"
}
