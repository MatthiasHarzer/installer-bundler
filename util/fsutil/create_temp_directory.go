package fsutil

import "os"

func CreateTempDirectory() (string, func() error, error) {
	dir, err := os.MkdirTemp("", "installer-bundler-*")
	if err != nil {
		return "", nil, err
	}
	cleanup := func() error {
		return os.RemoveAll(dir)
	}
	return dir, cleanup, nil
}
