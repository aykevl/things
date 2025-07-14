//go:build !bitbang3 && !bitbang2 && !bitbang1

package main

// #include "bitbang.h"
import "C"
import (
	"device/stm32"
	"runtime/interrupt"
	"unsafe"
)

func setLEDs(index int, led0, led1, led2 uint32) {
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
}

var bitplanes [12][2]uint32

// Putting updateLEDs in RAM saves a bit of current consumption.
//
//go:section .ramfuncs.updateLEDs
func updateLEDs() {
	mask := interrupt.Disable()

	otyper := &stm32.GPIOA.OTYPER
	out := &stm32.GPIOA.ODR

	// Update LED 0, 12, 24
	otyper.ClearBits(1 << 8) // clear bit for A4/PA8
	C.bitbang_show_leds(&bitplanes[0][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 8) // restore bits

	// Update LED 1, 13, 25
	otyper.ClearBits(1 << 7) // clear bit for A5/PA7
	C.bitbang_show_leds(&bitplanes[1][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 7) // restore bits

	// Update LED 2, 14, 26
	otyper.ClearBits(1 << 6) // clear bit for A6/PA6
	C.bitbang_show_leds(&bitplanes[2][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 6) // restore bits

	// Update LED 3, 15, 27
	otyper.ClearBits(1 << 5) // clear bit for A7/PA5
	C.bitbang_show_leds(&bitplanes[3][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 5) // restore bits

	// Update LED 4, 16, 28
	otyper.ClearBits(1 << 4) // clear bit for A8/PA4
	C.bitbang_show_leds(&bitplanes[4][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 4) // restore bits

	// Update LED 5, 17, 29
	otyper.ClearBits(1 << 3) // clear bit for A9/PA3
	C.bitbang_show_leds(&bitplanes[5][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 3) // restore bits

	// Update LED 6, 18, 30
	otyper.ClearBits(1 << 2) // clear bit for A10/PA2
	C.bitbang_show_leds(&bitplanes[6][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 2) // restore bits

	// Update LED 7, 19, 31
	otyper.ClearBits(1 << 1) // clear bit for A11/PA1
	C.bitbang_show_leds(&bitplanes[7][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 1) // restore bits

	// Update LED 8, 20, 32
	otyper.ClearBits(1 << 0) // clear bit for A12/PA0
	C.bitbang_show_leds(&bitplanes[8][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 0) // restore bits

	// Update LED 9, 21, 33
	otyper.ClearBits(1 << 15) // clear bit for A1/PA15
	C.bitbang_show_leds(&bitplanes[9][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 15) // restore bits

	// Update LED 10, 22, 34
	otyper.ClearBits(1 << 10) // clear bit for A2/PA10
	C.bitbang_show_leds(&bitplanes[10][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 10) // restore bits

	// Update LED 11, 23, 35
	otyper.ClearBits(1 << 9) // clear bit for A3/PA9
	C.bitbang_show_leds(&bitplanes[11][0], (*uint16)(unsafe.Pointer(out)))
	otyper.SetBits(1 << 9) // restore bits

	interrupt.Restore(mask)
}
