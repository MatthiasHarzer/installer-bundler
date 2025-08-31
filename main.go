package main

import (
	"embed"
	"fmt"
	"io/fs"

	"installer-bundler/commands/build"
	"installer-bundler/util/fsutil"

	"github.com/spf13/cobra"
)

//go:embed build/installer-runtime
var installerRuntime embed.FS

var version = "unknown"

var command = &cobra.Command{
	Use:   "installer-bundler",
	Short: "Bundles multiple installer files into a single executable",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("installer-bundler version", version)
	},
}

func init() {
	runtimeFS, err := fs.Sub(installerRuntime, "build/installer-runtime")
	if err != nil {
		panic(err)
	}

	embedFS := fsutil.GoModuleEmbedFS(runtimeFS, "go.mod.embed")
	//fs.WalkDir(embedFS, ".", func(path string, d fs.DirEntry, err error) error {
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Println("Embedded file:", path)
	//	return nil
	//})

	build.InstallerRuntime = embedFS
	command.AddCommand(build.Command)
}

func main() {
	err := command.Execute()
	if err != nil {
		panic(err)
	}

	// The following code is commented out to prevent automatic execution during tests.
	// Uncomment to enable building the installer.

	//println("Building installer with the following items:")
	//for title, link := range items {
	//	println("-", title+":", link)
	//}
	//var coreItems []core.Item
	//for title, link := range items {
	//	coreItems = append(coreItems, core.Item{
	//		Title: title,
	//		Link:  link,
	//	})
	//}
	//bundler := core.NewBundler(coreItems, "installer-runtime")
	//projectDir, err := bundler.GenerateProject()
	//if err != nil {
	//	panic(err)
	//}
	//
	//installerPath, err := core.BuildProject(projectDir)
	//if err != nil {
	//	panic(err)
	//}
	//
	//println("Installer built at:", installerPath)
}
