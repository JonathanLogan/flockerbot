package ringbuffer

// ByteSliceBuffer implements a ringbuffer over bytes
type ByteSliceBuffer struct {
	Buffer *Buffer
	data   [][]byte
}

// NewByteSliceBuffer returns a ringbuffer over bytes with size.
func NewByteSliceBuffer(size uint64) *ByteSliceBuffer {
	return &ByteSliceBuffer{
		Buffer: New(size, 0),
		data:   make([][]byte, size),
	}
}

// Push adds byteslice b to the buffer, returning the number of bytes added.
func (bb *ByteSliceBuffer) Push(b ...[]byte) (n int, ok bool) {
	return bb.PushSlice(b)
}

// PushSlice adds byteslices from slice b to the buffer, returning the number of bytes added.
func (bb *ByteSliceBuffer) PushSlice(b [][]byte) (n int, ok bool) {
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

// PopSlice up to len(b) byteslices from the buffer, returning the number of bytes read.
func (bb *ByteSliceBuffer) PopSlice(b [][]byte) (n int, ok bool) {
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

// Pop reads one byteslice from the buffer, returning success.
func (bb *ByteSliceBuffer) Pop() (d []byte, ok bool) {
	if pos, ok := bb.Buffer.GetReadPos(); ok {
		return bb.data[pos], true
	}
	return nil, false
}

// GetPos returns the entry at position pos, or !ok if it isnt contained.
func (bb *ByteSliceBuffer) GetPos(pos uint64) ([]byte, bool) {
	if n, ok := bb.Buffer.TransPos(pos); ok {
		return bb.data[n], true
	}
	return nil, false
}

// Resize the byteSlicebuffer to new size. Returns true on success, false otherwise.
// Copies of ByteSliceBuffer data become invalid with this operation.
func (bb *ByteSliceBuffer) Resize(size uint64) (newbuffer *ByteSliceBuffer, ok bool) {
	fill := bb.Buffer.Fill()
	if fill > size {
		return bb, false
	}
	bd := &ByteSliceBuffer{
		Buffer: New(size, fill),
		data:   make([][]byte, size),
	}
	fb, fe, ss, sb, se := bb.Buffer.CutPoints()
	copy(bd.data, bb.data[fb:fe])
	copy(bd.data[ss:], bb.data[sb:se])
	return bd, true
}
