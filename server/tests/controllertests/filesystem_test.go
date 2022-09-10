package controllertests

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileSystem(t *testing.T) {
	testServer.RefreshFileSystem()

	testData := []byte("some file data")
	err := testServer.Server.FileSystem.UpsertFileRaw(1, bytes.NewReader(testData))
	require.NoError(t, err)

	data, err := testServer.Server.FileSystem.GetFileRaw(1)
	require.NoError(t, err)

	b := make([]byte, len(testData))
	_, err = data.Read(b)
	require.NoError(t, err)
	assert.Equal(t, testData, b)

	newData := []byte("different file data")
	err = testServer.Server.FileSystem.UpsertFileRaw(1, bytes.NewReader(newData))
	require.NoError(t, err)

	data, err = testServer.Server.FileSystem.GetFileRaw(1)
	require.NoError(t, err)

	b = make([]byte, len(newData))
	_, err = data.Read(b)
	require.NoError(t, err)
	assert.Equal(t, newData, b)

	_, err = testServer.Server.FileSystem.Stat("1")
	require.True(t, os.IsNotExist(err))
}
