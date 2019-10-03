package cartridge

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadData(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "nebula-go-")
	require.NoError(t, err)

	defer func() {
		require.NoError(t, os.Remove(tmpFile.Name()))
	}()

	text := []byte("This is an example!")

	_, err = tmpFile.Write(text)
	require.NoError(t, err)
	require.NoError(t, tmpFile.Close())

	data, err := loadData(tmpFile.Name())
	require.NoError(t, err)
	require.Equal(t, text, data)
}
