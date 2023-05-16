//go:build attiny1616

package main

import (
	"runtime/volatile"
	"unsafe"

	"github.com/aykevl/tinygl/pixel"
)

// #include "bitbang.h"
import "C"

func makeColor[T pixel.Color](r, g, b uint8) T {
	return pixel.NewLinearColor[T](r, g, b)
}

// Hack until we fully support this chip.
type Port struct {
	DIR    volatile.Register8
	DIRSET volatile.Register8
	DIRCLR volatile.Register8
	DIRTGL volatile.Register8
	OUT    volatile.Register8
	OUTSET volatile.Register8
	OUTCLR volatile.Register8
	OUTTGL volatile.Register8
}

var (
	PORTA = (*Port)(unsafe.Pointer(uintptr(0x0400)))
	PORTB = (*Port)(unsafe.Pointer(uintptr(0x0420)))
	PORTC = (*Port)(unsafe.Pointer(uintptr(0x0440)))
)

var leds [18]pixel.LinearGRB888

func initLEDs() {
	PORTA.DIR.Set(0b0111_1110) // Configure PA1-PA6 as output.
	PORTB.DIR.Set(0b0011_1111) // Configure PB0-PB5 as output (R2, G2, B2, R3, G3, B3)
	PORTB.OUT.Set(0b0011_1111) // Set PB0-PB5 low
	PORTC.DIR.Set(0b0000_0111) // Configure PC0-PC2 as output (R1, G1, B1)
	PORTC.OUT.Set(0b0000_0111) // Set PC0-PC2 low
}

func updateLEDs() {
	// Pinout for the anodes:
	//   A1: PA2
	//   A2: PA1
	//   A3: PA6
	//   A4: PA5
	//   A5: PA4
	//   A6: PA3
	// Pinout for the cathodes:
	//   R1: PC0
	//   G1: PC1
	//   B1: PC2
	//   R2: PB0
	//   G2: PB1
	//   B2: PB2
	//   R3: PB3
	//   G3: PB4
	//   B3: PB5

	// R1
	PORTC.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[0+2].R,
		leds[0+3].R,
		leds[0+4].R,
		leds[0+5].R,
		leds[0+0].R,
		leds[0+1].R,
	)
	PORTC.OUTTGL.Set(1 << 0)

	// R2 (LEDs reversed)
	PORTB.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[6+3].R,
		leds[6+2].R,
		leds[6+1].R,
		leds[6+0].R,
		leds[6+5].R,
		leds[6+4].R,
	)
	PORTB.OUTTGL.Set(1 << 0)

	// R3
	PORTB.OUTTGL.Set(1 << 3)
	showLEDs(
		leds[12+2].R,
		leds[12+3].R,
		leds[12+4].R,
		leds[12+5].R,
		leds[12+0].R,
		leds[12+1].R,
	)
	PORTB.OUTTGL.Set(1 << 3)

	// G1
	PORTC.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[0+2].G,
		leds[0+3].G,
		leds[0+4].G,
		leds[0+5].G,
		leds[0+0].G,
		leds[0+1].G,
	)
	PORTC.OUTTGL.Set(1 << 1)

	// G2 (LEDs reversed)
	PORTB.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[6+3].G,
		leds[6+2].G,
		leds[6+1].G,
		leds[6+0].G,
		leds[6+5].G,
		leds[6+4].G,
	)
	PORTB.OUTTGL.Set(1 << 1)

	// G3
	PORTB.OUTTGL.Set(1 << 4)
	showLEDs(
		leds[12+2].G,
		leds[12+3].G,
		leds[12+4].G,
		leds[12+5].G,
		leds[12+0].G,
		leds[12+1].G,
	)
	PORTB.OUTTGL.Set(1 << 4)

	// B1
	PORTC.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[0+2].B,
		leds[0+3].B,
		leds[0+4].B,
		leds[0+5].B,
		leds[0+0].B,
		leds[0+1].B,
	)
	PORTC.OUTTGL.Set(1 << 2)

	// B2 (LEDs reversed)
	PORTB.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[6+3].B,
		leds[6+2].B,
		leds[6+1].B,
		leds[6+0].B,
		leds[6+5].B,
		leds[6+4].B,
	)
	PORTB.OUTTGL.Set(1 << 2)

	// B3
	PORTB.OUTTGL.Set(1 << 5)
	showLEDs(
		leds[12+2].B,
		leds[12+3].B,
		leds[12+4].B,
		leds[12+5].B,
		leds[12+0].B,
		leds[12+1].B,
	)
	PORTB.OUTTGL.Set(1 << 5)
}

// The bitbang function is very large, but without the //go:noinline it would be
// inlined. I guess LLVM doesn't recognize the inline assembly is very large.
//
//go:noinline
func showLEDs(c1, c2, c3, c4, c5, c6 uint8) {
	port := (*uint8)(unsafe.Pointer(uintptr(0x0001))) // VPORTA.OUT (alias of PORTA.OUT in I/O space)
	C.bitbang_show_leds(c1, c2, c3, c4, c5, c6, port)
}
