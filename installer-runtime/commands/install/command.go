package install

import (
	"fmt"
	"os/exec"
	"path"
	"sync"

	root "installer-runtime"
	"installer-runtime/commands/extract"
	"installer-runtime/config"
	"installer-runtime/core"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var directory string
var itemNames []string
var configFile string
var parallel bool

func init() {
	Command.Flags().BoolVarP(&parallel, "parallel", "p", false, "Install items in parallel instead of sequentially")
	Command.Flags().StringVarP(&directory, "directory", "d", "", "Directory to save files to")
	Command.Flags().StringArrayVarP(&itemNames, "item", "i", []string{}, "Names of items to install (if empty, all items are installed)")
	Command.Flags().StringVarP(&configFile, "config", "c", "", "Path to config file")
}

func getDownloadDir() string {
	downloadDir := directory
	if downloadDir == "" {
		downloadDir = path.Join(root.AppDataDir, "downloads")
	}
	return downloadDir
}

func getConfig() (*config.Config, error) {
	if configFile != "" {
		cfg, err := config.GetConfigFromFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load config from file: %w", err)
		}
		return cfg, nil
	}

	return config.GetConfig()
}

var Command = &cobra.Command{
	Use:   "install",
	Short: "Installs all available binaries",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := getConfig()
		if err != nil {
			return err
		}

		downloadDir := getDownloadDir()

		runtime := core.NewRuntime(*cfg, downloadDir, root.Files)
		filteredItems := runtime.GetItems(itemNames)

		switch cfg.Mode {
		case config.ModeURL:
			extract.DownloadFiles(runtime, filteredItems)
		case config.ModeEmbedded:
			extract.CopyFiles(runtime, filteredItems)
		default:
			return fmt.Errorf("unknown mode: %s", cfg.Mode)
		}

		cmds := make(map[config.Item]*exec.Cmd)

		fmt.Println()
		color.Green("Installing %d item(s)", len(filteredItems))
		for _, item := range filteredItems {
			fmt.Printf(`Installing "%s"`, item.Name)
			cmd, err := runtime.Install(*item, parallel)
			if err != nil {
				color.Red(" - failed: %v", err)
				continue
			}
			if parallel {
				color.Green(" - started")
			} else {
				color.Green(" - completed")
			}
			cmds[*item] = cmd
		}

		if parallel {
			fmt.Println()
			wg := sync.WaitGroup{}
			for item, cmd := range cmds {
				wg.Add(1)
				go func(cmd *exec.Cmd) {
					defer wg.Done()
					err := cmd.Wait()
					if err != nil {
						color.Red(`Installation of "%s" unsuccessful. Failed to execute command "%s": %v`, item.Name, cmd.String(), err)
					} else {
						color.Green(`Installation of "%s" successful`, item.Name)
					}
				}(cmd)
			}
			wg.Wait()
		}

		return nil
	},
}
