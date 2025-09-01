package core

import (
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/exec"
	"path"
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

func (r *Runtime) getFileName(item config.Item) (string, error) {
	if item.File != nil {
		return path.Base(*item.File), nil
	}

	if item.URL != nil {
		response, err := http.Head(*item.URL)
		if err != nil || response.StatusCode != http.StatusOK {
			return "", fmt.Errorf("failed to get file name from URL: %s", *item.URL)
		}

		return getFileNameFromHeader(response.Header, *item.URL)
	}

	return "", fmt.Errorf("item has neither URL nor File specified")
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
			if strings.EqualFold(item.Name, itemName) {
				filtered = append(filtered, item)
				break
			}
		}
	}
	return filtered
}

func (r *Runtime) IsExtracted(item config.Item) (bool, string) {
	filename, err := r.getFileName(item)
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

	filename, err := getFileNameFromHeader(response.Header, *item.URL)
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
	if err != nil {
		return destPath, nil
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return "", err
	}

	return destPath, nil
}

func (r *Runtime) Install(item config.Item, shouldDetach bool) (*exec.Cmd, error) {
	isExtracted, filePath := r.IsExtracted(item)
	if !isExtracted {
		return nil, fmt.Errorf("item not found: %s", item.Name)
	}

	cmd := exec.Command(filePath)

	if shouldDetach {
		err := cmd.Start()
		if err != nil {
			return nil, err
		}
		return cmd, nil
	}

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return cmd, nil
}
