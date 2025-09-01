package extract

import (
	"fmt"
	"os"
	"path"
	"runtime"

	root "installer-runtime"
	"installer-runtime/config"
	"installer-runtime/core"
	"installer-runtime/util/windowsutil"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var directory string
var itemNames []string

func init() {
	Command.Flags().StringVarP(&directory, "directory", "d", "", "Directory to save files to")
	Command.Flags().StringArrayVarP(&itemNames, "item", "i", []string{}, "Names of items to download (if empty, all items will be downloaded)")
}

func getDownloadDir() string {
	var err error
	p := directory
	if p == "" {
		switch runtime.GOOS {
		case "windows":
			p, err = windowsutil.GetKnownFolderPath(windowsutil.DownloadsFolder)
			if err != nil {
				p = path.Join(os.Getenv("USERPROFILE"), "Downloads")
			}
		default:
			p = path.Join(os.Getenv("HOME"), "Downloads")
		}
	}

	p = path.Clean(p)

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

func DownloadFiles(runtime *core.Runtime, items []*config.Item) {
	color.Green(`Downloading %d file(s) to "%s"`, len(items), runtime.OutputDirectory)

	for _, item := range items {
		isDownloaded, filePath := runtime.IsExtracted(*item)
		if isDownloaded {
			color.Yellow(`Skipping "%s" - already downloaded to "%s"`, item.Name, filePath)
			continue
		}

		fmt.Printf(color.WhiteString(`Downloading "%s"`), item.Name)

		filePath, err := runtime.DownloadItem(*item)
		if err != nil {
			fmt.Printf(" - ")
			color.Red("ailed: %s\n", err.Error())
			continue
		}

		downloadedFile := path.Clean(filePath)
		color.Green(` - saved to "%s"`, downloadedFile)
	}
}

func CopyFiles(runtime *core.Runtime, items []*config.Item) {
	color.Green(`Copying %d file(s) to "%s"`, len(items), runtime.OutputDirectory)

	for _, item := range items {
		isCopied, filePath := runtime.IsExtracted(*item)
		if isCopied {
			color.Yellow(`Skipping "%s" - already copied to "%s"`, item.Name, filePath)
			continue
		}

		fmt.Printf(color.WhiteString(`Copying "%s"`), item.Name)

		filePath, err := runtime.CopyItem(*item)
		if err != nil {
			fmt.Printf(" - ")
			color.Red("failed: %s\n", err.Error())
			continue
		}

		copiedFile := path.Clean(filePath)
		color.Green(` - saved to "%s"`, copiedFile)
	}
}

var Command = &cobra.Command{
	Use:   "extract",
	Short: "Downloads or copies the embedded files or URLs to the specified directory",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		downloadDir := getDownloadDir()

		runtime := core.NewRuntime(*cfg, downloadDir, root.Files)
		filteredItems := runtime.GetItems(itemNames)

		switch cfg.Mode {
		case config.ModeURL:
			DownloadFiles(runtime, filteredItems)
		case config.ModeEmbedded:
			CopyFiles(runtime, filteredItems)
		default:
			return fmt.Errorf("unknown mode: %s", cfg.Mode)
		}

		return nil
	},
}
