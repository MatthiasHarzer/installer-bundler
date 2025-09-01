package main

import (
	"fmt"

	root "installer-runtime"
	"installer-runtime/commands/extract"
	"installer-runtime/commands/install"
	"installer-runtime/commands/list"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(extract.Command)
	command.AddCommand(install.Command)
	command.AddCommand(list.Command)
}

var command = &cobra.Command{
	Use:   "installer-runtime",
	Short: "Downloads or runs installer files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installer-runtime version", root.Version)
	},
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
