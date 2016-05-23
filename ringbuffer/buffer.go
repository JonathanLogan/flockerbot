// Package ringbuffer implements ringbuffer calculations. It is limited to a maximum of maxUint64 operations
// and will panic then.
package ringbuffer

// Buffer implements a ring buffer calculator
type Buffer struct {
	size  uint64 // total size of the buffer
	read  uint64 // position to read from
	write uint64 // position to write to on next write
	empty bool   // if buffer is empty
}

// New returns a new buffer of size. Fill is the current write position (when copying buffers).
func New(size, fill uint64) *Buffer {
	if fill > 0 {
		return &Buffer{
			size:  size,
			empty: false,
			write: fill,
		}
	}
	return &Buffer{
		size:  size,
		empty: true,
	}
}

// GetReadPos returns the next read position, or 0,false.
func (buf *Buffer) GetReadPos() (uint64, bool) {
	if buf.Fill() > 0 {
		pr := buf.read % buf.size
		buf.read++
		return pr, true
	}
	return 0, false
}

// GetWritePos returns the next write position, or 0,false.
func (buf *Buffer) GetWritePos() (uint64, bool) {
	if buf.Fill() < buf.size {
		pw := buf.write % buf.size
		buf.empty = false
		buf.write++
		return pw, true
	}
	return 0, false
}

// Fill returns the number of elements in the buffer.
func (buf *Buffer) Fill() uint64 {
	if buf.empty {
		return 0
	}
	return buf.write - buf.read
}

// Available returns the bytes available in the buffer for writing.
func (buf *Buffer) Available() uint64 {
	return buf.size - buf.Fill()
}

// Stat returns the status of a buffer.
func (buf *Buffer) Stat() (size, readpos, writepos uint64) {
	return buf.size, buf.read, buf.write
}

// Set readpos and writepos.
func (buf *Buffer) Set(readpos, writepos uint64) {
	buf.read = readpos
	buf.write = writepos
}

// TransPos translates positiong to the correct address. Returns false if not contained in ring.
func (buf *Buffer) TransPos(pos uint64) (uint64, bool) {
	if pos >= buf.read && pos < buf.write {
		return pos % buf.size, true
	}
	return pos % buf.size, false
}

// CutPoints returns the cut points for reslicing.
// The cutpoints can be used for simple reslicing for grow/shrink operations like this:
//
// 		buffer_new:=New(newsize,buffer_old.Fill())
// 		firstSliceBegin, firstSliceEnd, secondSliceSkip, secondSliceBegin, secondSliceEnd=huffer_old.CutPoints()
// 		copy(buffer_new_data, buffer_old_data[firstSliceBegin:firstSliceEnd]
// 		copy(buffer_new_data[secondSliceSkip:], buffer_old_data[secondSliceBegin:secondSliceEnd])
//
// This results in a new ring buffer that has the correct read/write positions set.
func (buf *Buffer) CutPoints() (firstSliceBegin, firstSliceEnd, secondSliceSkip, secondSliceBegin, secondSliceEnd uint64) {
	pw := (buf.write % buf.size)
	pr := (buf.read % buf.size)
	if pw > pr {
		return pr, pw, 0, 0, 0
	}
	return pr, buf.size, (buf.size - pr), 0, pw
}
