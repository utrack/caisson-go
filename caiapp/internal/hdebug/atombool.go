package hdebug

import "sync/atomic"

type atomBool struct{ flag int32 }

func (b *atomBool) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), int32(i))
}

func (b *atomBool) Get() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
