package main

// #include "bitbang.h"
import "C"
import (
	"device/stm32"
	"machine"
	"unsafe"
)

const numLEDs = 64

func configureLEDs() {
	// Note: enabling a pullup would save a bit of current (because we'd avoid
	// floating inputs) but sadly also means some LEDs get turned on slightly
	// unintentionally.
	// It might be possible to fix this by controlling the MODER register
	// directly (instead of the ODR register), setting the pin to either output
	// mode (high/low depending on anode/cathode) or analog mode which disables
	// the output entirely.

	A1.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A4.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A5.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A6.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A7.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A8.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A9.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A10.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A11.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	A12.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

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
	A13.Low()
	A14.Low()
	A15.Low()
	A16.Low()
	A17.Low()
	A18.Low()
	A19.Low()

	// Set A1-A12 as open drain (and importantly, skip SWDIO/SWCLK)
	stm32.GPIOA.OTYPER.Set(0b0000_1111_1111_1111)
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
	switch index {
	case 0: // A4
		led0 := uint32(animate(mode, 36, frame))
		led1 := uint32(animate(mode, 58, frame))
		led2 := uint32(animate(mode, 16, frame))
		C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[0][0])
	case 1: // A5
		led0 := uint32(animate(mode, 49, frame))
		led1 := uint32(animate(mode, 37, frame))
		led2 := uint32(animate(mode, 29, frame))
		C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[1][0])
	case 2: // A6
		led0 := uint32(animate(mode, 62, frame))
		led1 := uint32(animate(mode, 50, frame))
		led2 := uint32(animate(mode, 21, frame))
		C.bitbang_update_bitplane_1(led0, led1, led2, &bitplanes[2][0])
	case 3: // A7
		led0 := uint32(animate(mode, 41, frame))
		led1 := uint32(animate(mode, 63, frame))
		led2 := uint32(animate(mode, 28, frame))
		C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[3][0])
	case 4: // A8
		led0 := uint32(animate(mode, 54, frame))
		led1 := uint32(animate(mode, 42, frame))
		led2 := uint32(animate(mode, 20, frame))
		C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[4][0])
	case 5: // A9
		led0 := uint32(animate(mode, 46, frame))
		led1 := uint32(animate(mode, 55, frame))
		led2 := uint32(animate(mode, 33, frame))
		C.bitbang_update_bitplane_2(led0, led1, led2, &bitplanes[5][0])
	case 6: // A10
		led0 := uint32(animate(mode, 59, frame))
		led1 := uint32(animate(mode, 22, frame))
		led2 := uint32(animate(mode, 25, frame))
		C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[6][0])
	case 7: // A11
		led0 := uint32(animate(mode, 38, frame))
		led1 := uint32(animate(mode, 14, frame))
		led2 := uint32(animate(mode, 17, frame))
		C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[7][0])
	case 8: // A12
		led0 := uint32(animate(mode, 51, frame))
		led1 := uint32(animate(mode, 27, frame))
		led2 := uint32(animate(mode, 9, frame))
		C.bitbang_update_bitplane_3(led0, led1, led2, &bitplanes[8][0])
	case 9: // A1
		led0 := uint32(animate(mode, 53, frame))
		led1 := uint32(animate(mode, 19, frame))
		led2 := uint32(animate(mode, 1, frame))
		C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[9][0])
	case 10: // A2
		led0 := uint32(animate(mode, 32, frame))
		led1 := uint32(animate(mode, 11, frame))
		led2 := uint32(animate(mode, 6, frame))
		C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[10][0])
	case 11: // A3
		led0 := uint32(animate(mode, 45, frame))
		led1 := uint32(animate(mode, 24, frame))
		led2 := uint32(animate(mode, 3, frame))
		C.bitbang_update_bitplane_4(led0, led1, led2, &bitplanes[11][0])
	case 12: // A13
		led0 := uint32(animate(mode, 30, frame))
		led1 := uint32(animate(mode, 8, frame))
		led2 := uint32(animate(mode, 13, frame))
		led3 := uint32(animate(mode, 34, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[12][0])
	case 13: // A14
		led0 := uint32(animate(mode, 43, frame))
		led1 := uint32(animate(mode, 0, frame))
		led2 := uint32(animate(mode, 26, frame))
		led3 := uint32(animate(mode, 47, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[13][0])
	case 14: // A15
		led0 := uint32(animate(mode, 56, frame))
		led1 := uint32(animate(mode, 2, frame))
		led2 := uint32(animate(mode, 18, frame))
		led3 := uint32(animate(mode, 60, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[14][0])
	case 15: // A16
		led0 := uint32(animate(mode, 35, frame))
		led1 := uint32(animate(mode, 15, frame))
		led2 := uint32(animate(mode, 5, frame))
		led3 := uint32(animate(mode, 39, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[15][0])
	case 16: // A17
		led0 := uint32(animate(mode, 48, frame))
		led1 := uint32(animate(mode, 7, frame))
		led2 := uint32(animate(mode, 10, frame))
		led3 := uint32(animate(mode, 52, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[16][0])
	case 17: // A18
		led0 := uint32(animate(mode, 61, frame))
		led1 := uint32(animate(mode, 12, frame))
		led2 := uint32(animate(mode, 23, frame))
		led3 := uint32(animate(mode, 31, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[17][0])
	case 18: // A19
		led0 := uint32(animate(mode, 40, frame))
		led1 := uint32(animate(mode, 4, frame))
		led2 := uint32(animate(mode, 57, frame))
		led3 := uint32(animate(mode, 44, frame))
		C.bitbang_update_bitplane_all4(led0, led1, led2, led3, &bitplanes[18][0])
	}
}

var bitplanes [19][3]uint32

// Putting updateLEDs in RAM saves a bit of current consumption.
// TODO: this goes through a thunk, which adds a few cycles. GCC has
// __attribute__((long_call)) for ARM, perhaps we can also add this to Clang?
// (It's
//
//go:section .ramfuncs.updateLEDs
func updateLEDs() {
	C.bitbang_show_leds(&bitplanes[0][0], unsafe.Pointer(stm32.GPIOA))
}
