package ringbuffer

import (
	"bytes"
	"testing"
)

func TestAny(t *testing.T) {
	td := []byte("Something")
	buf := NewAnyBuffer(10)
	buf.Push(td)
	x, _ := buf.Pop()
	td2 := x.([]byte)
	if !bytes.Equal(td, td2) {
		t.Error("No match")
	}
}
