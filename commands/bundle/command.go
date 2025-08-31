package bundle

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	root "installer-bundler"
	"installer-bundler/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var outputFile string
var embedded bool
var file string

func init() {
	Command.Flags().StringVarP(&outputFile, "output", "o", "", "Output file")
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

func getOutputFilePath() string {
	outputFilePath := outputFile
	if outputFile == "" {
		outputFilePath = "installer-runtime"
	}
	if !strings.HasSuffix(outputFilePath, ".exe") && runtime.GOOS == "windows" {
		outputFilePath += ".exe"
	}
	return outputFilePath
}

var Command = &cobra.Command{
	Use:   "bundle",
	Short: "Bundles the references executables",
	RunE: func(cmd *cobra.Command, args []string) error {
		outputFilePath := getOutputFilePath()

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
			color.Yellow("Embedding binaries into the executable")

			// Download files if not already downloaded
			for _, item := range items {
				isDownloaded, filePath := bundler.IsDownloaded(item)
				if !isDownloaded {
					fmt.Printf(`Downloading "%s"`, item.Title)
					filePath, err := bundler.Download(item)
					if err != nil {
						return err
					}
					color.Green(`- Saved to "%s"`, filePath)
				} else {
					color.Yellow(`Already downloaded "%s" to "%s"`, item.Title, filePath)
				}
			}

			err = bundler.BuildProjectEmbedded(outputFilePath)
		} else {
			err = bundler.BuildProjectURL(outputFilePath)
		}
		if err != nil {
			return err
		}

		fmt.Println("Build successful:", outputFilePath)

		return nil
	},
}
