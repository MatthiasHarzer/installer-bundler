package root

import (
	"embed"
	"io/fs"

	"installer-bundler/util/fsutil"
)

//go:embed generated/installer-runtime
var installerRuntime embed.FS
var InstallerRuntimeFS fs.FS

var Version = "unknown"

func init() {
	runtimeFS, err := fs.Sub(installerRuntime, "generated/installer-runtime")
	if err != nil {
		panic(err)
	}

	InstallerRuntimeFS = fsutil.GoModuleEmbedFS(runtimeFS, "go.mod.embed")
}
