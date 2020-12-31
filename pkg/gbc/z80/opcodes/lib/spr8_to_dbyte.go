package opcodeslib

import (
	"nebula-go/pkg/common/bitwise"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/z80/registers"
	registerslib "nebula-go/pkg/gbc/z80/registers/lib"
)

// Half-carry and carry are really odd here. See explanation from blargg below original found at URL:
//   http://forums.nesdev.com/viewtopic.php?p=42138
//
// > E8 ADD SP,+$00
// > F8 LD HL,SP+$00
// >
// > Both of these set carry and half-carry based on the low byte of SP added to the UNSIGNED immediate byte. The
// > Negative and Zero flags are always cleared. They also calculate SP + SIGNED immediate byte and put the result
// > into SP or HL, respectively.
func SPR8ToDByte(mmu memory.MMU, regs *registers.Registers, reg registerslib.DByte, clock uint16) OpcodeResult {
	sp := regs.SP.Get()

	d8, err := mmu.ReadByte(regs.PC + 1)
	if err != nil {
		return OpcodeError(err)
	}

	reg.Set(AddRelativeConst(sp, d8))

	sp8 := bitwise.Mask16(sp, 0xFF)
	carryResult := sp8 + uint16(d8)

	regs.F.Set(registers.FlagsCleared)
	regs.F.HC.SetBool(bitwise.Mask16(carryResult, 0xF) < bitwise.Mask16(sp8, 0xF))
	regs.F.CY.SetBool(bitwise.Mask16(carryResult, 0xFF) < bitwise.Mask16(sp8, 0xFF))

	return OpcodeSuccess(2, clock)
}
