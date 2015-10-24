package main

import "net/http"
import "github.com/codegangsta/cli"
import "io/ioutil"
import "encoding/json"
import "os"
import "io"
import "net/url"
import "path/filepath"
import "sync"
import "hash/fnv"

func getList(url string) []FileEntry {
	Info.Println("Getting List " + url)
	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	check(err)
	fileList := []FileEntry{}
	err = json.Unmarshal(body, &fileList)
	check(err)
	Info.Printf("Got %d entries", len(fileList))
	return fileList
}

func getFileList(host string, listname string, out chan FileEntry) {
	fileList := getList("http://" + host + ":1978/" + listname)
	for _, f := range fileList {
		out <- f
	}
	close(out)
}

func ensureDirectory(targetPath string, f *FileEntry) {
	dir := filepath.Dir(f.Path)
	err := os.MkdirAll(filepath.Join(targetPath, dir), 0777)
	check(err)
}

func downloadFile(targetPath string, host string, f *FileEntry) {
	Info.Printf("Starting to download %s", f.Path)
	fileUrl, err := url.Parse("http://" + host + ":1978")
	fileUrl.Path += "/files/" + f.Path
	resp, err := http.Get(fileUrl.String())
	check(err)
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		Error.Printf("Failed to download %s : %s", f.Path, resp.Status)
	} else {
		targetFilePath := filepath.Join(targetPath, f.Path)
		out, err := os.Create(targetFilePath + ".part")
		check(err)
		hasher := fnv.New64a()
		defer out.Close()
		io.Copy(io.MultiWriter(out, hasher), resp.Body)
		if hasher.Sum64() != f.Hash {
			err := os.Remove(targetFilePath + ".part")
			check(err)
			Error.Printf("Failed to download %s.  Hash Error", f.Path)
		} else {
			os.Remove(targetFilePath)
			os.Rename(targetFilePath+".part", targetFilePath)
			os.Chtimes(targetFilePath, f.Updated, f.Updated)
			Info.Printf("Finished downloading %s", f.Path)
		}
	}
}

func needsDownloading(targetPath string, f *FileEntry) bool {
	targetFilePath := filepath.Join(targetPath, f.Path)
	fileInfo, err := os.Stat(targetFilePath)
	if os.IsNotExist(err) {
		return true
	}
	check(err)
	return !f.IsUptoDate(fileInfo.ModTime(), fileInfo.Size())
}

func downloadFiles(host string, targetPath string, remoteFiles chan FileEntry) {
	for f := range remoteFiles {
		ensureDirectory(targetPath, &f)
		if needsDownloading(targetPath, &f) {
			downloadFile(targetPath, host, &f)
		}
	}
}

func startDownloadWorker(host string, targetPath string, remoteFiles chan FileEntry, wg *sync.WaitGroup) {
	downloadFiles(host, targetPath, remoteFiles)
	wg.Done()
}

func getFiles(host string, targetPath string, listname string) {
	remoteFiles := make(chan FileEntry, 100)
	go getFileList(host, listname, remoteFiles)
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		Info.Printf("Starting Worker %d", i)
		wg.Add(1)
		go startDownloadWorker(host, targetPath, remoteFiles, &wg)
	}
	wg.Wait()
}

func newCopyCommand() func(c *cli.Context) {
	return func(c *cli.Context) {

		if c.GlobalIsSet("log-file") {
			InitFileLogging(c.GlobalString("log-file"))
		} else {
			InitStdOutLogging()
		}

		host := c.Args()[0]
		targetPath := c.Args()[1]
		getFiles(host, targetPath, "new-file-list")
		getFiles(host, targetPath, "file-list")
	}
}
