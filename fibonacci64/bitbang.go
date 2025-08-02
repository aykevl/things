package main

// #include "bitbang.h"
import "C"
import (
	"device/stm32"
	"machine"
	"unsafe"
)

func configureLEDs() {
	// Note: enabling a pullup would save a bit of current (because we'd avoid
	// floating inputs) but sadly also means some LEDs get turned on slightly
	// unintentionally.
	// It might be possible to fix this by controlling the MODER register
	// directly (instead of the ODR register), setting the pin to either output
	// mode (high/low depending on anode/cathode) or analog mode which disables
	// the output entirely.

	A1.High()
	A2.High()
	A3.High()
	A4.High()
	A5.High()
	A6.High()
	A7.High()
	A8.High()
	A9.High()
	A10.High()
	A11.High()
	A12.High()
	A13.High()
	A14.High()
	A15.High()
	A16.High()
	A17.High()
	A18.High()
	A19.High()
	A1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A3.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A4.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A5.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A6.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A7.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A8.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A9.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A10.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A11.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A12.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A13.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A14.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A15.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A16.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A17.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A18.Configure(machine.PinConfig{Mode: machine.PinOutput})
	A19.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func disableLEDs() {
	A1.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A2.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A3.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A4.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A5.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A6.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A7.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A8.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A9.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A10.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A11.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A12.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A13.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A14.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A15.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A16.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A17.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A18.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
	A19.Configure(machine.PinConfig{Mode: machine.PinInputAnalog})
}

func animateLEDs(mode, index, frame int) {
	if index < 12 {
		led0 := uint32(animate(mode, index+0, frame))
		led1 := uint32(animate(mode, index+12, frame))
		led2 := uint32(animate(mode, index+24, frame))
		switch index {
		case 0:
			C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[0][0])
		case 1:
			C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[1][0])
		case 2:
			C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[2][0])
		case 3:
			C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[3][0])
		case 4:
			C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[4][0])
		case 5:
			C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[5][0])
		case 6:
			C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[6][0])
		case 7:
			C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[7][0])
		case 8:
			C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[8][0])
		case 9:
			C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[9][0])
		case 10:
			C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[10][0])
		case 11:
			C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[11][0])
		}
	} else {

	}
}

var bitplanes [12][3]uint32

// Putting updateLEDs in RAM saves a bit of current consumption.
// TODO: this goes through a thunk, which adds a few cycles. GCC has
// __attribute__((long_call)) for ARM, perhaps we can also add this to Clang?
// (It's supported in Clang, but not in the ARM backend).
//
//go:section .ramfuncs.updateLEDs
func updateLEDs() {
	C.bitbang_show_leds(&bitplanes[0][0], unsafe.Pointer(stm32.GPIOA))
}
