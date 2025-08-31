package build

import (
	"fmt"
	"io/fs"
	"strings"

	"installer-bundler/core"

	"github.com/spf13/cobra"
)

var items = map[string]string{
	"Chrome":  "https://dl.google.com/chrome/install/375.126/chrome_installer.exe",
	"Firefox": "https://download-installer.cdn.mozilla.net/pub/firefox/releases/113.0/win64/en-US/Firefox%20Setup%20113.0.exe",
	"VLC":     "https://get.videolan.org/vlc/3.0.18/win64/vlc-3.0.18-win64.exe",
	"7-Zip":   "https://www.7-zip.org/a/7z1900-x64.exe",
}

var InstallerRuntime fs.FS
var outputFile string

func init() {
	Command.Flags().StringVarP(&outputFile, "output", "o", "output.exe", "Output file")
}

var Command = &cobra.Command{
	Use:   "build",
	Short: "Builds the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !strings.HasSuffix(outputFile, ".exe") {
			outputFile += ".exe"
		}

		var coreItems []core.Item
		for title, link := range items {
			coreItems = append(coreItems, core.Item{
				Title: title,
				Link:  link,
			})
		}

		bundler := core.NewBundler(coreItems, InstallerRuntime)
		err := bundler.BuildProject(outputFile)
		if err != nil {
			return err
		}

		fmt.Println("Build successful:", outputFile)

		return nil
	},
}
