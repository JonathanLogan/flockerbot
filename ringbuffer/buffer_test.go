// Package ringbuffer implements ringbuffer calculations
package ringbuffer

import "testing"

func TestBuffer(t *testing.T) {
	b := New(3, 0)
	if _, ok := b.GetReadPos(); ok {
		t.Error("Read from empty must fail")
	}
	if b.Fill() != 0 {
		t.Error("Fill must be zero on empty ring")
	}
	n, ok := b.GetWritePos()
	if !ok || n != 0 {
		t.Errorf("Write 1, expected 0: %d", n)
	}
	if b.Fill() != 1 {
		t.Error("Fill(1) != 1")
	}
	n, ok = b.GetReadPos()
	if !ok || n != 0 {
		t.Errorf("Read 1, expected 0: %d", n)
	}
	if b.Fill() != 0 {
		t.Error("Fill(2) != 0")
	}
	_, ok = b.GetReadPos()
	if ok {
		t.Error("Read 2, expected error")
	}
	n, ok = b.GetWritePos()
	if !ok || n != 1 {
		t.Errorf("Write 2, expected 1: %d", n)
	}
	if b.Fill() != 1 {
		t.Error("Fill(3) != 1")
	}
	n, ok = b.GetWritePos()
	if !ok || n != 2 {
		t.Errorf("Write 3, expected 2: %d", n)
	}
	if b.Fill() != 2 {
		t.Error("Fill(4) != 2")
	}
	n, ok = b.GetWritePos()
	if !ok || n != 0 {
		t.Errorf("Write 4, expected 0: %d", n)
	}
	if b.Fill() != 3 {
		t.Errorf("Fill(5) != 3. %d", b.Fill())
	}
	n, ok = b.GetReadPos()
	if !ok || n != 1 {
		t.Errorf("Read 1, expected 1: %d", n)
	}

	for i := 0; i < 100; i++ {
		n1, ok1 := b.GetReadPos()
		n2, ok2 := b.GetWritePos()
		if !(ok1 && ok2) {
			t.Fatalf("Fill fail: %d, %d,%d", n1, n2, i)
		}
	}
	for b.Fill() > 0 {
		_, ok := b.GetReadPos()
		if !ok {
			t.Fatal("Deplete failed")
		}
	}
	if _, ok := b.GetReadPos(); ok {
		t.Error("Empty must fail")
	}
	for i := uint64(0); i < b.size; i++ {
		_, ok := b.GetWritePos()
		if !ok {
			t.Fatal("Fill failed")
		}
	}
	if _, ok := b.GetWritePos(); ok {
		t.Error("Full must fail")
	}
}

func TestReslice(t *testing.T) {
	data1 := make([]byte, 10)
	buffer1 := New(uint64(len(data1)), 0)
	for i := 1; i < 13; i++ {
		n, _ := buffer1.GetWritePos()
		data1[n] = byte(i + 100)
		buffer1.GetReadPos()
	}
	buffer1.write = 8
	buffer1.read = 2
	fb, fe, ss, sb, se := buffer1.CutPoints()

	data2 := make([]byte, 15)
	buffer2 := New(uint64(len(data2)), buffer1.Fill())
	copy(data2, data1[fb:fe])
	copy(data2[ss:], data1[sb:se])
	r1, _ := buffer1.GetReadPos()
	r2, _ := buffer2.GetReadPos()
	w2, _ := buffer2.GetWritePos()

	if data1[r1] != data2[r2] {
		t.Error("1: Read/Begin reslice failed")
	}
	if data2[w2] != 0x00 {
		t.Error("1: Write over size")
	}

	buffer1.write = 15
	buffer1.read = 8
	fb, fe, ss, sb, se = buffer1.CutPoints()

	data2 = make([]byte, 15)
	buffer2 = New(uint64(len(data2)), buffer1.Fill())
	copy(data2, data1[fb:fe])
	copy(data2[ss:], data1[sb:se])
	r1, _ = buffer1.GetReadPos()
	r2, _ = buffer2.GetReadPos()
	w2, _ = buffer2.GetWritePos()

	if data1[r1] != data2[r2] {
		t.Error("2: Read/Begin reslice failed")
	}
	if data2[w2] != 0x00 {
		t.Error("2: Write over size")
	}

}
