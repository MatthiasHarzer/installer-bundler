package root

import (
	"embed"
)

//go:embed files
var Files embed.FS

//var Files fs.FS

var Version = "unknown"

//func init() {
//	var err error
//	Files, err = fs.Sub(files, "files")
//	if err != nil {
//		panic(err)
//	}
//}
