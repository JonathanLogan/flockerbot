// Package fixbuffer implements a buffered reader over a fixed size buffer
package fixbuffer

import (
	"errors"
	"io"

	"github.com/JonathanLogan/flockerbot/ringbuffer"
)

var (
	// ErrNotFound is returned when a search does not return data
	ErrNotFound = errors.New("fixbuffer: Seperator not found")
	// ErrBuffer is returned if there is an issue with the buffer
	ErrBuffer = errors.New("fixbuffer: Buffer error")
	// ErrFull is returned if the buffer is filled
	ErrFull = errors.New("fixbuffer: Buffer full")
)

// FixBuffer implements a buffered reader over a fixed size buffer
type FixBuffer struct {
	r   io.Reader
	b   *ringbuffer.ByteBuffer
	s   uint64 // skip
	Sep []byte // Read until encountering Sep.
}

// New returns a FixBuffer of size oveErrBufferFullr the reader.
// Sep is the seperator for reads.
func New(reader io.Reader, size int, sep []byte) *FixBuffer {
	return &FixBuffer{
		b:   ringbuffer.NewByteBuffer(uint64(size)),
		r:   reader,
		Sep: sep,
	}
}

// ReadBytes returns a byteslice up to and including Sep.
// If Sep isnt found, it returns nil and ErrNotFound.
// If the buffer is filled without finding Sep, ErrFull is returned.
func (fb *FixBuffer) ReadBytes() (d []byte, err error) {
	d, err = fb.FindSep()
	if d != nil {
		return d, nil
	}
	if err != ErrNotFound {
		return nil, err
	}
	nl := fb.b.Buffer.Available()
	if nl <= 1 {
		return nil, ErrFull
	}
	nd := make([]byte, nl)
	n, err := fb.r.Read(nd)
	if n > 0 {
		_, ok := fb.b.PushSlice(nd[:n])
		if !ok {
			return nil, ErrBuffer
		}
	}
	d, err2 := fb.FindSep()
	if d != nil {
		return d, err
	}
	if err != nil {
		return nil, err
	}
	return nil, err2
}

// FindSep finds the seperator in the buffer and returns all bytes up and including
// the seperator. If not found, it returns ErrNotFound. Should
func (fb *FixBuffer) FindSep() (d []byte, err error) {
	var n uint64
	var ok bool
	if n, ok = fb.b.FindByte(fb.Sep, fb.s); ok {
		fb.s = 0
		d, ok := fb.b.CutToPos(n)
		if !ok {
			return nil, ErrBuffer
		}
		return d, nil
	}
	// fb.s = n
	return nil, ErrNotFound
}
