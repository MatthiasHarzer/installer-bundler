BUILD_VERSION ?= "unknown"

build:
	@GOOS=windows GOARCH=amd64 go build -o ./build/installer-runtime.exe -ldflags "-X main.version=$(BUILD_VERSION)" ./main.go

.PHONY: build