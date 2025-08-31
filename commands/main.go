package main

import (
	"fmt"

	root "installer-bundler"
	"installer-bundler/commands/build"

	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "installer-bundler",
	Short: "Bundles multiple installer files into a single executable",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installer-bundler version", root.Version)
	},
}

func init() {
	command.AddCommand(build.Command)
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
