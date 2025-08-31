package core

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"

	"installer-bundler/util/fsutil"
)

type Item struct {
	Title string
	Link  string
}

type Bundler struct {
	items            []Item
	runtimeProjectFS fs.FS
}

func NewBundler(items []Item, runtimeProject fs.FS) *Bundler {
	return &Bundler{
		items:            items,
		runtimeProjectFS: runtimeProject,
	}
}

func (b *Bundler) GetItems() []Item {
	return b.items
}

func (b *Bundler) IsDownloaded(item Item) (bool, string) {
	response, err := http.Head(item.Link)
	if err != nil || response.StatusCode != http.StatusOK {
		return false, ""
	}

	filename, err := getFileName(response.Header, item.Link)
	if err != nil {
		return false, ""
	}

	path := filePath(filename)
	exists := fsutil.FileExists(path)
	if !exists {
		return false, ""
	}

	return true, path
}

func (b *Bundler) Download(item Item) (string, error) {
	isDownloaded, fp := b.IsDownloaded(item)
	if isDownloaded {
		return fp, nil
	}

	response, err := http.Get(item.Link)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s", response.Status)
	}

	filename, err := getFileName(response.Header, item.Link)
	if err != nil {
		return "", err
	}

	temporaryDownloadFile, cleanup, err := fsutil.CreateTemporaryFile(filename)
	if err != nil {
		return "", err
	}
	defer cleanup()
	defer temporaryDownloadFile.Close()

	_, err = io.Copy(temporaryDownloadFile, response.Body)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filesBaseDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	temporaryDownloadFile.Close()
	fp = filePath(filename)
	err = fsutil.MoveFile(temporaryDownloadFile.Name(), fp)
	if err != nil {
		return "", err
	}

	return fp, nil
}
