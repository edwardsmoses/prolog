package log

import (
	"fmt"
	"io"
	"os"
	"testing"

	api "github.com/edwardsmoses/proglog/api/v1"
	"github.com/stretchr/testify/require"
)

func TestSegment(t *testing.T) {

	// create a temporary directory to store the segment files
	dir, _ := os.MkdirTemp("", "segment_test")
	defer os.RemoveAll(dir)

	want := &api.Record{Value: []byte("hello world")}

	//configure the segment
	c := Config{}
	c.Segment.MaxIndexBytes = 1024
	c.Segment.MaxStoreBytes = entWidth * 3

	s, err := newSegment((dir), 16, c)
	require.NoError(t, err)

	require.Equal(t, uint64(16), s.nextOffset, s.nextOffset)
	require.False(t, s.IsMaxed())

	for i := uint64(0); i < 3; i++ {
		off, err := s.Append(want)
		require.NoError(t, err)
		require.Equal(t, 16+i, off)

		fmt.Println("is there err when apedning", err, s.IsMaxed(), s.store.size, c.Segment)

		got, err := s.Read(off)
		require.NoError(t, err)
		require.Equal(t, want.Value, got.Value)
	}

	_, err = s.Append(want)
	fmt.Println("error when apedning", err, s.IsMaxed(), s.store.size, c.Segment)
	require.Equal(t, io.EOF, err)

	//maxed index
	require.True(t, s.IsMaxed())

	c.Segment.MaxStoreBytes = uint64(len(want.Value) * 3)
	c.Segment.MaxIndexBytes = 1024

	s, err = newSegment((dir), 16, c)
	require.NoError(t, err)

	//maxed store
	require.True(t, s.IsMaxed())

	err = s.Remove()
	require.NoError(t, err)

	s, err = newSegment((dir), 16, c)
	require.NoError(t, err)
	require.False(t, s.IsMaxed())

}