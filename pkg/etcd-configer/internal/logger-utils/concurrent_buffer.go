package loggerutils

import (
	"bytes"
	"sync"
)

// ConcurrentBuffer ...
type ConcurrentBuffer struct {
	locker      sync.RWMutex
	innerBuffer bytes.Buffer
}

// Write ...
func (buffer *ConcurrentBuffer) Write(p []byte) (n int, err error) {
	buffer.locker.Lock()
	defer buffer.locker.Unlock()

	return buffer.innerBuffer.Write(p)
}

// String ...
func (buffer *ConcurrentBuffer) String() string {
	buffer.locker.RLock()
	defer buffer.locker.RUnlock()

	return buffer.innerBuffer.String()
}
