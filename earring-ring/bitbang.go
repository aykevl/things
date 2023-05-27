//go:build attiny1616

package main

import (
	"device/avr"
	"machine"
	"runtime/interrupt"
	"unsafe"
)

// #include "bitbang.h"
import "C"

func init() {
	// Use 20MHz/32 = 625kHz
	// This results in a current consumption of around 0.41mA with all LEDs off.
	// It *should* result in a current consumption of 187.5µA, but I think the
	// 20MHz oscillator uses a fair bit of current too.
	avr.CPU.CCP.Set(0xD8)                 // unlock protected registers
	avr.CLKCTRL.MCLKCTRLB.Set(0x4<<1 | 1) // prescaler of 32
}

const button = machine.PC3

func initHardware() {
	avr.PORTA.DIR.Set(0b0111_1110) // Configure PA1-PA6 as output.
	avr.PORTB.DIR.Set(0b0011_1111) // Configure PB0-PB5 as output (R2, G2, B2, R3, G3, B3)
	avr.PORTB.OUT.Set(0b0011_1111) // Set PB0-PB5 low
	avr.PORTC.DIR.Set(0b0000_0111) // Configure PC0-PC2 as output (R1, G1, B1)
	avr.PORTC.OUT.Set(0b0000_0111) // Set PC0-PC2 low

	button.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// This is the only remaining pin. Configure it as output just in case that
	// helps with standby current consumption.
	machine.PA7.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

func disableLEDs() {
	avr.PORTA.OUTCLR.Set(0b0111_1110) // Set PA1-PA6 high
	avr.PORTB.OUTCLR.Set(0b0011_1111) // Set PB0-PB5 high
	avr.PORTC.OUTCLR.Set(0b0000_0111) // Set PC0-PC2 high
}

func isButtonPressed() bool {
	return !button.Get()
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

	state := interrupt.Disable()

	// R1
	avr.PORTC.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[0+2].R,
		leds[0+3].R,
		leds[0+4].R,
		leds[0+5].R,
		leds[0+0].R,
		leds[0+1].R,
	)
	avr.PORTC.OUTTGL.Set(1 << 0)

	// R2 (LEDs reversed)
	avr.PORTB.OUTTGL.Set(1 << 0)
	showLEDs(
		leds[6+3].R,
		leds[6+2].R,
		leds[6+1].R,
		leds[6+0].R,
		leds[6+5].R,
		leds[6+4].R,
	)
	avr.PORTB.OUTTGL.Set(1 << 0)

	// R3
	avr.PORTB.OUTTGL.Set(1 << 3)
	showLEDs(
		leds[12+2].R,
		leds[12+3].R,
		leds[12+4].R,
		leds[12+5].R,
		leds[12+0].R,
		leds[12+1].R,
	)
	avr.PORTB.OUTTGL.Set(1 << 3)

	// G1
	avr.PORTC.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[0+2].G,
		leds[0+3].G,
		leds[0+4].G,
		leds[0+5].G,
		leds[0+0].G,
		leds[0+1].G,
	)
	avr.PORTC.OUTTGL.Set(1 << 1)

	// G2 (LEDs reversed)
	avr.PORTB.OUTTGL.Set(1 << 1)
	showLEDs(
		leds[6+3].G,
		leds[6+2].G,
		leds[6+1].G,
		leds[6+0].G,
		leds[6+5].G,
		leds[6+4].G,
	)
	avr.PORTB.OUTTGL.Set(1 << 1)

	// G3
	avr.PORTB.OUTTGL.Set(1 << 4)
	showLEDs(
		leds[12+2].G,
		leds[12+3].G,
		leds[12+4].G,
		leds[12+5].G,
		leds[12+0].G,
		leds[12+1].G,
	)
	avr.PORTB.OUTTGL.Set(1 << 4)

	// B1
	avr.PORTC.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[0+2].B,
		leds[0+3].B,
		leds[0+4].B,
		leds[0+5].B,
		leds[0+0].B,
		leds[0+1].B,
	)
	avr.PORTC.OUTTGL.Set(1 << 2)

	// B2 (LEDs reversed)
	avr.PORTB.OUTTGL.Set(1 << 2)
	showLEDs(
		leds[6+3].B,
		leds[6+2].B,
		leds[6+1].B,
		leds[6+0].B,
		leds[6+5].B,
		leds[6+4].B,
	)
	avr.PORTB.OUTTGL.Set(1 << 2)

	// B3
	avr.PORTB.OUTTGL.Set(1 << 5)
	showLEDs(
		leds[12+2].B,
		leds[12+3].B,
		leds[12+4].B,
		leds[12+5].B,
		leds[12+0].B,
		leds[12+1].B,
	)
	avr.PORTB.OUTTGL.Set(1 << 5)

	interrupt.Restore(state)
}

// The bitbang function is very large, but without the //go:noinline it would be
// inlined. I guess LLVM doesn't recognize the inline assembly is very large.
//
//go:noinline
func showLEDs(c1, c2, c3, c4, c5, c6 uint8) {
	port := (*uint8)(unsafe.Pointer(uintptr(0x0001))) // VPORTA.OUT (alias of PORTA.OUT in I/O space)
	C.bitbang_show_leds(c1, c2, c3, c4, c5, c6, port)
}
