package segments

type AddressRange struct {
	Start uint16
	End   uint16
}

func (r AddressRange) size() uint {
	return uint(r.End - r.Start + 1)
}

func (r AddressRange) containsAddress(addr uint16) bool {
	return r.Start <= addr && addr <= r.End
}

func (r AddressRange) transposeAddress(other AddressRange, addr uint16) uint16 {
	return other.Start + uint16(r.asOffset(addr))
}

func (r AddressRange) asOffset(addr uint16) uint {
	return uint(addr - r.Start)
}

func (r AddressRange) hasCapacityFromAddress(addr uint16, count uint) bool {
	return r.asOffset(addr)+count <= r.size()
}
