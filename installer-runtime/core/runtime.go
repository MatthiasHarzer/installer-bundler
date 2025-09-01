package core

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"installer-runtime/config"
	"installer-runtime/util/fsutil"
)

type Runtime struct {
	cfg             config.Config
	files           fs.FS
	OutputDirectory string
}

func NewRuntime(cfg config.Config, outputDirectory string, filesFS fs.FS) *Runtime {
	return &Runtime{
		cfg:             cfg,
		files:           filesFS,
		OutputDirectory: outputDirectory,
	}
}

func (r *Runtime) FilePath(fileName string) string {
	return fmt.Sprintf("%s/%s", r.OutputDirectory, fileName)
}

func (r *Runtime) GetItems(names []string) []*config.Item {
	if len(names) == 0 {
		return r.cfg.Items
	}

	var filtered []*config.Item
	for _, itemName := range names {
		for _, item := range r.cfg.Items {
			if strings.EqualFold(itemName, itemName) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func (r *Runtime) IsDownloaded(item config.Item) (bool, string) {
	response, err := http.Head(*item.URL)
	if err != nil || response.StatusCode != http.StatusOK {
		return false, ""
	}

	filename, err := getFileName(response.Header, *item.URL)
	if err != nil {
		return false, ""
	}

	filePath := r.FilePath(filename)
	exists := fsutil.FileExists(filePath)
	if !exists {
		return false, ""
	}

	return true, filePath
}

func (r *Runtime) DownloadItem(item config.Item) (string, error) {
	response, err := http.Get(*item.URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s", response.Status)
	}

	filename, err := getFileName(response.Header, *item.URL)
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

	if !fsutil.FileExists(r.OutputDirectory) {
		err = os.MkdirAll(r.OutputDirectory, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	temporaryDownloadFile.Close()
	filePath := r.FilePath(filename)
	err = fsutil.MoveFile(temporaryDownloadFile.Name(), filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

func (r *Runtime) IsCopied(item config.Item) (bool, string) {
	if item.File == nil {
		return false, ""
	}

	sourcePath := *item.File
	file, err := r.files.Open(sourcePath)
	if err != nil {
		return false, ""
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return false, ""
	}

	destPath := r.FilePath(stat.Name())
	exists := fsutil.FileExists(destPath)
	if !exists {
		return false, ""
	}

	return true, destPath
}

func (r *Runtime) CopyItem(item config.Item) (string, error) {
	sourcePath := *item.File

	file, err := r.files.Open(sourcePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	if !fsutil.FileExists(r.OutputDirectory) {
		err = os.MkdirAll(r.OutputDirectory, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	destPath := r.FilePath(stat.Name())
	destFile, err := os.Create(destPath)
	if err == nil {
		destFile.Close()
		return destPath, nil
	}

	return destPath, nil
}
