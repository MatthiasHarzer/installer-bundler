package root

import (
	"embed"
	"io/fs"
)

//go:embed files/*
var files embed.FS
var Files fs.FS

var Version = "unknown"

func init() {
	var err error
	Files, err = fs.Sub(files, "files")
	if err != nil {
		panic(err)
	}
}
