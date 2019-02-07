package agent

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/nicktitle/rc/pkg/fileformat"
)

const testFileName = "definitely-not-taken.txt"

func TestFileCreate(t *testing.T) {
	setupFileTests(t)
	// defer call will wipe the file out again at the end
	defer setupFileTests(t)

	event, err := createFile(sharedFileSettings(t))
	require.NoError(t, err)
	require.Equal(t, event["activity"], fileformat.FileCreate)

	_, err = os.Stat(testFileName)
	require.NoError(t, err)
}

func TestFileModify(t *testing.T) {
	setupFileTests(t)
	// defer call will wipe the file out again at the end
	defer setupFileTests(t)

	_, err := os.Create(testFileName)
	info, err := os.Stat(testFileName)
	initSize := info.Size()

	event, err := modifyFile(sharedFileSettings(t))
	require.NoError(t, err)
	require.Equal(t, event["activity"], fileformat.FileModify)

	info, err = os.Stat(testFileName)
	// we modify by adding bytes to the file, so check that the size is different
	require.NotEqual(t, initSize, info.Size())
}

func TestFileDelete(t *testing.T) {
	setupFileTests(t)
	// defer call will wipe the file out again at the end
	defer setupFileTests(t)

	_, err := os.Create(testFileName)

	event, err := deleteFile(sharedFileSettings(t))
	require.NoError(t, err)
	require.Equal(t, event["activity"], fileformat.FileDelete)

	_, err = os.Stat(testFileName)
	require.Error(t, err)
}

// designating these as helpers will ensure that they are ignored in reporting
func setupFileTests(t *testing.T) {
	t.Helper()
	os.Remove(testFileName)
}

func sharedFileSettings(t *testing.T) fileformat.FileSettings {
	t.Helper()
	return fileformat.FileSettings{Path: testFileName}
}
