package root

import (
	"embed"
	"os"
)

//go:embed all:files
var Files embed.FS

var Version = "unknown"

var AppDataDir string

func init() {
	appDataBaseDir := os.Getenv("APPDATA")
	if appDataBaseDir == "" {
		appDataBaseDir = os.Getenv("HOME")
	}

	AppDataDir = appDataBaseDir + string(os.PathSeparator) + ".installer-bundler"
	err := os.MkdirAll(AppDataDir, 0777)
	if err != nil {
		panic(err)
	}
}
