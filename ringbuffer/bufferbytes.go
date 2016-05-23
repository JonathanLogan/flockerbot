package ringbuffer

// ByteBuffer implements a ringbuffer over bytes
type ByteBuffer struct {
	Buffer *Buffer
	data   []byte
}

// NewByteBuffer returns a ringbuffer over bytes with size.
func NewByteBuffer(size uint64) *ByteBuffer {
	return &ByteBuffer{
		Buffer: New(size, 0),
		data:   make([]byte, size),
	}
}

// Push adds bytes b to the buffer, returning the number of bytes added.
func (bb *ByteBuffer) Push(b ...byte) (n int, ok bool) {
	return bb.PushSlice(b)
}

// PushSlice adds bytes from slice b to the buffer, returning the number of bytes added.
func (bb *ByteBuffer) PushSlice(b []byte) (n int, ok bool) {
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

// PopSlice up to len(b) bytes from the buffer, returning the number of bytes read.
func (bb *ByteBuffer) PopSlice(b []byte) (n int, ok bool) {
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

// Pop reads one byte from the buffer, returning success.
func (bb *ByteBuffer) Pop() (d byte, ok bool) {
	if pos, ok := bb.Buffer.GetReadPos(); ok {
		return bb.data[pos], true
	}
	return 0x00, false
}

// GetPos returns the entry at position pos, or !ok if it isnt contained.
func (bb *ByteBuffer) GetPos(pos uint64) (byte, bool) {
	if n, ok := bb.Buffer.TransPos(pos); ok {
		return bb.data[n], true
	}
	return 0x00, false
}

// FindByte finds s in buffer and returns the position. If the byte was
// not found, return the number of bytes searched.
func (bb *ByteBuffer) FindByte(sep []byte, skip uint64) (uint64, bool) {
	var i uint64
	l := uint64(len(sep))
	if bb.Buffer.Fill() == 0 {
		return 0, false
	}
FirstFind:
	for i = bb.Buffer.read + skip; i < bb.Buffer.write; i++ {
		if bb.data[i%bb.Buffer.size] == sep[0] {
			r := int64(-1)
			for j := uint64(0); (i+j < bb.Buffer.write) && (j < l); j++ {
				if bb.data[(i+j)%bb.Buffer.size] != sep[j] {
					continue FirstFind
				}
				r++
			}
			return i + uint64(r), true
		}
	}
	return i - bb.Buffer.read, false // Not found
}

// CutToPos returns a byteslice containing all bytes up to and including pos. Returns false on failure
func (bb *ByteBuffer) CutToPos(pos uint64) ([]byte, bool) {
	if pos >= bb.Buffer.write {
		return nil, false
	}
	if pos < bb.Buffer.read {
		return nil, false
	}
	l := 1 + pos - bb.Buffer.read
	d := make([]byte, l)
	for x := uint64(0); x < l; x++ {
		d[x] = bb.data[(bb.Buffer.read+x)%bb.Buffer.size]
	}
	bb.Buffer.read += l
	return d, true
}

// Resize the bytebuffer to new size. Returns true on success, false otherwise.
// Copies of ByteBuffer data become invalid with this operation.
func (bb *ByteBuffer) Resize(size uint64) (newbuffer *ByteBuffer, ok bool) {
	fill := bb.Buffer.Fill()
	if fill > size {
		return bb, false
	}
	bd := &ByteBuffer{
		Buffer: New(size, fill),
		data:   make([]byte, size),
	}
	fb, fe, ss, sb, se := bb.Buffer.CutPoints()
	copy(bd.data, bb.data[fb:fe])
	copy(bd.data[ss:], bb.data[sb:se])
	return bd, true
}
