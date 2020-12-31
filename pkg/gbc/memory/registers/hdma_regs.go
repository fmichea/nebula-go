package registers

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory/lib"
	registerslib "nebula-go/pkg/gbc/memory/registers/lib"
)

type hdma5Reg struct {
	mmu lib.MemoryIO

	inner registerslib.Byte

	bit7   registerslib.Flag
	length registerslib.BitProxy

	index uint8

	hdma1 registerslib.Byte
	hdma2 registerslib.Byte
	hdma3 registerslib.Byte
	hdma4 registerslib.Byte
}

func NewHDMA1Reg(ptr *uint8) registerslib.Byte {
	return registerslib.NewByte(ptr, 0x00)
}

func NewHDMA2Reg(ptr *uint8) registerslib.Byte {
	return registerslib.NewByteWithMask(ptr, 0x00, 0xF0)
}

func NewHDMA3Reg(ptr *uint8) registerslib.Byte {
	return registerslib.NewByteWithMask(ptr, 0x00, 0x1F)
}

func NewHDMA4Reg(ptr *uint8) registerslib.Byte {
	return registerslib.NewByteWithMask(ptr, 0x00, 0xF0)
}

type HDMA5Reg interface {
	registerslib.ByteWithErrors

	MaybeDoHDMA() error
	SourceAddr() uint16
	DestAddr() uint16
}

func NewHDMA5Reg(ptr *uint8, mmu lib.MemoryIO, hdma1, hdma2, hdma3, hdma4 registerslib.Byte) HDMA5Reg {
	reg := registerslib.NewByte(ptr, 0x00)

	return &hdma5Reg{
		mmu: mmu,

		inner: reg,

		bit7:   registerslib.NewFlag(reg, 7),
		length: registerslib.NewBitProxy(reg, 0, 0x7F),

		index: 0,

		hdma1: hdma1,
		hdma2: hdma2,
		hdma3: hdma3,
		hdma4: hdma4,
	}
}

func (r *hdma5Reg) Get() (uint8, error) {
	if r.isActive() {
		return r.length.Get() - r.index, nil
	}
	return r.inner.Get(), nil
}

func (r *hdma5Reg) Set(value uint8) error {
	r.inner.Set(value)

	if r.bit7.GetBool() {
		// Bit 7 == 1, H-Blank DMA switch. This actually sets the bit7 to 0, keeps the length, and for each H-Blank
		// the length will be decremented until everything is transferred.
		r.bit7.SetBool(false)
		r.index = 0
		return nil
	} else {
		// Bit 7 == 0, General Purpose DMA. This transfers everything at once and sets HDMA5 to FFh.
		length := (uint(r.length.Get()) + 1) * 0x10
		r.inner.Set(0xFF)
		return r.doCopyOfLength(r.SourceAddr(), r.DestAddr(), length)
	}
}

func (r *hdma5Reg) MaybeDoHDMA() error {
	if !r.isActive() {
		return nil
	}

	offset := uint16(r.index) * 0x10
	r.index++

	if err := r.doCopyOfLength(r.SourceAddr()+offset, r.DestAddr()+offset, 0x10); err != nil {
		return err
	}

	if r.length.Get() < r.index {
		r.inner.Set(0xFF)
	}
	return nil
}

func (r *hdma5Reg) isActive() bool {
	return !r.bit7.GetBool()
}

func (r *hdma5Reg) doCopyOfLength(sourceAddr, destAddr uint16, length uint) error {
	values, err := r.mmu.ReadByteSlice(sourceAddr, length)
	if err != nil {
		return err
	}

	return r.mmu.WriteByteSlice(destAddr, values)
}

func (r *hdma5Reg) SourceAddr() uint16 {
	// FIXME: add check that address is in ROM/RAM?
	return bitwise.ConvertHighLow8To16(r.hdma1.Get(), r.hdma2.Get())
}

func (r *hdma5Reg) DestAddr() uint16 {
	return 0x8000 | bitwise.ConvertHighLow8To16(r.hdma3.Get(), r.hdma4.Get())
}
