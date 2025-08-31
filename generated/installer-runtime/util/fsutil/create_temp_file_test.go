package fsutil_test

import (
	"os"
	"strings"
	"testing"

	"installer-runtime/util/fsutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTemporaryFile(t *testing.T) {
	t.Run("creates a file with the given pattern", func(t *testing.T) {
		file, cleanup, err := fsutil.CreateTemporaryFile("testfile_*.txt")
		require.NoError(t, err)

		defer file.Close()
		defer func() {
			err := cleanup()
			assert.NoError(t, err)
		}()

		require.NotNil(t, file)

		stat, err := os.Stat(file.Name())
		require.NoError(t, err)
		assert.False(t, stat.IsDir(), "expected a file, but got a directory")

		expectedPattern := "testfile_"
		assert.True(t, strings.HasPrefix(stat.Name(), expectedPattern))
	})

	t.Run("cleanup removes the file", func(t *testing.T) {
		file, cleanup, err := fsutil.CreateTemporaryFile("tempfile_*.tmp")
		require.NoError(t, err)

		filePath := file.Name()
		_, err = os.Stat(filePath)
		require.NoError(t, err, "expected the temporary file to exist")

		err = cleanup()
		require.NoError(t, err)

		_, err = os.Stat(filePath)
		assert.True(t, os.IsNotExist(err), "expected the temporary file to be removed")
	})
}
