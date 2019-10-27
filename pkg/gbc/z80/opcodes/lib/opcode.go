package opcodeslib

type OpcodeResult struct {
	Size  uint16
	Clock uint16

	Err error
}

type Opcode func() OpcodeResult

func OpcodeSuccess(size, clock uint16) OpcodeResult {
	return OpcodeResult{
		Size:  size,
		Clock: clock,
		Err:   nil,
	}
}

func OpcodeError(err error) OpcodeResult {
	return OpcodeResult{
		Size:  0,
		Clock: 0,
		Err:   err,
	}
}
