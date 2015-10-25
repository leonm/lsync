package main

import "testing"
import "time"

func TestNewListGetsAddedToFirst(t *testing.T) {
	entries := make(chan *FileEntry, 100)
	entries <- &FileEntry{"test1", 1, time.Now(), 0}

	fileList := NewFileList()
	go fileList.update(entries)
	time.Sleep(100 * time.Millisecond)

	if len(fileList.Current) > 0 {
		t.Errorf("current files should be empty")
	}
	if len(fileList.New) != 1 {
		t.Errorf("new list should have an entry")
	}

	close(entries)
}

func TestCurrentListGetsAddedTo(t *testing.T) {
	entries := make(chan *FileEntry, 100)
	entries <- &FileEntry{"test1", 1, time.Now(), 0}
	close(entries)

	fileList := NewFileList()
	fileList.update(entries)

	if len(fileList.Current) != 1 {
		t.Errorf("current list should have an entry")
	}
	if len(fileList.New) != 1 {
		t.Errorf("new list should have an entry")
	}

}

func TestFileEntryIsUpToDateWhenSizeAndTimeMatches(t *testing.T) {
	now := time.Now()
	entry := &FileEntry{"test1", 55, now, 0}
	if !entry.IsUptoDate(now, 55) {
		t.Errorf("file entry should be up to date")
	}
}

func TestFileEntryIsNotUpToDateWhenSizeDiffers(t *testing.T) {
	now := time.Now()
	entry := &FileEntry{"test1", 55, now, 0}
	if entry.IsUptoDate(now, 77) {
		t.Errorf("file entry should not be up to date")
	}
}

func TestFileEntryIsNotUpToDateWhenTimeDiffers(t *testing.T) {
	now := time.Now()
	entry := &FileEntry{"test1", 55, now, 0}
	if entry.IsUptoDate(time.Now(), 55) {
		t.Errorf("file entry should not be up to date")
	}
}

func TestFileEntryLocation(t *testing.T) {
	entry := &FileEntry{"test1", 55, time.Now(), 0}
	if entry.Location("/some/path") != "/some/path/test1" {
		t.Errorf("file entry location should be correct")
	}
}

func TestFileEntryTempLocation(t *testing.T) {
	entry := &FileEntry{"test1", 55, time.Now(), 0}
	if entry.TempLocation("/some/path") != "/some/path/test1.part" {
		t.Errorf("file entry temp location should be correct")
	}
}
