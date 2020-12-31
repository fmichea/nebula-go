package joypad

import (
	"fmt"
	"nebula-go/pkg/common/frontends"
	"nebula-go/pkg/gbc/memory"
	"nebula-go/pkg/gbc/memory/registers"
)

type Joypad interface{}

type joypad struct {
	mmuRegs *memory.Registers

	keyMapping map[frontends.KeyboardKey]*registers.JOYPButton
}

func New(mmu memory.MMU, window frontends.MainWindow) Joypad {
	mmuRegs := mmu.Registers()

	j := &joypad{
		mmuRegs: mmuRegs,

		keyMapping: map[frontends.KeyboardKey]*registers.JOYPButton{
			frontends.DownKey:  mmuRegs.JOYP.DownButton,
			frontends.UpKey:    mmuRegs.JOYP.UpButton,
			frontends.RightKey: mmuRegs.JOYP.RightButton,
			frontends.LeftKey:  mmuRegs.JOYP.LeftButton,

			frontends.SpaceKey:  mmuRegs.JOYP.StartButton,
			frontends.ReturnKey: mmuRegs.JOYP.SelectButton,
			frontends.AKey:      mmuRegs.JOYP.AButton,
			frontends.ZKey:      mmuRegs.JOYP.BButton,
		},
	}

	window.SubscribeKeyboardStateChanges(j.callback)

	return j
}

func (j *joypad) callback(key frontends.KeyboardKey, state frontends.KeyboardState) {
	if state == frontends.DownState && key == frontends.EscapeKey {
		j.mmuRegs.Stopped = true
	} else {
		flag, ok := j.keyMapping[key]
		if !ok {
			fmt.Println("Reiceived unknown key:", key)
			return
		}

		if state == frontends.DownState {
			flag.Press()
			j.mmuRegs.IF.Joypad.Request()
		} else {
			flag.Release()
		}
	}
}
