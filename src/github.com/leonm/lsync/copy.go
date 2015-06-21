package main

import "net/http"
import "github.com/codegangsta/cli"
import "io/ioutil"
import "encoding/json"
import "os"
import "io"
import "path/filepath"

func getList(url string) []FileEntry {
  res, err := http.Get(url)
  check(err)
  defer res.Body.Close()
  body, err := ioutil.ReadAll(res.Body)
  check(err)
  fileList := []FileEntry{}
  err = json.Unmarshal(body, &fileList)
  check(err)
  return fileList
}

func getFileList(host string, listname string, out chan FileEntry) {
  fileList := getList("http://"+host+":1978/"+listname)
  for _,f := range fileList {
    out <- f
  }
  close(out)
}

func ensureDirectory(targetPath string, f *FileEntry) {
  dir := filepath.Dir(f.Path)
  err := os.MkdirAll(filepath.Join(targetPath,dir),0777)
  check(err)
}

func getFiles(host string, targetPath string, listname string) {
  remoteFiles := make(chan FileEntry, 100)
  go getFileList(host, listname, remoteFiles)

  for f := range remoteFiles {
    ensureDirectory(targetPath, &f)
    targetFilePath := filepath.Join(targetPath,f.Path)
    out, err := os.Create(targetFilePath)
    check(err)
    defer out.Close()
    resp, err := http.Get("http://"+host+":1978/files/"+f.Path)
    defer resp.Body.Close()
    io.Copy(out, resp.Body)
    os.Chtimes(targetFilePath, f.Updated, f.Updated)
  }
}

func newCopyCommand() func (c *cli.Context) {
  return func (c *cli.Context) {
    host := c.Args()[0]
    targetPath := c.Args()[1]
    getFiles(host, targetPath, "new-file-list")
    getFiles(host, targetPath, "file-list")
  }
}
