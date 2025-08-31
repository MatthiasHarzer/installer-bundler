package download

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	root "installer-runtime"
	"installer-runtime/config"
	"installer-runtime/core"
	"installer-runtime/util/fsutil"
	"installer-runtime/util/windowsutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var directory string
var itemNames []string

func init() {
	defaultDownloadDir, err := windowsutil.GetKnownFolderPath(windowsutil.DownloadsFolder)
	if err != nil {
		defaultDownloadDir = "."
	}

	Command.Flags().StringVarP(&directory, "directory", "d", defaultDownloadDir, "Directory to save files to")
	Command.Flags().StringArrayVarP(&itemNames, "item", "i", []string{}, "Names of items to download (if empty, all items will be downloaded)")
}

func getDownloadDir() string {
	p := path.Clean(directory)

	if path.IsAbs(p) {
		return p
	}

	cwd, err := os.Getwd()
	if err != nil {
		return p
	}

	path.Join(cwd, p)
	return p
}

func download(runtime *core.Runtime, items []*config.Item) {
	color.Green("Downloading %d file(s)\n", len(items))

	for _, item := range items {
		isDownloaded, file := runtime.IsDownloaded(*item)
		if isDownloaded {
			color.Yellow("Skipping \"%s\" - already downloaded to %s\n", item.Name, file)
			continue
		}

		fmt.Printf(color.WhiteString(`Downloading "%s"`), item.Name)

		file, err := runtime.DownloadItem(*item)
		if err != nil {
			fmt.Printf(" - ")
			color.Red("Failed: %s\n", err.Error())
			continue
		}

		downloadedFile := path.Clean(file)
		color.Green(" - Saved to %s\n", downloadedFile)
	}
}

func copyFiles(outputDirectory string, items []*config.Item) {
	color.Green("Copying %d file(s)\n", len(items))

	err := fs.WalkDir(root.Files, ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fpath, err := filepath.Localize(p)
		if err != nil {
			return err
		}
		newPath := path.Join(outputDirectory, fpath)

		if fsutil.FileExists(newPath) {
			return nil
		}

		switch d.Type() {
		case os.ModeDir:
			return os.MkdirAll(newPath, 0777)
		case 0:
			fmt.Printf(color.WhiteString(`Copying "%s"`), p)

			r, err := root.Files.Open(p)
			if err != nil {
				return err
			}
			defer r.Close()
			info, err := r.Stat()
			if err != nil {
				return err
			}
			w, err := os.OpenFile(newPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666|info.Mode()&0777)
			if err != nil {
				return err
			}

			if _, err := io.Copy(w, r); err != nil {
				w.Close()
				return &os.PathError{Op: "Copy", Path: newPath, Err: err}
			}

			color.Green(" - Saved to %s\n", newPath)

			return w.Close()
		default:
			return &os.PathError{Op: "CopyFS", Path: p, Err: os.ErrInvalid}
		}
	})

	if err != nil {
		color.Red("Failed to copy files: %s\n", err.Error())
		return
	}
}

var Command = &cobra.Command{
	Use:   "download",
	Short: "Download files from URLs",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		downloadDir := getDownloadDir()

		runtime := core.NewRuntime(*cfg, downloadDir)
		filteredItems := runtime.GetItems(itemNames)

		switch cfg.Mode {
		case config.ModeURL:
			download(runtime, filteredItems)
		case config.ModeEmbedded:
			copyFiles(downloadDir, filteredItems)
		default:
			return fmt.Errorf("unknown mode: %s", cfg.Mode)
		}

		return nil
	},
}
