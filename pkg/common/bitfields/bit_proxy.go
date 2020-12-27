package bitfields

import (
	"nebula-go/pkg/common/bitwise"
)

type ROBitProxy interface {
	Get() uint8
}

type BitProxy interface {
	Get() uint8
	Set(value uint8)
}

type bitproxy struct {
	reg    Byte
	offset uint8
	mask   uint8
}

func NewBitProxy(reg Byte, offset, mask uint8) BitProxy {
	return &bitproxy{
		reg:    reg,
		offset: offset,
		mask:   mask,
	}
}

func (p *bitproxy) Get() uint8 {
	return bitwise.GetBits8(p.reg.Get(), p.offset, p.mask)
}

func (p *bitproxy) Set(value uint8) {
	p.reg.Set(bitwise.SetBits8(p.reg.Get(), p.offset, p.mask, value))
}
