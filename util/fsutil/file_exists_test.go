package fsutil_test

import (
	"testing"

	"installer-runtime/util/fsutil"

	"github.com/stretchr/testify/require"
)

func TestFileExists(t *testing.T) {
	t.Run("returns true for existing file", func(t *testing.T) {
		file, cleanup, err := fsutil.CreateTemporaryFile("existingfile_*.txt")
		require.NoError(t, err)

		defer func() {
			err := cleanup()
			require.NoError(t, err)
		}()

		exists := fsutil.FileExists(file.Name())
		require.True(t, exists, "expected file to exist")
	})

	t.Run("returns false for non-existing file", func(t *testing.T) {
		nonExistentFilePath := "nonexistentfile_123456.txt"
		exists := fsutil.FileExists(nonExistentFilePath)
		require.False(t, exists, "expected file to not exist")
	})
}
