

generate-runtime:
	# It is required to rename the go.mod file in order for it to be embedded. To keep development of
	# the runtime module easier, we generate a copy of the runtime with a renamed go.mod file.
	rm -rf generated/installer-runtime
	mkdir -p build
	cp -r installer-runtime generated
	mv generated/installer-runtime/go.mod generated/installer-runtime/go.mod.embed

build: generate-runtime
	@GOOS=windows GOARCH=amd64 go build -o build/bundler-windows-amd64.exe commands/main.go

.PHONY:	build \
		generate-runtime \