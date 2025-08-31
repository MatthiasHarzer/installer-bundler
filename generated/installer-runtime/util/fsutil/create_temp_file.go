package fsutil

import "os"

func CreateTemporaryFile(pattern string) (*os.File, func() error, error) {
	file, err := os.CreateTemp("", pattern)
	if err != nil {
		return nil, nil, err
	}
	cleanup := func() error {
		file.Close()
		err := os.Remove(file.Name())
		if err != nil {
			return err
		}
		return nil
	}
	return file, cleanup, nil
}
