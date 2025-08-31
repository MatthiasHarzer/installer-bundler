package root

import (
	"embed"
)

//go:embed files/*
var Files embed.FS

var Version = "unknown"
