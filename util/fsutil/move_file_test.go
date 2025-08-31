package fsutil_test

import (
	"os"
	"testing"

	"installer-bundler/util/fsutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMoveFile(t *testing.T) {
	t.Run("moves file successfully", func(t *testing.T) {
		srcFile, cleanup, err := fsutil.CreateTemporaryFile("movefile_src_*.txt")
		require.NoError(t, err)
		defer func() { _ = cleanup() }()

		_, err = srcFile.WriteString("test content")
		require.NoError(t, err)
		srcFile.Close()

		destPath := srcFile.Name() + "_dest"
		err = fsutil.MoveFile(srcFile.Name(), destPath)
		require.NoError(t, err)

		_, err = os.Stat(destPath)
		require.NoError(t, err)
		assert.False(t, fsutil.FileExists(srcFile.Name()), "source file should not exist after move")

		file, err := os.ReadFile(destPath)
		require.NoError(t, err)
		assert.Equal(t, "test content", string(file))

		err = os.Remove(destPath)
		require.NoError(t, err)
	})

	t.Run("returns error when source file does not exist", func(t *testing.T) {
		err := fsutil.MoveFile("nonexistent_src.txt", "dest.txt")
		require.Error(t, err)
	})
}
