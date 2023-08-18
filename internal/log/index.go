package log

import (
	"fmt"
	"io"
	"os"

	"github.com/tysonmote/gommap"
)

var (
	offWidth uint64 = 4
	posWidth uint64 = 8
	entWidth        = offWidth + posWidth
)

type index struct {
	file *os.File
	mmap gommap.MMap
	size uint64
}

func newIndex(f *os.File, c Config) (*index, error) {
	idx := &index{
		file: f,
	}

	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	idx.size = uint64(fi.Size())
	if err := os.Truncate(f.Name(), int64(c.Segment.MaxIndexBytes)); err != nil {
		return nil, err
	}

	if idx.mmap, err = gommap.Map(
		idx.file.Fd(),
		gommap.PROT_READ|gommap.PROT_WRITE,
		gommap.MAP_SHARED,
	); err != nil {
		return nil, err
	}

	return idx, nil
}

func (i *index) Close() error {
	if err := i.mmap.Sync(gommap.MS_SYNC); err != nil {
		return err
	}

	if err := i.file.Sync(); err != nil {
		return err
	}

	if err := i.file.Truncate((int64)(i.size)); err != nil {
		return err
	}

	return i.file.Close()
}

func (i *index) Read(in int64) (out uint32, pos uint64, err error) {
	if i.size == 0 {
		return 0, 0, io.EOF
	}

	// if in is -1, then we want to read the last entry in the index
	if in == -1 {
		out = uint32((uint32(i.size) / uint32(entWidth)) - uint32(1))
	} else {
		out = uint32(in) // otherwise, we want to read the entry at the given index
	}

	// calculate the position of the entry in the index
	pos = uint64(out) * entWidth
	if i.size < pos+entWidth {
		return 0, 0, io.EOF // if the position is greater than the size of the index, then we return an error
	}

	// read the offset and position of the entry
	out = enc.Uint32(i.mmap[pos : pos+offWidth])
	pos = enc.Uint64(i.mmap[pos+offWidth : pos+entWidth])

	//return the offset and position
	return out, pos, nil
}

func (i *index) Write(off uint32, pos uint64) error {

	//we're here maybe?
	fmt.Println("index.go: Write()", len(i.mmap), i.size, entWidth, "result: ", uint64(len(i.mmap)) < i.size+entWidth)

	// check if the mmap is large enough to hold the entry
	if uint64(len(i.mmap)) < i.size+entWidth {
		return io.EOF
	}

	// encode the offset and position
	enc.PutUint32(i.mmap[i.size:i.size+offWidth], off)

	// encode the position
	enc.PutUint64((i.mmap[i.size+offWidth : i.size+entWidth]), pos)

	// increase the size of the index
	i.size += uint64(entWidth)
	return nil
}

func (i *index) Name() string {
	return i.file.Name()
}
