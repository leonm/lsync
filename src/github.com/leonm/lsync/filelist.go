package main

type FileList struct {
  Current []FileEntry
  New []FileEntry
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
