package root

import (
	"embed"
	"io/fs"
	"os"

	"installer-bundler/util/fsutil"
)

//go:embed generated/installer-runtime
var installerRuntime embed.FS
var InstallerRuntimeFS fs.FS

var Version = "unknown"

var AppDataDir string

func init() {
	runtimeFS, err := fs.Sub(installerRuntime, "generated/installer-runtime")
	if err != nil {
		panic(err)
	}

	InstallerRuntimeFS = fsutil.GoModuleEmbedFS(runtimeFS, "go.mod.embed")

	appDataBaseDir := os.Getenv("APPDATA")
	if appDataBaseDir == "" {
		appDataBaseDir = os.Getenv("HOME")
	}

	AppDataDir = appDataBaseDir + string(os.PathSeparator) + ".installer-bundler"
	err = os.MkdirAll(AppDataDir, 0777)
	if err != nil {
		panic(err)
	}
}
