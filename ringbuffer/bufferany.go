package ringbuffer

// Any type you want
type Any interface{}

// AnyBuffer implements a ringbuffer over bytes
type AnyBuffer struct {
	Buffer *Buffer
	data   []Any
}

// NewAnyBuffer returns a ringbuffer over any interface with size.
func NewAnyBuffer(size uint64) *AnyBuffer {
	return &AnyBuffer{
		Buffer: New(size, 0),
		data:   make([]Any, size),
	}
}

// Push adds b to the buffer, returning the number of elements added.
func (bb *AnyBuffer) Push(b ...Any) (n int, ok bool) {
	return bb.PushSlice(b)
}

// PushSlice adds elements from slice b to the buffer, returning the number of elemnts added.
func (bb *AnyBuffer) PushSlice(b []Any) (n int, ok bool) {
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

// PopSlice up to len(b) elements from the buffer, returning the number of elements read.
func (bb *AnyBuffer) PopSlice(b []Any) (n int, ok bool) {
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

// Pop reads one element from the buffer, returning success.
func (bb *AnyBuffer) Pop() (d Any, ok bool) {
	if pos, ok := bb.Buffer.GetReadPos(); ok {
		return bb.data[pos], true
	}
	return nil, false
}

// GetPos returns the entry at position pos, or !ok if it isnt contained.
func (bb *AnyBuffer) GetPos(pos uint64) (Any, bool) {
	if n, ok := bb.Buffer.TransPos(pos); ok {
		return bb.data[n], true
	}
	return nil, false
}

// Resize the AnyBuffer to new size. Returns true on success, false otherwise.
// Copies of AnyBuffer data become invalid with this operation.
func (bb *AnyBuffer) Resize(size uint64) (newbuffer *AnyBuffer, ok bool) {
	fill := bb.Buffer.Fill()
	if fill > size {
		return bb, false
	}
	bd := &AnyBuffer{
		Buffer: New(size, fill),
		data:   make([]Any, size),
	}
	fb, fe, ss, sb, se := bb.Buffer.CutPoints()
	copy(bd.data, bb.data[fb:fe])
	copy(bd.data[ss:], bb.data[sb:se])
	return bd, true
}
