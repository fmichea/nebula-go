package registerslib

import (
	"sync/atomic"
)

type threadSafeByte struct {
	d32  uint32
	mask uint8
}

func NewThreadSafeByte(value uint8) Byte {
	return NewThreadSafeByteWithMask(value, 0xFF)
}

func NewThreadSafeByteWithMask(value, mask uint8) Byte {
	reg := &threadSafeByte{
		d32:  0,
		mask: mask,
	}
	reg.SetNoMask(value)
	return reg
}

func (b *threadSafeByte) SetNoMask(value uint8) {
	atomic.StoreUint32(&b.d32, uint32(value))
}

func (b *threadSafeByte) Set(value uint8) {
	currentValue := b.Get()
	newValue := setMasked(value, currentValue, b.mask)
	atomic.CompareAndSwapUint32(&b.d32, uint32(currentValue), uint32(newValue))
}

func (b *threadSafeByte) Get() uint8 {
	return uint8(atomic.LoadUint32(&b.d32))
}
