package fields

import (
	"math/big"
	"sync/atomic"
)

type counter uint64

func (counter *counter) increment() {
	atomic.AddUint64((*uint64)(counter), 1)
}

func (counter counter) text() string {
	var wrapper big.Int
	wrapper.SetUint64(uint64(counter))

	return wrapper.Text(62)
}
