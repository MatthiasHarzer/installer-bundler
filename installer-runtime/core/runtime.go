package core

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"installer-runtime/config"
	"installer-runtime/util/fsutil"
)

type Runtime struct {
	cfg         config.Config
	downloadDir string
}

func NewRuntime(cfg config.Config, downloadDir string) *Runtime {
	return &Runtime{
		cfg:         cfg,
		downloadDir: downloadDir,
	}
}

func (r *Runtime) filePath(fileName string) string {
	return fmt.Sprintf("%s/%s", r.downloadDir, fileName)
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
	response, err := http.Head(item.URL)
	if err != nil || response.StatusCode != http.StatusOK {
		return false, ""
	}

	filename, err := getFileName(response.Header, item.URL)
	if err != nil {
		return false, ""
	}

	filePath := r.filePath(filename)
	exists := fsutil.FileExists(filePath)
	if !exists {
		return false, ""
	}

	return true, filePath
}

func (r *Runtime) DownloadItem(item config.Item) (string, error) {
	response, err := http.Get(item.URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to download file: %s", response.Status)
	}

	filename, err := getFileName(response.Header, item.URL)
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

	if !fsutil.FileExists(r.downloadDir) {
		err = os.MkdirAll(r.downloadDir, os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	temporaryDownloadFile.Close()
	filePath := r.filePath(filename)
	err = fsutil.MoveFile(temporaryDownloadFile.Name(), filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}
