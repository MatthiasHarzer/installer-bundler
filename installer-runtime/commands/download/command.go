package download

import (
	"os"
	"path"

	"installer-runtime/config"
	"installer-runtime/core"
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

		color.Green("Downloading %d file(s)\n", len(filteredItems))

		for _, item := range filteredItems {
			isDownloaded, file := runtime.IsDownloaded(*item)
			if isDownloaded {
				color.Yellow("Skipping \"%s\" - already downloaded to %s\n", item.Name, file)
				continue
			}

			cmd.Printf(color.WhiteString(`Downloading "%s"`), item.Name)

			file, err := runtime.DownloadItem(*item)
			if err != nil {
				cmd.Printf(" - ")
				color.Red("Failed: %s\n", err.Error())
				continue
			}

			downloadedFile := path.Clean(file)
			color.Green(" - Saved to %s\n", downloadedFile)
		}

		return nil
	},
}
