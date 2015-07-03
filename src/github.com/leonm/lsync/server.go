package main

import "github.com/codegangsta/cli"
import "net/http"
import "time"
import "encoding/json"

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
    rootDirectory := NewRootDirectory (rootPath)
    go rootDirectory.scan(scannedFiles)
    go rootDirectory.hash(scannedFiles, hashedFiles)
    fileList.update(hashedFiles)

    time.Sleep(time.Duration(240)*time.Second)
  }
}

func newServerCommand() func (c *cli.Context) {
  return func (c *cli.Context) {
    rootPath := c.Args()[0]

    fileList := &FileList{[]FileEntry{},[]FileEntry{}}

    go keepUpToDate(fileList, rootPath)

    http.HandleFunc("/file-list", filesHandler(fileList))
    http.HandleFunc("/new-file-list", newFilesHandler(fileList))
    http.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir(rootPath))))

    err := http.ListenAndServe(":1978", nil)
    check(err)
  }
}
