package ringbuffer

// StringBuffer implements a ringbuffer over bytes
type StringBuffer struct {
	Buffer *Buffer
	data   []string
}

// NewStringBuffer returns a ringbuffer over bytes with size.
func NewStringBuffer(size uint64) *StringBuffer {
	return &StringBuffer{
		Buffer: New(size, 0),
		data:   make([]string, size),
	}
}

// Push adds string b to the buffer, returning the number of bytes added.
func (bb *StringBuffer) Push(b ...string) (n int, ok bool) {
	return bb.PushSlice(b)
}

// PushSlice adds stringslice from slice b to the buffer, returning the number of bytes added.
func (bb *StringBuffer) PushSlice(b []string) (n int, ok bool) {
	for _, x := range b {
		if pos, ok := bb.Buffer.GetWritePos(); ok {
			bb.data[pos] = x
			n++
		} else {
			return n, false
		}
	}
	return n, true
}

// PopSlice up to len(b) stringslice from the buffer, returning the number of bytes read.
func (bb *StringBuffer) PopSlice(b []string) (n int, ok bool) {
	for wpos := range b {
		if pos, ok := bb.Buffer.GetReadPos(); ok {
			b[wpos] = bb.data[pos]
			n++
		} else {
			return n, false
		}
	}
	return n, true
}

// Pop reads one string from the buffer, returning success.
func (bb *StringBuffer) Pop() (d string, ok bool) {
	if pos, ok := bb.Buffer.GetReadPos(); ok {
		return bb.data[pos], true
	}
	return "", false
}

// GetPos returns the entry at position pos, or !ok if it isnt contained.
func (bb *StringBuffer) GetPos(pos uint64) (string, bool) {
	if n, ok := bb.Buffer.TransPos(pos); ok {
		return bb.data[n], true
	}
	return "", false
}

// Resize the Stringbuffer to new size. Returns true on success, false otherwise.
// Copies of StringBuffer data become invalid with this operation.
func (bb *StringBuffer) Resize(size uint64) (newbuffer *StringBuffer, ok bool) {
	fill := bb.Buffer.Fill()
	if fill > size {
		return bb, false
	}
	bd := &StringBuffer{
		Buffer: New(size, fill),
		data:   make([]string, size),
	}
	fb, fe, ss, sb, se := bb.Buffer.CutPoints()
	copy(bd.data, bb.data[fb:fe])
	copy(bd.data[ss:], bb.data[sb:se])
	return bd, true
}
