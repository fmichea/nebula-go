package registers

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
	return (p.reg.Get() >> p.offset) & p.mask
}

func (p *bitproxy) Set(value uint8) {
	value = (value & p.mask) << p.offset

	current := p.reg.Get()
	current &= 0xFF ^ (p.mask << p.offset)
	current |= value

	p.reg.Set(current)
}
