package segments

type addressRange struct {
	start uint16
	end   uint16
}

func (r addressRange) size() uint16 {
	return r.end - r.start + 1
}

func (r addressRange) containsAddress(addr uint16) bool {
	return r.start <= addr && addr <= r.end
}

func (r addressRange) transposeAddress(other addressRange, addr uint16) uint16 {
	return other.start + r.asOffset(addr)
}

func (r addressRange) asOffset(addr uint16) uint16 {
	return addr - r.start
}
