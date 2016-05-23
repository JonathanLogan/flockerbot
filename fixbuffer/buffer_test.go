package fixbuffer

import (
	"testing"

	"github.com/davecgh/go-spew/spew"
)

type rtest struct {
	data []byte
}

func (r *rtest) Read(p []byte) (n int, err error) {
	copy(p, r.data)
	if len(p) > len(r.data) {
		return len(r.data), nil
	}
	return len(p), nil
}

func TestBufRead(t *testing.T) {
	r := &rtest{
		data: []byte("Test data\n"),
	}
	buf := New(r, 30, []byte("\n"))
	_ = buf
	spew.Dump(buf.ReadBytes())
}
