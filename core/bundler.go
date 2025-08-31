package core

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path"

	"installer-bundler/util/fsutil"
)

type Item struct {
	Title string
	Link  string
}

type Bundler struct {
	items            []Item
	runtimeProjectFS fs.FS
	fileCacheDir     string
}

func NewBundler(items []Item, runtimeProject fs.FS, fileCacheDir string) *Bundler {
	return &Bundler{
		items:            items,
		runtimeProjectFS: runtimeProject,
		fileCacheDir:     fileCacheDir,
	}
}

func (b *Bundler) filePath(fileName string) string {
	return path.Join(b.fileCacheDir, fileName)
}

func (b *Bundler) runtimeFilesPath(filename string) string {
	return fmt.Sprintf("%s/%s", runtimeFilesDir, filename)
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

	path := b.filePath(filename)
	exists := fsutil.FileExists(path)
	if !exists {
		return false, ""
	}

	return true, b.runtimeFilesPath(filename)
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

	err = os.MkdirAll(b.fileCacheDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	temporaryDownloadFile.Close()
	fp = b.filePath(filename)
	err = fsutil.MoveFile(temporaryDownloadFile.Name(), fp)
	if err != nil {
		return "", err
	}

	return b.runtimeFilesPath(filename), nil
}
