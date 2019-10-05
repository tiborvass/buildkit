package snapshot

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/containerd/containerd/mount"
	"github.com/stretchr/testify/require"
)

func TestLocalMounterWithMounts(t *testing.T) {
	src, err := ioutil.TempDir("", "testlocalmounter-src")
	require.NoError(t, err)
	defer os.RemoveAll(src)
	err = ioutil.WriteFile(filepath.Join(src, "foo"), []byte("bar"), 0644)
	require.NoError(t, err)

	m := LocalMounterWithMounts([]mount.Mount{
		{
			Type:   "bind",
			Source: src,
		},
	})

	dest, err := m.Mount()
	require.NoError(t, err)

	b, err := ioutil.ReadFile(filepath.Join(dest, "foo"))
	require.NoError(t, err)
	require.Equal(t, string(b), "bar")
}
