package build

import (
	"fmt"
	"os"
	"path"
	"strings"

	root "installer-bundler"
	"installer-bundler/core"

	"github.com/spf13/cobra"
)

var outputFile string
var embedded bool
var file string

func init() {
	Command.Flags().StringVarP(&outputFile, "output", "o", "output.exe", "Output file")
	Command.Flags().BoolVarP(&embedded, "embedded", "e", false, "Embedded binaries")
	Command.Flags().StringVarP(&file, "file", "f", "", "File containing list of items (title and link separated by comma, one item per line)")
}

func loadItemsFromFile(filePath string) ([]core.Item, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(content), "\n")
	var result []core.Item
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ",", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format: %s", line)
		}
		title := strings.TrimSpace(parts[0])
		link := strings.TrimSpace(parts[1])
		result = append(result, core.Item{
			Title: title,
			Link:  link,
		})
	}

	return result, nil
}

var Command = &cobra.Command{
	Use:   "build",
	Short: "Builds the project",
	RunE: func(cmd *cobra.Command, args []string) error {
		if outputFile == "" {
			outputFile = "installer-runtime.exe"
		}

		if !strings.HasSuffix(outputFile, ".exe") {
			outputFile += ".exe"
		}

		if file == "" {
			return fmt.Errorf("please provide a file containing the list of items using the --file flag")
		}

		items, err := loadItemsFromFile(file)
		if err != nil {
			return fmt.Errorf("failed to load items from file: %w", err)
		}

		fmt.Println("Loaded", len(items), "items from file")

		filesDir := path.Join(root.AppDataDir, "files")

		bundler := core.NewBundler(items, root.InstallerRuntimeFS, filesDir)

		if embedded {
			for _, item := range items {
				isDownloaded, _ := bundler.IsDownloaded(item)
				if !isDownloaded {
					fmt.Println("Downloading:", item.Title)
					_, err := bundler.Download(item)
					if err != nil {
						return err
					}
				} else {
					fmt.Println("Already downloaded:", item.Title)
				}
			}

			err = bundler.BuildProjectEmbedded(outputFile)
		} else {
			err = bundler.BuildProjectURL(outputFile)
		}
		if err != nil {
			return err
		}

		fmt.Println("Build successful:", outputFile)

		return nil
	},
}
