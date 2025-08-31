package list

import (
	"installer-runtime/config"

	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "list",
	Short: "Lists available installer items",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return err
		}

		println("Mode:", cfg.Mode)
		println("Available installer items:")
		for _, item := range cfg.Items {
			if cfg.Mode == config.ModeURL {
				println("-", item.Name+":", *item.URL)
			} else {
				println("-", item.Name+":", *item.File)
			}
		}

		return nil
	},
}
