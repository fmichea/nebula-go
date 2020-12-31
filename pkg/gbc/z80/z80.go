package z80

import (
	"fmt"
	"os"

	"nebula-go/pkg/gbc/graphics"
	"nebula-go/pkg/gbc/memory"
	z80lib "nebula-go/pkg/gbc/z80/lib"
	"nebula-go/pkg/gbc/z80/opcodes"
	opcodeslib "nebula-go/pkg/gbc/z80/opcodes/lib"
	"nebula-go/pkg/gbc/z80/registers"
)

type CPU struct {
	MMU     memory.MMU
	MMURegs *memory.Registers

	GPU graphics.GPU

	Regs    *registers.Registers
	Opcodes []opcodeslib.Opcode

	noopOpcode        opcodeslib.Opcode
	diOpcode          opcodeslib.Opcode
	inPlaceInterrupts map[z80lib.Interrupt]opcodeslib.Opcode
}

func NewCPU(mmu memory.MMU, gpu graphics.GPU) *CPU {
	regs := registers.New()

	factory := opcodes.NewFactory(mmu, regs)

	noopOpcode := factory.Miscellaneous.Noop()
	diOpcode := factory.Miscellaneous.DI()
	inplaceInterrupts := map[z80lib.Interrupt]opcodeslib.Opcode{}
	for _, i := range []z80lib.Interrupt{z80lib.Rst40h, z80lib.Rst48h, z80lib.Rst50h, z80lib.Rst58h, z80lib.Rst60h} {
		inplaceInterrupts[i] = factory.ControlFlow.CallInterruptInplace(i)
	}

	return &CPU{
		MMU:     mmu,
		MMURegs: mmu.Registers(),

		GPU: gpu,

		Regs: regs,

		// We could probably get smart and generate this automatically, but after the initial toil of writing this,
		// it's pretty much fixed for the whole life of the emulator, can be tested once and then never changed.
		// This has the advantage to be fairly easy to read and verified.
		Opcodes: []opcodeslib.Opcode{
			// 0x00 - 0x07
			noopOpcode,
			factory.Load.ConstToDByte(regs.BC),
			factory.Load.AToBCPtr(),
			factory.ALU.IncrementDByte(regs.BC),
			factory.ALU.IncrementByte(regs.B),
			factory.ALU.DecrementByte(regs.B),
			factory.Load.ConstToByte(regs.B),
			factory.Miscellaneous.RLCA(),

			// 0x08 - 0x0F
			factory.Load.SPToAddress(),
			factory.ALU.AddDByteToHL(regs.BC),
			factory.Load.BCPtrToA(),
			factory.ALU.DecrementDByte(regs.BC),
			factory.ALU.IncrementByte(regs.C),
			factory.ALU.DecrementByte(regs.C),
			factory.Load.ConstToByte(regs.C),
			factory.Miscellaneous.RRCA(),

			// 0x10 - 0x17
			factory.Miscellaneous.Stop(),
			factory.Load.ConstToDByte(regs.DE),
			factory.Load.AToDEPtr(),
			factory.ALU.IncrementDByte(regs.DE),
			factory.ALU.IncrementByte(regs.D),
			factory.ALU.DecrementByte(regs.D),
			factory.Load.ConstToByte(regs.D),
			factory.Miscellaneous.RLA(),

			// 0x18 - 0x1F
			factory.ControlFlow.JumpRelative(),
			factory.ALU.AddDByteToHL(regs.DE),
			factory.Load.DEPtrToA(),
			factory.ALU.DecrementDByte(regs.DE),
			factory.ALU.IncrementByte(regs.E),
			factory.ALU.DecrementByte(regs.E),
			factory.Load.ConstToByte(regs.E),
			factory.Miscellaneous.RRA(),

			// 0x20 - 0x27
			factory.ControlFlow.JumpRelativeIfNot(regs.F.ZF),
			factory.Load.ConstToDByte(regs.HL),
			factory.Load.AToHLInc(),
			factory.ALU.IncrementDByte(regs.HL),
			factory.ALU.IncrementByte(regs.H),
			factory.ALU.DecrementByte(regs.H),
			factory.Load.ConstToByte(regs.H),
			factory.Miscellaneous.DAA(),

			// 0x28 - 0x2F
			factory.ControlFlow.JumpRelativeIf(regs.F.ZF),
			factory.ALU.AddDByteToHL(regs.HL),
			factory.Load.HLIncToA(),
			factory.ALU.DecrementDByte(regs.HL),
			factory.ALU.IncrementByte(regs.L),
			factory.ALU.DecrementByte(regs.L),
			factory.Load.ConstToByte(regs.L),
			factory.Miscellaneous.CPL(),

			// 0x30 - 0x37
			factory.ControlFlow.JumpRelativeIfNot(regs.F.CY),
			factory.Load.ConstToDByte(regs.SP),
			factory.Load.AToHLDec(),
			factory.ALU.IncrementDByte(regs.SP),
			factory.ALU.IncrementHLPtr(),
			factory.ALU.DecrementHLPtr(),
			factory.Load.D8ToHLPtr(),
			factory.Miscellaneous.SCF(),

			// 0x38 - 0x3F
			factory.ControlFlow.JumpRelativeIf(regs.F.CY),
			factory.ALU.AddDByteToHL(regs.SP),
			factory.Load.HLDecToA(),
			factory.ALU.DecrementDByte(regs.SP),
			factory.ALU.IncrementByte(regs.A),
			factory.ALU.DecrementByte(regs.A),
			factory.Load.ConstToByte(regs.A),
			factory.Miscellaneous.CCF(),

			// 0x40 - 0x47
			factory.Load.ByteToByte(regs.B, regs.B),
			factory.Load.ByteToByte(regs.B, regs.C),
			factory.Load.ByteToByte(regs.B, regs.D),
			factory.Load.ByteToByte(regs.B, regs.E),
			factory.Load.ByteToByte(regs.B, regs.H),
			factory.Load.ByteToByte(regs.B, regs.L),
			factory.Load.HLPtrToByte(regs.B),
			factory.Load.ByteToByte(regs.B, regs.A),

			// 0x48 - 0x4F
			factory.Load.ByteToByte(regs.C, regs.B),
			factory.Load.ByteToByte(regs.C, regs.C),
			factory.Load.ByteToByte(regs.C, regs.D),
			factory.Load.ByteToByte(regs.C, regs.E),
			factory.Load.ByteToByte(regs.C, regs.H),
			factory.Load.ByteToByte(regs.C, regs.L),
			factory.Load.HLPtrToByte(regs.C),
			factory.Load.ByteToByte(regs.C, regs.A),

			// 0x50 - 0x57
			factory.Load.ByteToByte(regs.D, regs.B),
			factory.Load.ByteToByte(regs.D, regs.C),
			factory.Load.ByteToByte(regs.D, regs.D),
			factory.Load.ByteToByte(regs.D, regs.E),
			factory.Load.ByteToByte(regs.D, regs.H),
			factory.Load.ByteToByte(regs.D, regs.L),
			factory.Load.HLPtrToByte(regs.D),
			factory.Load.ByteToByte(regs.D, regs.A),

			// 0x58 - 0x5F
			factory.Load.ByteToByte(regs.E, regs.B),
			factory.Load.ByteToByte(regs.E, regs.C),
			factory.Load.ByteToByte(regs.E, regs.D),
			factory.Load.ByteToByte(regs.E, regs.E),
			factory.Load.ByteToByte(regs.E, regs.H),
			factory.Load.ByteToByte(regs.E, regs.L),
			factory.Load.HLPtrToByte(regs.E),
			factory.Load.ByteToByte(regs.E, regs.A),

			// 0x60 - 0x67
			factory.Load.ByteToByte(regs.H, regs.B),
			factory.Load.ByteToByte(regs.H, regs.C),
			factory.Load.ByteToByte(regs.H, regs.D),
			factory.Load.ByteToByte(regs.H, regs.E),
			factory.Load.ByteToByte(regs.H, regs.H),
			factory.Load.ByteToByte(regs.H, regs.L),
			factory.Load.HLPtrToByte(regs.H),
			factory.Load.ByteToByte(regs.H, regs.A),

			// 0x68 - 0x6F
			factory.Load.ByteToByte(regs.L, regs.B),
			factory.Load.ByteToByte(regs.L, regs.C),
			factory.Load.ByteToByte(regs.L, regs.D),
			factory.Load.ByteToByte(regs.L, regs.E),
			factory.Load.ByteToByte(regs.L, regs.H),
			factory.Load.ByteToByte(regs.L, regs.L),
			factory.Load.HLPtrToByte(regs.L),
			factory.Load.ByteToByte(regs.L, regs.A),

			// 0x70 - 0x77
			factory.Load.ByteToHLPtr(regs.B),
			factory.Load.ByteToHLPtr(regs.C),
			factory.Load.ByteToHLPtr(regs.D),
			factory.Load.ByteToHLPtr(regs.E),
			factory.Load.ByteToHLPtr(regs.H),
			factory.Load.ByteToHLPtr(regs.L),
			factory.Miscellaneous.Halt(),
			factory.Load.ByteToHLPtr(regs.A),

			// 0x78 - 0x7F
			factory.Load.ByteToByte(regs.A, regs.B),
			factory.Load.ByteToByte(regs.A, regs.C),
			factory.Load.ByteToByte(regs.A, regs.D),
			factory.Load.ByteToByte(regs.A, regs.E),
			factory.Load.ByteToByte(regs.A, regs.H),
			factory.Load.ByteToByte(regs.A, regs.L),
			factory.Load.HLPtrToByte(regs.A),
			factory.Load.ByteToByte(regs.A, regs.A),

			// 0x80 - 0x87
			factory.ALU.AddByteToA(regs.B),
			factory.ALU.AddByteToA(regs.C),
			factory.ALU.AddByteToA(regs.D),
			factory.ALU.AddByteToA(regs.E),
			factory.ALU.AddByteToA(regs.H),
			factory.ALU.AddByteToA(regs.L),
			factory.ALU.AddHLPtrToA(),
			factory.ALU.AddByteToA(regs.A),

			// 0x88 - 0x8F
			factory.ALU.AdcByteToA(regs.B),
			factory.ALU.AdcByteToA(regs.C),
			factory.ALU.AdcByteToA(regs.D),
			factory.ALU.AdcByteToA(regs.E),
			factory.ALU.AdcByteToA(regs.H),
			factory.ALU.AdcByteToA(regs.L),
			factory.ALU.AdcHLPtrToA(),
			factory.ALU.AdcByteToA(regs.A),

			// 0x90 - 0x97
			factory.ALU.SubByteToA(regs.B),
			factory.ALU.SubByteToA(regs.C),
			factory.ALU.SubByteToA(regs.D),
			factory.ALU.SubByteToA(regs.E),
			factory.ALU.SubByteToA(regs.H),
			factory.ALU.SubByteToA(regs.L),
			factory.ALU.SubHLPtrToA(),
			factory.ALU.SubByteToA(regs.A),

			// 0x98 - 0x9F
			factory.ALU.SbcByteToA(regs.B),
			factory.ALU.SbcByteToA(regs.C),
			factory.ALU.SbcByteToA(regs.D),
			factory.ALU.SbcByteToA(regs.E),
			factory.ALU.SbcByteToA(regs.H),
			factory.ALU.SbcByteToA(regs.L),
			factory.ALU.SbcHLPtrToA(),
			factory.ALU.SbcByteToA(regs.A),

			// 0xA0 - 0xA7
			factory.ALU.AndByteToA(regs.B),
			factory.ALU.AndByteToA(regs.C),
			factory.ALU.AndByteToA(regs.D),
			factory.ALU.AndByteToA(regs.E),
			factory.ALU.AndByteToA(regs.H),
			factory.ALU.AndByteToA(regs.L),
			factory.ALU.AndHLPtrToA(),
			factory.ALU.AndByteToA(regs.A),

			// 0xA8 - 0xAF
			factory.ALU.XorByteToA(regs.B),
			factory.ALU.XorByteToA(regs.C),
			factory.ALU.XorByteToA(regs.D),
			factory.ALU.XorByteToA(regs.E),
			factory.ALU.XorByteToA(regs.H),
			factory.ALU.XorByteToA(regs.L),
			factory.ALU.XorHLPtrToA(),
			factory.ALU.XorByteToA(regs.A),

			// 0xB0 - 0xB7
			factory.ALU.OrByteToA(regs.B),
			factory.ALU.OrByteToA(regs.C),
			factory.ALU.OrByteToA(regs.D),
			factory.ALU.OrByteToA(regs.E),
			factory.ALU.OrByteToA(regs.H),
			factory.ALU.OrByteToA(regs.L),
			factory.ALU.OrHLPtrToA(),
			factory.ALU.OrByteToA(regs.A),

			// 0xB8 - 0xBF
			factory.ALU.CompareByteToA(regs.B),
			factory.ALU.CompareByteToA(regs.C),
			factory.ALU.CompareByteToA(regs.D),
			factory.ALU.CompareByteToA(regs.E),
			factory.ALU.CompareByteToA(regs.H),
			factory.ALU.CompareByteToA(regs.L),
			factory.ALU.CompareHLPtrToA(),
			factory.ALU.CompareByteToA(regs.A),

			// 0xC0 - 0xC7
			factory.ControlFlow.ReturnIfNot(regs.F.ZF),
			factory.Load.PopDByte(regs.BC),
			factory.ControlFlow.JumpIfNot(regs.F.ZF),
			factory.ControlFlow.Jump(),
			factory.ControlFlow.CallIfNot(regs.F.ZF),
			factory.Load.PushDByte(regs.BC),
			factory.ALU.AddD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst00h),

			// 0xC8 - 0xCF
			factory.ControlFlow.ReturnIf(regs.F.ZF),
			factory.ControlFlow.Return(),
			factory.ControlFlow.JumpIf(regs.F.ZF),
			factory.CB.CB(),
			factory.ControlFlow.CallIf(regs.F.ZF),
			factory.ControlFlow.Call(),
			factory.ALU.AdcD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst08h),

			// 0xD0 - 0xD7
			factory.ControlFlow.ReturnIfNot(regs.F.CY),
			factory.Load.PopDByte(regs.DE),
			factory.ControlFlow.JumpIfNot(regs.F.CY),
			nil,
			factory.ControlFlow.CallIfNot(regs.F.CY),
			factory.Load.PushDByte(regs.DE),
			factory.ALU.SubD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst10h),

			// 0xD8 - 0xDF
			factory.ControlFlow.ReturnIf(regs.F.CY),
			factory.ControlFlow.ReturnInterrupt(),
			factory.ControlFlow.JumpIf(regs.F.CY),
			nil,
			factory.ControlFlow.CallIf(regs.F.CY),
			nil,
			factory.ALU.SbcD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst18h),

			// 0xE0 - 0xE7
			factory.Load.AToHighRAM(),
			factory.Load.PopDByte(regs.HL),
			factory.Load.AToCPtrInHighRAM(),
			nil,
			nil,
			factory.Load.PushDByte(regs.HL),
			factory.ALU.AndD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst20h),

			// 0xE8 - 0xEF
			factory.ALU.AddR8ToSP(),
			factory.ControlFlow.JumpHL(),
			factory.Load.AToAddress(),
			nil,
			nil,
			nil,
			factory.ALU.XorD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst28h),

			// 0xF0 - 0xF7
			factory.Load.HighRAMToA(),
			factory.Load.PopDByte(regs.AF),
			factory.Load.CPtrInHighRAMToA(),
			diOpcode,
			nil,
			factory.Load.PushDByte(regs.AF),
			factory.ALU.OrD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst30h),

			// 0xF8 - 0xFF
			factory.Load.SPR8ToHL(),
			factory.Load.HLToSP(),
			factory.Load.AddressToA(),
			factory.Miscellaneous.EI(),
			nil,
			nil,
			factory.ALU.CompareD8ToA(),
			factory.ControlFlow.CallInterrupt(z80lib.Rst38h),
		},

		noopOpcode:        noopOpcode,
		diOpcode:          diOpcode,
		inPlaceInterrupts: inplaceInterrupts,
	}
}

func (c *CPU) doCycle(withDebug bool) error {
	clock := uint16(4)

	if !c.Regs.HaltMode {
		opcode, err := c.MMU.ReadByte(c.Regs.PC)
		if err != nil {
			return err
		}

		if withDebug {
			mem1, _ := c.MMU.ReadByte(c.Regs.PC + 1)
			mem2, _ := c.MMU.ReadByte(c.Regs.PC + 2)

			fmt.Printf(
				"PC: %04X | OPCODE: %02X | MEM: %02X%02X\n",
				c.Regs.PC, opcode, mem1, mem2)
		}

		opcodeResult := c.Opcodes[opcode]()
		if opcodeResult.Err != nil {
			return opcodeResult.Err
		}

		c.Regs.PC += opcodeResult.Size

		clock = opcodeResult.Clock
	}

	// GPU still works at normal speed even in double speed mode.
	gpuClock := clock
	if c.MMURegs.KEY1.CurrentSpeed.IsDoubleSpeed() {
		gpuClock /= 2
	}

	if err := c.GPU.DoCycles(gpuClock); err != nil {
		return err
	}

	c.manageTimers(clock)
	c.manageInterruptRequests()
	return nil
}

func (c *CPU) Run() error {
	withDebug := os.Getenv("DEBUG") != ""

	for !c.MMU.Registers().Stopped {
		if err := c.doCycle(withDebug); err != nil {
			return err
		}
	}
	return nil
}
