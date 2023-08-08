package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	write = []byte("hello world")
	width = uint64(len(write)) + lenWidth
)

func TestStoreAppendRead(t *testing.T) {
	f, err := ioutil.TempFile("", "store_append_read_test")
	require.NoError(t, err)

	// defer the removal of the file until after the test
	defer os.Remove((f.Name()))

	s, err := newStore(f)
	require.NoError(t, err)

	testAppend(t, s)
	// testRead(t, s)
	// testReadAt(t, s)

	s, err = newStore((f))
	require.NoError(t, err)
	// testRead(t, s)

}

func testAppend(t *testing.T, s *store) {
	t.Helper()
	for i := uint64(1); i < 4; i++ {
		n, pos, err := s.Append(write)
		require.NoError(t, err)
		require.Equal(t, pos+n, i*width)
	}
}
