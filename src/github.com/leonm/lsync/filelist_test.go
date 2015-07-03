package main

import "testing"
import "time"

func TestNewListGetsAddedToFirst(t *testing.T) {
  entries := make(chan *FileEntry, 100)
  entries <- &FileEntry{"test1",1,time.Now(),0}

  fileList := NewFileList()
  go fileList.update(entries)
  time.Sleep(100 * time.Millisecond)

  if (len(fileList.Current) > 0) { t.Errorf("current files should be empty") }
  if (len(fileList.New) != 1) { t.Errorf("new list should have an entry") }

  close(entries)
}

func TestCurrentListGetsAddedTo(t *testing.T) {
  entries := make(chan *FileEntry, 100)
  entries <- &FileEntry{"test1",1,time.Now(),0}
  close(entries)

  fileList := NewFileList()
  fileList.update(entries)

  if (len(fileList.Current) != 1) { t.Errorf("current list should have an entry") }
  if (len(fileList.New) != 1) { t.Errorf("new list should have an entry") }

}
