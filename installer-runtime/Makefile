BUILD_VERSION ?= "unknown"

build:
	@GOOS=windows GOARCH=amd64 go build -o ./build/installer-runtime.exe -ldflags "-X root.Version=$(BUILD_VERSION)" ./commands/main.go
	@GOOS=linux GOARCH=amd64 go build -o ./build/installer-runtime -ldflags "-X root.Version=$(BUILD_VERSION)" ./commands/main.go

.PHONY: build