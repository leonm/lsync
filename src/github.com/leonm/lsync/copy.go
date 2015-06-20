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

func getFiles(host string, listname string, out chan FileEntry) {
  fileList := getList("http://"+host+":1978/"+listname)
  for _,f := range fileList {
    // println(f.Path)
    out <- f
  }
  close(out)
}

func newCopyCommand() func (c *cli.Context) {
  return func (c *cli.Context) {
    host := c.Args()[0]
    targetPath := c.Args()[1]
    remoteFiles := make(chan FileEntry, 100)
    go getFiles(host,"new-file-list",remoteFiles)

    for f := range remoteFiles {
      dir := filepath.Dir(f.Path)
      err := os.MkdirAll(filepath.Join(targetPath,dir),0777)
      check(err)
      out, err := os.Create(filepath.Join(targetPath,f.Path))
      check(err)
      defer out.Close()
      resp, err := http.Get("http://"+host+":1978/files/"+f.Path)
      defer resp.Body.Close()
      io.Copy(out, resp.Body)
    }
  }
}
