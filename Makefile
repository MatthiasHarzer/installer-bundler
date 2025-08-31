
generate-runtime:
	# It is required to rename the go.mod file in order for it to be embedded. To keep development of
	# the runtime module easier, we generate a copy of the runtime with a renamed go.mod file.
	cp -r installer-runtime build
	mv build/installer-runtime/go.mod build/installer-runtime/go.mod.embed

build: generate-runtime
	@GOOS=windows GOARCH=amd64 go build -o build/installer-windows-amd64.exe main.go

.PHONY:	build \
		generate-runtime