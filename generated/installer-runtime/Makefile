BUILD_VERSION ?= "unknown"

build:
	@GOOS=windows GOARCH=amd64 go build -o ./build/installer-runtime.exe -ldflags "-X installer-runtime.Version=$(BUILD_VERSION)" ./commands/main.go
	@GOOS=linux GOARCH=amd64 go build -o ./build/installer-runtime -ldflags "-X installer-runtime.Version=$(BUILD_VERSION)" ./commands/main.go

generate-manifest:
	@rsrc -manifest installer-runtime.exe.manifest -o rsrc.syso

.PHONY: build \
		generate-manifest