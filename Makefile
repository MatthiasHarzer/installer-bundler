BUILD_VERSION ?= "unknown"

generate-runtime:
	# It is required to rename the go.mod file in order for it to be embedded. To keep development of
	# the runtime module easier, we generate a copy of the runtime with a renamed go.mod file.
	rm -rf generated/installer-runtime
	mkdir -p generated
	cp -r installer-runtime generated/
	mv generated/installer-runtime/go.mod generated/installer-runtime/go.mod.embed

build: generate-runtime
	@GOOS=windows GOARCH=amd64 go build -o build/bundler-windows-amd64.exe -ldflags "-X installer-bundler.Version=$(BUILD_VERSION)" commands/main.go
	@GOOS=linux GOARCH=amd64 go build -o build/bundler-linux-amd64 -ldflags "-X installer-bundler.Version=$(BUILD_VERSION)" commands/main.go

.PHONY:	build \
		generate-runtime \