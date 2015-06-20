package main

import "net/http"
import "time"
import "encoding/json"
import "os"

func filesHandler(fileList *FileList) func(w http.ResponseWriter, r *http.Request) {
  return func (w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(fileList.Current)
  }
}

func newFilesHandler(fileList *FileList) func(w http.ResponseWriter, r *http.Request) {
  return func (w http.ResponseWriter, req *http.Request) {
    json.NewEncoder(w).Encode(fileList.New)
  }
}

func keepUpToDate(fileList *FileList, rootPath string) {
  for {
    println("Updating...")
    scannedFiles := make(chan *FileEntry, 100)
    hashedFiles :=  make(chan *FileEntry, 100)
    go scan(rootPath, scannedFiles)
    go hash(rootPath, scannedFiles, hashedFiles)
    fileList.update(hashedFiles)

    time.Sleep(time.Duration(240)*time.Second)
  }
}

func main() {

  rootPath := os.Args[1]

  fileList := &FileList{[]FileEntry{},[]FileEntry{}}

  go keepUpToDate(fileList, rootPath)

  http.HandleFunc("/file-list", filesHandler(fileList))
  http.HandleFunc("/new-file-list", newFilesHandler(fileList))
  http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(rootPath))))
  //
  err := http.ListenAndServe(":1978", nil)
  check(err)
}
