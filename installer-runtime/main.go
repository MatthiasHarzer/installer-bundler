package main

import (
	"fmt"

	"installer-runtime/commands/download"
	"installer-runtime/commands/list"

	"github.com/spf13/cobra"
)

func init() {
	command.AddCommand(download.Command)
	command.AddCommand(list.Command)
}

var version = "unknown"

var command = &cobra.Command{
	Use:   "installer-runtime",
	Short: "Downloads or runs installer files",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installer-runtime", version)
	},
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}
}
