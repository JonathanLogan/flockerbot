package ringbuffer

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBytes(t *testing.T) {
	td1 := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	td2 := []byte{0x07, 0x08, 0x09, 0x10, 0x11}
	buf := NewByteBuffer(10)
	n, _ := buf.PushSlice(td1)
	if n != 6 {
		t.Errorf("Expected to write 6, wrote %d", n)
	}
	n, _ = buf.PushSlice(td2)
	if n != 4 {
		t.Errorf("Expected to write 4, wrote %d", n)
	}
	rb := make([]byte, 6)
	n, _ = buf.PopSlice(rb)
	if n != 6 {
		t.Errorf("Expected to read 6, read %d", n)
	}
	if !bytes.Equal(rb, td1) {
		t.Error("Data error first read")
	}
	rb = make([]byte, 4)
	n, _ = buf.PopSlice(rb)
	if n != 4 {
		t.Errorf("Expected to read 4, read %d", n)
	}
	if !bytes.Equal(rb, td2[0:4]) {
		t.Error("Data error second read")
	}
	n, _ = buf.PushSlice(td1)
	if n != 6 {
		t.Errorf("Expected to write 6, wrote %d", n)
	}
	fill := buf.Buffer.Fill()
	var ok bool
	if buf, ok = buf.Resize(fill); !ok {
		t.Fatal("Resize to fill failed")
	}
	rb = make([]byte, 6)
	n, _ = buf.PopSlice(rb)
	if n != 6 {
		t.Errorf("Expected to read 6, read %d", n)
	}
	if !bytes.Equal(rb, td1) {
		t.Error("Data error third read")
	}
	size, _, _ := buf.Buffer.Stat()
	if size != fill {
		t.Errorf("New buffer did not overwrite old: %d", size)
	}
}

func TestBytesFind(t *testing.T) {
	td1 := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05}
	buf := NewByteBuffer(10)
	buf.PushSlice(td1)
	n, ok := buf.FindByte([]byte{0x01}, 0)
	if n != 1 || !ok {
		t.Errorf("0x01 on position 1 not found: %d", n)
	}
	d, ok := buf.CutToPos(n)
	if !ok {
		t.Error("Cutting error")
	}
	if !bytes.Equal(d, td1[0:2]) {
		t.Errorf("Data error: %+v", d)
	}
	b, _ := buf.Pop()
	if b != td1[2] {
		t.Error("Read pointer advance error")
	}
	n, ok = buf.FindByte([]byte{0x04}, 0)
	if !ok {
		t.Errorf("Follow read not found: %d", n)
	}
	d, _ = buf.CutToPos(n)
	if !bytes.Equal(d, td1[3:5]) {
		t.Error("Follow read data error")
	}
}

func TestCutToPos(t *testing.T) {
	var d []byte
	fmt.Println()
	buf := NewByteBuffer(10)
	buf.PushSlice(bytes.Repeat([]byte{0x03}, 7))
	d, _ = buf.CutToPos(5)
	fmt.Println(d)
	buf.PushSlice(bytes.Repeat([]byte{0x02}, 8))
	d, _ = buf.CutToPos(11)
	fmt.Println(d)
	d, _ = buf.CutToPos(14)
	fmt.Println(d)
}
