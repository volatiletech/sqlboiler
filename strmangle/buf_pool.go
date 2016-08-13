package strmangle

import (
	"bytes"
	"sync"
)

var bufPool = sync.Pool{
	New: newBuffer,
}

func newBuffer() interface{} {
	return &bytes.Buffer{}
}

// GetBuffer retrieves a buffer from the buffer pool
func GetBuffer() *bytes.Buffer {
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()

	return buf
}

// PutBuffer back into the buffer pool
func PutBuffer(buf *bytes.Buffer) {
	bufPool.Put(buf)
}
